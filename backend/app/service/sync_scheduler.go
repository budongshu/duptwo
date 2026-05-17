package service

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"datauptwo/app/repo"
	"datauptwo/global"
	"datauptwo/middleware"
)

// SyncScheduler 同步调度器
type SyncScheduler struct {
	mu              sync.RWMutex
	regMu           sync.RWMutex
	registerMu      sync.Mutex // 注册互斥锁，防止并发注册产生多个 key
	running         bool
	stopChan        chan struct{}
	processSignal   chan struct{} // 触发立即处理队列（401 后立即重试）
	retrySignal     chan struct{} // 触发重新注册（由 handleTaskError 触发，让主循环执行注册）
	queue           []SyncTask
	client          *SyncClient
	interval        time.Duration
	heartbeatInterval time.Duration
	retryCount      int
	lastErrorType   string
	lastErrorTime   time.Time
	errorCount      int
	registered      bool
	registering     bool   // 正在注册中，processQueue 应等待
	lastRegisterErr string
	lastSyncAt      *time.Time // 上次同步时间（断点）
	lastSerialNo    string     // 上次同步的最后一条 SerialNo
	batchSize       int       // 每批同步数量
	syncFilter      *SyncFilter // 同步过滤器
	recordRepo      *repo.UploadRecordRepo // 用于查询待同步记录
	stationRepo     *repo.SyncStationRepo // 用于 MySQL 持久化断点
}

// SyncTask 同步任务
type SyncTask struct {
	ID          uint       `json:"id"`
	SerialNo    string     `json:"serialNo"`
	ProjectName string     `json:"projectName"`
	Status      string     `json:"status"`
	RetryCount  int        `json:"retryCount"`
	MaxRetries  int        `json:"maxRetries"`
	NextRetryAt *time.Time `json:"nextRetryAt"`
	CreatedAt   time.Time  `json:"createdAt"`
	ErrorMsg    string     `json:"errorMsg"`
}

// SyncFilter 同步过滤器
type SyncFilter struct {
	ProjectNames []string // 只同步这些项目，空=全部
	StartTime    *time.Time
	EndTime      *time.Time
	Status       string // 只同步特定状态的记录
}

// NewSyncScheduler 创建同步调度器
func NewSyncScheduler() *SyncScheduler {
	s := &SyncScheduler{
		running:           false,
		stopChan:          make(chan struct{}),
		processSignal:     make(chan struct{}, 1), // 非阻塞，最多 1 个待处理信号
		retrySignal:       make(chan struct{}, 1),  // 触发重新注册（handleTaskError 通知主循环执行）
		queue:             make([]SyncTask, 0),
		interval:          5 * time.Minute,
		retryCount:        3,
		batchSize:         100, // 默认每批100条
		stationRepo:       repo.NewSyncStationRepo(),
	}
	// 从 MySQL 加载断点
	s.loadCheckpoint()
	return s
}

// RecordInfo 记录信息（完整字段）
type RecordInfo struct {
	SerialNo    string
	ProjectName string
	DiskLabel   string
	DestPath    string
	FileSize    int64
	Uploader    string
	Status      string
	Remark      string
	Data        string
	CreatedAt   time.Time
}

// SyncCheckpoint 断点数据结构（内存用）
type SyncCheckpoint struct {
	LastSyncAt   *time.Time `json:"lastSyncAt"`
	LastSerialNo string     `json:"lastSerialNo"`
}

// loadCheckpoint 从 MySQL 加载断点
func (s *SyncScheduler) loadCheckpoint() {
	s.mu.Lock()
	defer s.mu.Unlock()

	stationIDStr := global.CONF.Sync.Agent.StationID
	if stationIDStr == "" {
		global.AppLogger.Info("未配置 station_id，跳过断点加载")
		return
	}

	stationID, err := strconv.ParseUint(stationIDStr, 10, 64)
	if err != nil {
		global.AppLogger.Warn("无效的 station_id: %s，跳过断点加载", stationIDStr)
		return
	}

	station, err := s.stationRepo.GetByID(uint(stationID))
	if err != nil {
		global.AppLogger.Info("未找到站点记录，跳过断点加载: %v", err)
		return
	}

	s.lastSyncAt = station.LastSyncAt
	s.lastSerialNo = station.LastSerialNo
	global.AppLogger.Info("已从 MySQL 加载断点: lastSyncAt=%v, lastSerialNo=%s", station.LastSyncAt, station.LastSerialNo)
}

