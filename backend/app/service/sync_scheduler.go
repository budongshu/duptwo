package service

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"datauptwo/global"
	"datauptwo/app/repo"
)

// SyncScheduler 同步调度器
type SyncScheduler struct {
	mu            sync.RWMutex
	running       bool
	stopChan      chan struct{}
	queue         []SyncTask
	client        *SyncClient
	interval      time.Duration
	retryCount    int
	lastErrorType string
	lastErrorTime time.Time
	errorCount    int
	registered    bool
	lastRegisterErr string
	lastSyncAt    *time.Time       // 上次同步时间（断点）
	lastSerialNo  string           // 上次同步的最后一条 SerialNo
	batchSize     int             // 每批同步数量
	syncFilter    *SyncFilter     // 同步过滤器
	recordRepo    *repo.UploadRecordRepo // 用于查询待同步记录
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
	return &SyncScheduler{
		running:    false,
		stopChan:   make(chan struct{}),
		queue:      make([]SyncTask, 0),
		interval:   5 * time.Minute,
		retryCount: 3,
		batchSize:  100, // 默认每批100条
	}
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
	if global.CONF.Sync.Enabled {
		s.interval = parseDuration(global.CONF.Sync.Interval)
		s.retryCount = global.CONF.Sync.RetryCount

		// 读取批量大小配置
		if global.CONF.Sync.BatchSize > 0 {
			s.batchSize = global.CONF.Sync.BatchSize
		}

		// Agent 模式：自动注册到 Center
		if global.CONF.Sync.Mode == "agent" && global.CONF.Sync.CenterURL != "" {
			s.doAutoRegister()
		}
	}

	s.running = true
	s.stopChan = make(chan struct{})

	// 启动调度循环
	go s.run()

	global.AppLogger.Info("同步调度器已启动，模式: %s，间隔: %v，批量大小: %d",
		global.CONF.Sync.Mode, s.interval, s.batchSize)
	return nil
}

// doAutoRegister 自动注册到 Center 站点
func (s *SyncScheduler) doAutoRegister() {
	// 如果已有 API Key，直接使用
	if global.CONF.Sync.APIKey != "" {
		s.client = NewSyncClient(global.CONF.Sync.CenterURL, global.CONF.Sync.APIKey)
		s.registered = true
		global.AppLogger.Info("使用已有 API Key 连接到 Center")
		return
	}

	// 执行自动注册
	s.registerToCenter()
}

// registerToCenter 注册到 Center 站点
func (s *SyncScheduler) registerToCenter() {
	// 创建临时客户端（无 API Key）
	tempClient := NewSyncClient(global.CONF.Sync.CenterURL, "")

	// 准备注册请求
	registerReq := struct {
		StationCode string `json:"stationCode"`
		StationName string `json:"stationName"`
		URL         string `json:"url"`
	}{
		StationCode: global.CONF.Sync.StationID,
		StationName: global.CONF.Sync.StationName,
		URL:         fmt.Sprintf("http://localhost:%s", global.CONF.Sync.StationID),
	}

	respBody, statusCode, err := tempClient.Post("/api/sync/register", registerReq)
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
			StationID string `json:"stationId"`
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
	global.CONF.Sync.StationID = resp.Data.StationID

	s.client = NewSyncClient(global.CONF.Sync.CenterURL, apiKey)
	s.registered = true

	global.AppLogger.Info("自动注册成功，站点ID: %s", resp.Data.StationID)
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
		}
	}
}

// processQueue 处理队列
func (s *SyncScheduler) processQueue() {
	// Agent 未注册，跳过
	if global.CONF.Sync.Mode == "agent" && !s.registered {
		s.tryReRegister()
		return
	}

	if s.client == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

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

// RecordInfo 记录信息
type RecordInfo struct {
	SerialNo    string
	ProjectName string
	CreatedAt   time.Time
	Status      string
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
			CreatedAt:   r.CreatedAt,
			Status:      r.Status,
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

// tryReRegister 尝试重新注册
func (s *SyncScheduler) tryReRegister() {
	if s.registered {
		return
	}

	if time.Since(s.lastErrorTime) < s.interval {
		return
	}

	s.mu.Lock()
	s.registerToCenter()
	s.mu.Unlock()
}

// executeTask 执行同步任务
func (s *SyncScheduler) executeTask(task *SyncTask) error {
	record := SyncRecordData{
		SerialNo:    task.SerialNo,
		ProjectName: task.ProjectName,
	}

	req := &SyncUploadRequest{
		Records: []SyncRecordData{record},
	}

	resp, err := s.client.UploadRecords(req)
	if err != nil {
		return fmt.Errorf("同步失败: %w", err)
	}

	if resp.Code != 200 {
		return fmt.Errorf("同步返回错误: %s", resp.Message)
	}

	// 更新断点
	s.lastSerialNo = task.SerialNo
	now := time.Now()
	s.lastSyncAt = &now

	return nil
}

// handleTaskError 处理任务错误
func (s *SyncScheduler) handleTaskError(task *SyncTask, err error) {
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
		retryInterval := parseDuration(global.CONF.Sync.RetryInterval)
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
		"enabled":      global.CONF.Sync.Enabled,
		"mode":         global.CONF.Sync.Mode,
		"running":      s.running,
		"interval":     s.interval.String(),
		"batchSize":    s.batchSize,
		"registered":   s.registered,
		"lastSyncAt":   s.lastSyncAt,
		"lastSerialNo": s.lastSerialNo,
		"total":        len(s.queue),
		"pending":      pending,
		"processing":   processing,
		"completed":    completed,
		"failed":       failed,
		"lastErrorAt":   s.lastErrorTime,
		"filter":        s.syncFilter,
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