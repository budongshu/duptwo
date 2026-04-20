package service

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"datauptwo/global"
)

// SyncScheduler 同步调度器
type SyncScheduler struct {
	mu              sync.RWMutex
	running         bool
	stopChan        chan struct{}
	queue           []SyncTask
	client          *SyncClient
	interval        time.Duration
	retryCount      int
	lastErrorType   string          // 上次错误类型，用于去重告警
	lastErrorTime   time.Time       // 上次错误时间
	errorCount      int             // 连续错误次数
	registered      bool            // 是否已注册到 Center
	lastRegisterErr string          // 注册错误信息（避免重复告警）
}

// SyncTask 同步任务
type SyncTask struct {
	ID          uint       `json:"id"`
	SerialNo    string     `json:"serialNo"`
	ProjectName string     `json:"projectName"`
	Status      string     `json:"status"` // pending/processing/completed/failed
	RetryCount  int        `json:"retryCount"`
	MaxRetries  int        `json:"maxRetries"`
	NextRetryAt *time.Time `json:"nextRetryAt"`
	CreatedAt   time.Time  `json:"createdAt"`
	ErrorMsg    string     `json:"errorMsg"`
}

// NewSyncScheduler 创建同步调度器
func NewSyncScheduler() *SyncScheduler {
	return &SyncScheduler{
		running:    false,
		stopChan:   make(chan struct{}),
		queue:      make([]SyncTask, 0),
		interval:   5 * time.Minute,
		retryCount: 3,
	}
}

// Start 启动调度器（Agent模式自动注册）
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

		// Agent 模式：自动注册到 Center
		if global.CONF.Sync.Mode == "agent" && global.CONF.Sync.CenterURL != "" {
			s.doAutoRegister()
		}
	}

	s.running = true
	s.stopChan = make(chan struct{})

	// 启动调度循环
	go s.run()

	global.AppLogger.Info("同步调度器已启动，模式: %s，间隔: %v", global.CONF.Sync.Mode, s.interval)
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
		StationID   string `json:"stationId"`
		StationName string `json:"stationName"`
		URL         string `json:"url"`
	}{
		StationID:   global.CONF.Sync.StationID,
		StationName: global.CONF.Sync.StationName,
		URL:         fmt.Sprintf("http://localhost:%s", global.CONF.Sync.StationID), // 本地地址
	}

	// 调用注册接口
	respBody, statusCode, err := tempClient.Post("/api/sync/register", registerReq)
	if err != nil {
		errMsg := fmt.Sprintf("自动注册失败: %v", err)
		// 避免重复告警：相同错误1小时内只告警一次
		if s.lastRegisterErr != errMsg || time.Since(s.lastErrorTime) > time.Hour {
			global.AppLogger.Error(errMsg)
			s.lastRegisterErr = errMsg
			s.lastErrorTime = time.Now()
		}
		// 注册失败不影响启动，继续运行
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

	// 解析响应获取 API Key
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

	// 保存 API Key 到配置（通过 viper）
	apiKey := resp.Data.APIKey
	global.CONF.Sync.APIKey = apiKey
	global.CONF.Sync.StationID = resp.Data.StationID

	// 创建客户端
	s.client = NewSyncClient(global.CONF.Sync.CenterURL, apiKey)
	s.registered = true

	global.AppLogger.Info("自动注册成功，站点ID: %s", resp.Data.StationID)
	global.AppLogger.Info("API Key 已保存，请在配置文件中备份: %s...", apiKey[:16])
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
		// 尝试重新注册
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
			// 检查是否到了重试时间
			if task.NextRetryAt != nil && time.Now().Before(*task.NextRetryAt) {
				continue
			}
			pendingTasks = append(pendingTasks, task)
		}
	}

	if len(pendingTasks) == 0 {
		return
	}

	// 只在首次或有新任务时记录
	if s.errorCount == 0 {
		global.AppLogger.Info("开始处理同步任务，共 %d 个", len(pendingTasks))
	}

	// 处理每个任务
	for _, task := range pendingTasks {
		if err := s.executeTask(&task); err != nil {
			s.handleTaskError(&task, err)
		} else {
			s.markTaskCompleted(&task)
		}
	}
}

// tryReRegister 尝试重新注册
func (s *SyncScheduler) tryReRegister() {
	if s.registered {
		return
	}

	// 每5分钟尝试一次注册
	if time.Since(s.lastErrorTime) < s.interval {
		return
	}

	s.mu.Lock()
	s.registerToCenter()
	s.mu.Unlock()
}

// executeTask 执行同步任务
func (s *SyncScheduler) executeTask(task *SyncTask) error {
	// 准备同步数据
	record := SyncRecordData{
		SerialNo:    task.SerialNo,
		ProjectName: task.ProjectName,
	}

	// 调用中心站点API
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

	return nil
}

// handleTaskError 处理任务错误（去重告警）
func (s *SyncScheduler) handleTaskError(task *SyncTask, err error) {
	task.RetryCount++
	task.ErrorMsg = err.Error()

	// 获取错误类型标识（用于去重）
	errType := fmt.Sprintf("%s:%s", task.ProjectName, err.Error()[:min(50, len(err.Error()))])

	if task.RetryCount >= task.MaxRetries {
		task.Status = "failed"
		// 去重告警：相同错误1小时内只记录一次
		if s.lastErrorType != errType || time.Since(s.lastErrorTime) > time.Hour {
			global.AppLogger.Error("同步任务 %s 失败，已达到最大重试次数: %s", task.SerialNo, err.Error())
			s.lastErrorType = errType
			s.lastErrorTime = time.Now()
		}
	} else {
		// 设置下次重试时间
		retryInterval := parseDuration(global.CONF.Sync.RetryInterval)
		nextRetry := time.Now().Add(retryInterval)
		task.NextRetryAt = &nextRetry

		// 去重告警：相同错误1小时内只记录一次
		if s.lastErrorType != errType || time.Since(s.lastErrorTime) > time.Hour {
			global.AppLogger.Warn("同步任务 %s 失败，将在 %v 后重试 (第 %d/%d 次): %s",
				task.SerialNo, retryInterval, task.RetryCount, task.MaxRetries, err.Error())
			s.lastErrorType = errType
			s.lastErrorTime = time.Now()
		}
	}

	// 更新队列中的任务
	s.updateTask(task)
}

// markTaskCompleted 标记任务完成
func (s *SyncScheduler) markTaskCompleted(task *SyncTask) {
	task.Status = "completed"
	task.NextRetryAt = nil
	s.updateTask(task)
	// 不打印成功日志，减少噪音
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

// AddTask 添加同步任务到队列
func (s *SyncScheduler) AddTask(serialNo, projectName string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查是否已存在
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
	// 不打印添加任务日志，减少噪音
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
		"enabled":     global.CONF.Sync.Enabled,
		"mode":        global.CONF.Sync.Mode,
		"running":     s.running,
		"interval":   s.interval.String(),
		"registered":  s.registered,
		"total":       len(s.queue),
		"pending":     pending,
		"processing":  processing,
		"completed":   completed,
		"failed":      failed,
		"lastErrorAt": s.lastErrorTime,
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

// ============ 辅助函数 ============

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

// MarshalJSON 自定义序列化
func (t SyncTask) MarshalJSON() ([]byte, error) {
	type Alias SyncTask
	return json.Marshal(&struct {
		Alias
		CreatedAtStr string `json:"createdAt"`
	}{
		Alias:       Alias(t),
		CreatedAtStr: t.CreatedAt.Format(time.RFC3339),
	})
}