// saveCheckpoint 保存断点到 MySQL
func (s *SyncScheduler) saveCheckpoint() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.saveCheckpointLocked()
}

// saveCheckpointLocked 保存断点（内部方法，不加锁，需在持有锁时调用）
func (s *SyncScheduler) saveCheckpointLocked() {
	stationIDStr := global.CONF.Sync.Agent.StationID
	if stationIDStr == "" {
		return
	}

	stationID, err := strconv.ParseUint(stationIDStr, 10, 64)
	if err != nil {
		return
	}

	if err := s.stationRepo.UpdateCheckpoint(uint(stationID), s.lastSyncAt, s.lastSerialNo); err != nil {
		global.AppLogger.Error("保存断点失败: %v", err)
		return
	}
	global.AppLogger.Debug("断点已保存: lastSyncAt=%v, lastSerialNo=%s", s.lastSyncAt, s.lastSerialNo)
}

// SetFilter 设置同步过滤器
func (s *SyncScheduler) SetFilter(filter *SyncFilter) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.syncFilter = filter
}

// GetFilter 获取同步过滤器
func (s *SyncScheduler) GetFilter() *SyncFilter {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.syncFilter
}

// Start 启动调度器
func (s *SyncScheduler) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return fmt.Errorf("调度器已在运行")
	}

	// 从配置读取设置
	if global.CONF.Sync.Mode == "agent" {
		s.interval = parseDuration(global.CONF.Sync.Agent.Interval)
		s.retryCount = global.CONF.Sync.Agent.RetryCount

		// 读取批量大小配置
		if global.CONF.Sync.Agent.BatchSize > 0 {
			s.batchSize = global.CONF.Sync.Agent.BatchSize
		}

		// 读取心跳间隔（默认 60 秒）
		hb := global.CONF.Sync.Agent.HeartbeatInterval
		if hb == "" {
			hb = "60s"
		}
		s.heartbeatInterval = parseDuration(hb)

		// 自动注册到 Center
		if global.CONF.Sync.Agent.CenterURL != "" {
			s.doAutoRegister()
		}
	}

	s.running = true
	s.stopChan = make(chan struct{})

	// 启动调度循环
	go s.run()

	// 启动心跳上报（Agent 模式每 30 秒向 Center 发送心跳）
	if global.CONF.Sync.Mode == "agent" {
		go s.heartbeatLoop()
	}

	global.AppLogger.Info("同步调度器已启动，Mode: %s，间隔: %v，批量大小: %d",
		global.CONF.Sync.Mode, s.interval, s.batchSize)
	return nil
}

// doAutoRegister 自动注册到 Center 站点（优先级：环境变量 > app.yaml > MySQL）
func (s *SyncScheduler) doAutoRegister() {
	// 1. 最高优先：环境变量（Docker 部署推荐，容器重启不丢失）
	if apiKey := os.Getenv("SYNC_API_KEY"); apiKey != "" {
		s.regMu.Lock()
		s.client = NewSyncClient(global.CONF.Sync.Agent.CenterURL, apiKey)
		s.registered = true
		s.regMu.Unlock()
		global.AppLogger.Info("使用环境变量 SYNC_API_KEY 连接到 Center")
		return
	}

	// 2. 其次：app.yaml 中的配置
	if global.CONF.Sync.APIKey != "" {
		s.regMu.Lock()
		s.client = NewSyncClient(global.CONF.Sync.Agent.CenterURL, global.CONF.Sync.APIKey)
		s.registered = true
		s.regMu.Unlock()
		global.AppLogger.Info("使用 app.yaml 中的 API Key 连接到 Center")
		return
	}

	// 3. 最后：MySQL 中的明文 API Key（上次注册时存储的）
	stationIDStr := global.CONF.Sync.Agent.StationID
	if stationIDStr != "" {
		if stationID, err := strconv.ParseUint(stationIDStr, 10, 64); err == nil {
			if plainKey, err := s.stationRepo.GetPlainAPIKey(uint(stationID)); err == nil && plainKey != "" {
				global.CONF.Sync.APIKey = plainKey
				s.regMu.Lock()
				s.client = NewSyncClient(global.CONF.Sync.Agent.CenterURL, plainKey)
				s.registered = true
				s.regMu.Unlock()
				global.AppLogger.Info("使用 MySQL 中的 API Key 连接到 Center")
				return
			}
		}
	}

	// 4. 都没有，互斥执行注册（同步等待完成）
	s.safeRegisterAndBlock()
}

// safeRegister 线程安全的注册：使用互斥锁防止并发多次注册产生多个 key
func (s *SyncScheduler) safeRegister() {
	s.registerMu.Lock()
	defer s.registerMu.Unlock()

	// 再次检查：如果已经有 client 和 registered=true，说明已有其他 goroutine 注册成功
	s.regMu.RLock()
	alreadyRegistered := s.registered && s.client != nil
	s.regMu.RUnlock()
	if alreadyRegistered {
		global.AppLogger.Info("已有有效注册，跳过本次注册请求")
		return
	}

	// 标记正在注册，防止 doAutoRegister 的同步调用也进来
	s.regMu.Lock()
	s.registering = true
	s.regMu.Unlock()

	s.registerToCenter()

	// 注册完成，取消标记
	s.regMu.Lock()
	s.registering = false
	s.regMu.Unlock()
}

// safeRegisterAndBlock 互斥注册并等待完成（用于 doAutoRegister 同步路径）
func (s *SyncScheduler) safeRegisterAndBlock() {
	s.registerMu.Lock()
	defer s.registerMu.Unlock()

	// 再次检查
	s.regMu.RLock()
	alreadyRegistered := s.registered && s.client != nil
	s.regMu.RUnlock()
	if alreadyRegistered {
		global.AppLogger.Info("已有有效注册，跳过本次注册请求")
		return
	}

	// 标记正在注册
	s.regMu.Lock()
	s.registering = true
	s.regMu.Unlock()

	s.registerToCenter()

	// 注册完成，取消标记
	s.regMu.Lock()
	s.registering = false
	s.regMu.Unlock()
}

// registerToCenter 注册到 Center 站点
func (s *SyncScheduler) registerToCenter() {
	// 创建临时客户端（无 API Key）
	tempClient := NewSyncClient(global.CONF.Sync.Agent.CenterURL, "")

	// 准备注册请求
	registerReq := struct {
		StationCode string `json:"stationCode"`
		StationName string `json:"stationName"`
		URL         string `json:"url"`
	}{
		StationCode: global.CONF.Sync.Agent.StationID,
		StationName: global.CONF.Sync.Agent.StationName,
		URL:         global.CONF.Sync.Agent.URL,
	}

	respBody, statusCode, err := tempClient.Post("/api/sync/register", registerReq)

	// 注册状态写入需要加锁（防止与 heartbeatLoop/processQueue 并发读写 registered/client）
	s.regMu.Lock()
	defer s.regMu.Unlock()

	if err != nil {
		errMsg := fmt.Sprintf("自动注册失败: %v", err)
		if s.lastRegisterErr != errMsg || time.Since(s.lastErrorTime) > time.Hour {
			global.AppLogger.Error(errMsg)
			s.lastRegisterErr = errMsg
			s.lastErrorTime = time.Now()
		}
		s.registered = false
		return
	}

	if statusCode != 200 {
		errMsg := fmt.Sprintf("自动注册失败，状态码: %d，响应: %s", statusCode, string(respBody))
		if s.lastRegisterErr != errMsg || time.Since(s.lastErrorTime) > time.Hour {
			global.AppLogger.Error(errMsg)
			s.lastRegisterErr = errMsg
			s.lastErrorTime = time.Now()
		}
		s.registered = false
		return
	}

	var resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			StationID uint   `json:"stationId"`
			APIKey    string `json:"apiKey"`
		} `json:"data"`
	}

	if err := json.Unmarshal(respBody, &resp); err != nil {
		global.AppLogger.Error("解析注册响应失败: %v", err)
		s.registered = false
		return
	}

	if resp.Data.APIKey == "" {
		global.AppLogger.Error("注册成功但未返回 API Key")
		s.registered = false
		return
	}

	apiKey := resp.Data.APIKey
	global.CONF.Sync.APIKey = apiKey
	global.CONF.Sync.Agent.StationID = strconv.FormatUint(uint64(resp.Data.StationID), 10)

	// 持久化明文 API Key 到 MySQL（与断点同表）
	if stationID, err := strconv.ParseUint(global.CONF.Sync.Agent.StationID, 10, 64); err == nil {
		s.stationRepo.UpdatePlainAPIKey(uint(stationID), apiKey)
	}

	s.client = NewSyncClient(global.CONF.Sync.Agent.CenterURL, apiKey)
	s.registered = true

	// 刷新 Center 的 API Key 缓存，使新站点立即可用（无缓存则下次心跳才刷新）
	middleware.RefreshCache()

	global.AppLogger.Info("自动注册成功，站点ID: %d", resp.Data.StationID)
}

// Stop 停止调度器
func (s *SyncScheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return
	}

	close(s.stopChan)
	s.running = false

	if s.client != nil {
		s.client.Close()
	}

	global.AppLogger.Info("同步调度器已停止")
}

// run 调度循环
func (s *SyncScheduler) run() {
	// 首次执行延迟，等待注册完成
	if global.CONF.Sync.Mode == "agent" && !s.registered {
		time.Sleep(5 * time.Second)
	}

	// 立即执行一次同步（不要等待定时器）
	go s.processQueue()

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			s.processQueue()
		case <-s.processSignal:
			s.processQueue()
		case <-s.retrySignal:
			s.safeRegisterAndBlock()
			// 注册完成后，立即触发队列处理
			select {
			case s.processSignal <- struct{}{}:
			default:
			}
		}
	}
}

// heartbeatLoop 心跳循环（已注册时按配置间隔发送心跳）
func (s *SyncScheduler) heartbeatLoop() {
	ticker := time.NewTicker(s.heartbeatInterval)
	defer ticker.Stop()
	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			s.sendHeartbeat()
		}
	}
}

// sendHeartbeat 发送心跳到 Center
func (s *SyncScheduler) sendHeartbeat() {
	s.regMu.RLock()
	registered := s.registered
	client := s.client
	s.regMu.RUnlock()

	if !registered || client == nil {
		return
	}

	if err := client.Heartbeat(); err != nil {
		// 检测 401：说明当前 key 已失效（可能是旧 key，或 Center 重启后 key 丢失）
		// 清除状态并通过 retrySignal 通知主循环执行注册（避免并发多个 goroutine 同时注册）
		if strings.Contains(err.Error(), "401") || strings.Contains(err.Error(), "无效的 API Key") {
			global.AppLogger.Warn("心跳返回 401，API Key 已失效，清除状态，等待重新注册")
			s.regMu.Lock()
			s.registered = false
			s.client = nil
			s.regMu.Unlock()
			// 清除旧 API Key，防止后续请求使用
			oldKey := global.CONF.Sync.APIKey
			global.CONF.Sync.APIKey = ""
			// 清除 MySQL 中的明文 Key
			stationIDStr := global.CONF.Sync.Agent.StationID
			if stationIDStr != "" {
				if stationID, err := strconv.ParseUint(stationIDStr, 10, 64); err == nil {
					s.stationRepo.UpdatePlainAPIKey(uint(stationID), "")
				}
			}
			_ = oldKey
			// 通过 retrySignal 通知主循环执行注册
			select {
			case s.retrySignal <- struct{}{}:
			default:
			}
			return
		}

		if s.errorCount == 0 || time.Since(s.lastErrorTime) > time.Hour {
			global.AppLogger.Warn("心跳发送失败: %v", err)
			s.lastErrorTime = time.Now()
		}
		s.errorCount++
	} else if s.errorCount > 0 {
		// 心跳成功后重置错误计数
		s.errorCount = 0
	}
}

// processQueue 处理队列
func (s *SyncScheduler) processQueue() {
	// Agent 未注册，跳过（等待注册完成）
	if global.CONF.Sync.Mode == "agent" {
		s.regMu.RLock()
		notRegistered := !s.registered
		registering := s.registering
		s.regMu.RUnlock()
		if notRegistered || registering {
			// 正在注册中，等待后再检查；或未注册，尝试注册
			if registering {
				global.AppLogger.Info("正在注册中，等待...")
			}
			s.tryReRegister()
			return
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// 再次检查 client（持有队列锁后）
	s.regMu.RLock()
	client := s.client
	s.regMu.RUnlock()
	if client == nil {
		return
	}

	// 获取待处理的任务
	var pendingTasks []SyncTask
	for _, task := range s.queue {
		if task.Status == "pending" || (task.Status == "failed" && task.RetryCount < task.MaxRetries) {
			if task.NextRetryAt != nil && time.Now().Before(*task.NextRetryAt) {
				continue
			}
			pendingTasks = append(pendingTasks, task)
		}
	}

	if len(pendingTasks) == 0 {
		// 没有排队的任务，扫描新增记录
		s.scanAndAddRecords()
		return
	}

	if s.errorCount == 0 {
		global.AppLogger.Info("开始处理同步任务，共 %d 个", len(pendingTasks))
	}

	for _, task := range pendingTasks {
		if err := s.executeTask(&task); err != nil {
			s.handleTaskError(&task, err)
		} else {
			s.markTaskCompleted(&task)
		}
	}
}

// scanAndAddRecords 扫描并添加新记录到队列（增量同步）
func (s *SyncScheduler) scanAndAddRecords() {
	// 获取需要同步的记录（从上次同步位置之后）
	records := s.fetchRecordsToSync()
	if len(records) == 0 {
		// 打印断点信息，方便排查扫描为空的原因
		lastSyncStr := "<nil>"
		if s.lastSyncAt != nil {
			lastSyncStr = s.lastSyncAt.Format("2006-01-02 15:04:05.000")
		}
		global.AppLogger.Debug("[Scan] 无待同步记录，断点: lastSyncAt=%s, lastSerialNo=%s", lastSyncStr, s.lastSerialNo)
		return
	}

	global.AppLogger.Info("发现 %d 条新记录待同步", len(records))

	for _, record := range records {
		// 检查是否已在队列中
		exists := false
		for _, t := range s.queue {
			if t.SerialNo == record.SerialNo && t.Status != "completed" {
				exists = true
				break
			}
		}
		if !exists {
			task := SyncTask{
				ID:          uint(len(s.queue) + 1),
				SerialNo:    record.SerialNo,
				ProjectName: record.ProjectName,
				Status:      "pending",
				RetryCount:  0,
				MaxRetries:  s.retryCount,
				CreatedAt:   time.Now(),
			}
			s.queue = append(s.queue, task)
		}
	}
}

// fetchRecordsToSync 获取需要同步的记录（从数据库查询，增量同步）
func (s *SyncScheduler) fetchRecordsToSync() []RecordInfo {
	if s.recordRepo == nil {
		s.recordRepo = repo.NewUploadRecordRepo()
	}

	// 获取项目过滤列表
	var projectNames []string
	if s.syncFilter != nil && len(s.syncFilter.ProjectNames) > 0 {
		projectNames = s.syncFilter.ProjectNames
	}

	// 查询记录
	records, err := s.recordRepo.GetRecordsSince(s.lastSyncAt, s.lastSerialNo, projectNames, s.batchSize)
	if err != nil {
		if s.errorCount == 0 || time.Since(s.lastErrorTime) > time.Hour {
			global.AppLogger.Error("查询待同步记录失败: %v", err)
			s.lastErrorTime = time.Now()
			s.errorCount++
		}
		return nil
	}

	result := make([]RecordInfo, len(records))
	for i, r := range records {
		result[i] = RecordInfo{
			SerialNo:    r.SerialNo,
			ProjectName: r.ProjectName,
			DiskLabel:   r.DiskLabel,
			DestPath:    r.DestPath,
			FileSize:    r.FileSize,
			Uploader:    r.Uploader,
			Status:      r.Status,
			Remark:      r.Remark,
			Data:        r.Data,
			CreatedAt:   r.CreatedAt,
		}
	}
	return result
}

// AddTask 手动添加同步任务
func (s *SyncScheduler) AddTask(serialNo, projectName string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, task := range s.queue {
		if task.SerialNo == serialNo && task.Status != "completed" {
			return
		}
	}

	task := SyncTask{
		ID:          uint(len(s.queue) + 1),
		SerialNo:    serialNo,
		ProjectName: projectName,
		Status:      "pending",
		RetryCount:  0,
		MaxRetries:  s.retryCount,
		CreatedAt:   time.Now(),
	}

	s.queue = append(s.queue, task)
}

// tryReRegister 尝试重新注册（不在持有队列锁时调用，避免与 processQueue 争抢）
func (s *SyncScheduler) tryReRegister() {
	s.regMu.RLock()
	registered := s.registered
	registering := s.registering
	s.regMu.RUnlock()
	if registered || registering {
		return
	}

	if time.Since(s.lastErrorTime) < s.interval {
		return
	}

	// safeRegister 使用互斥锁，防止并发多次注册
	s.safeRegister()
}

// executeTask 执行同步任务
func (s *SyncScheduler) executeTask(task *SyncTask) error {
	// 每次执行时从 DB 查完整记录，确保发送最新数据
	if s.recordRepo == nil {
		s.recordRepo = repo.NewUploadRecordRepo()
	}
	dbRecord, err := s.recordRepo.GetBySerialNo(task.SerialNo)
	if err != nil {
		global.AppLogger.Error("同步任务 [%s] 执行失败：本地记录不存在或已删除，%v", task.SerialNo, err)
		return fmt.Errorf("记录不存在或已删除: %s", task.SerialNo)
	}

	// 解析 Data JSON 字段
	var data map[string]interface{}
	if dbRecord.Data != "" {
		if parseErr := json.Unmarshal([]byte(dbRecord.Data), &data); parseErr != nil {
			global.AppLogger.Warn("记录 [%s] Data 字段解析失败: %v", task.SerialNo, parseErr)
		}
	}

	record := SyncRecordData{
		SerialNo:    dbRecord.SerialNo,
		ProjectName: dbRecord.ProjectName,
		DiskLabel:   dbRecord.DiskLabel,
		DestPath:    dbRecord.DestPath,
		FileSize:    dbRecord.FileSize,
		Uploader:    dbRecord.Uploader,
		Status:      dbRecord.Status,
		Remark:      dbRecord.Remark,
		Data:        data,
	}

	req := &SyncUploadRequest{
		Records: []SyncRecordData{record},
	}

	global.AppLogger.Info("[Agent] 推送记录到 Center: %s/api/sync/upload-records, serialNo=%s, project=%s",
		global.CONF.Sync.Agent.CenterURL, record.SerialNo, record.ProjectName)

	s.regMu.RLock()
	client := s.client
	s.regMu.RUnlock()
	if client == nil {
		return fmt.Errorf("client 未初始化")
	}

	resp, err := client.UploadRecords(req)
	if err != nil {
		global.AppLogger.Error("同步任务 [%s] 执行失败：网络错误，%v", task.SerialNo, err)
		return err
	}

	if resp.Code != 200 {
		global.AppLogger.Error("同步任务 [%s] 执行失败：Center 返回错误 [%d] %s", task.SerialNo, resp.Code, resp.Message)
		return fmt.Errorf("同步返回错误: %s", resp.Message)
	}

	// 检查响应体中的处理结果，不只是 HTTP 状态码
	if resp.Data.FailCount > 0 {
		global.AppLogger.Warn("同步任务 [%s] 部分失败：success=%d, fail=%d, conflict=%d",
			task.SerialNo, resp.Data.SuccessCount, resp.Data.FailCount, resp.Data.ConflictCount)
		return fmt.Errorf("Center 处理失败: success=%d, fail=%d", resp.Data.SuccessCount, resp.Data.FailCount)
	}

	// 更新断点
	s.lastSerialNo = task.SerialNo
	now := time.Now()
	s.lastSyncAt = &now
	s.saveCheckpointLocked()

	global.AppLogger.Info("[Agent] 同步任务 [%s] 成功，Center 已接收", task.SerialNo)
	return nil
}

// handleTaskError 处理任务错误
func (s *SyncScheduler) handleTaskError(task *SyncTask, err error) {
	// 检测 401：API Key 无效 → 触发重新注册，不再重试该任务
	if strings.Contains(err.Error(), "401") || strings.Contains(err.Error(), "无效的 API Key") {
		global.AppLogger.Warn("同步任务 [%s] 返回 401，API Key 已失效，将触发重新注册", task.SerialNo)

		s.regMu.Lock()
		s.registered = false
		s.client = nil
		s.regMu.Unlock()
		global.CONF.Sync.APIKey = ""
		// 清除 MySQL 中的明文 Key
		stationIDStr := global.CONF.Sync.Agent.StationID
		if stationIDStr != "" {
			if stationID, err := strconv.ParseUint(stationIDStr, 10, 64); err == nil {
				s.stationRepo.UpdatePlainAPIKey(uint(stationID), "")
			}
		}

		// 任务重置为 pending，等待重新注册后重试
		task.Status = "pending"
		task.RetryCount = 0
		task.ErrorMsg = "API Key 失效，等待重新注册"
		task.NextRetryAt = nil
		s.updateTask(task)

		// 通过 retrySignal 通知主循环执行注册（避免并发多个 goroutine 同时注册）
		select {
		case s.retrySignal <- struct{}{}:
		default:
		}
		return
	}

	task.RetryCount++
	task.ErrorMsg = err.Error()

	errType := fmt.Sprintf("%s:%s", task.ProjectName, err.Error()[:min(50, len(err.Error()))])

	if task.RetryCount >= task.MaxRetries {
		task.Status = "failed"
		if s.lastErrorType != errType || time.Since(s.lastErrorTime) > time.Hour {
			global.AppLogger.Error("同步任务 %s 失败，已达到最大重试次数: %s", task.SerialNo, err.Error())
			s.lastErrorType = errType
			s.lastErrorTime = time.Now()
			s.errorCount++
		}
	} else {
		retryInterval := parseDuration(global.CONF.Sync.Agent.RetryInterval)
		nextRetry := time.Now().Add(retryInterval)
		task.NextRetryAt = &nextRetry

		if s.lastErrorType != errType || time.Since(s.lastErrorTime) > time.Hour {
			global.AppLogger.Warn("同步任务 %s 失败，将在 %v 后重试 (第 %d/%d 次): %s",
				task.SerialNo, retryInterval, task.RetryCount, task.MaxRetries, err.Error())
			s.lastErrorType = errType
			s.lastErrorTime = time.Now()
		}
	}

	s.updateTask(task)
}

// markTaskCompleted 标记任务完成
func (s *SyncScheduler) markTaskCompleted(task *SyncTask) {
	task.Status = "completed"
	task.NextRetryAt = nil
	s.updateTask(task)
}

// updateTask 更新队列中的任务
func (s *SyncScheduler) updateTask(updatedTask *SyncTask) {
	for i, task := range s.queue {
		if task.ID == updatedTask.ID {
			s.queue[i] = *updatedTask
			return
		}
	}
}

// GetQueueStatus 获取队列状态
func (s *SyncScheduler) GetQueueStatus() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.regMu.RLock()
	registered := s.registered
	s.regMu.RUnlock()

	pending := 0
	processing := 0
	completed := 0
	failed := 0

	for _, task := range s.queue {
		switch task.Status {
		case "pending":
			pending++
		case "processing":
			processing++
		case "completed":
			completed++
		case "failed":
			failed++
		}
	}

	return map[string]interface{}{
		"enabled":      global.CONF.Sync.Mode == "agent",
		"running":      s.running,
		"interval":     s.interval.String(),
		"batchSize":    s.batchSize,
		"registered":   registered,
		"lastSyncAt":   s.lastSyncAt,
		"lastSerialNo": s.lastSerialNo,
		"total":        len(s.queue),
		"pending":      pending,
		"processing":   processing,
		"completed":    completed,
		"failed":       failed,
		"lastErrorAt":  s.lastErrorTime,
		"lastError":    s.lastErrorType,
		"filter":       s.syncFilter,
	}
}

// GetQueue 获取队列中的所有任务
func (s *SyncScheduler) GetQueue() []SyncTask {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]SyncTask, len(s.queue))
	copy(result, s.queue)
	return result
}

// ClearCompleted 清除已完成的任务
func (s *SyncScheduler) ClearCompleted() {
	s.mu.Lock()
	defer s.mu.Unlock()

	var newQueue []SyncTask
	for _, task := range s.queue {
		if task.Status != "completed" {
			newQueue = append(newQueue, task)
		}
	}
	s.queue = newQueue
}

// RetryFailed 重试失败的任务
func (s *SyncScheduler) RetryFailed() {
	s.mu.Lock()
	defer s.mu.Unlock()

	count := 0
	for i, task := range s.queue {
		if task.Status == "failed" {
			s.queue[i].Status = "pending"
			s.queue[i].RetryCount = 0
			s.queue[i].NextRetryAt = nil
			count++
		}
	}

	if count > 0 {
		global.AppLogger.Info("已重置 %d 个失败任务", count)
	}
}

// parseDuration 解析时间字符串
func parseDuration(s string) time.Duration {
	if s == "" {
		return 5 * time.Minute
	}

	d, err := time.ParseDuration(s)
	if err != nil {
		return 5 * time.Minute
	}

	return d
}
