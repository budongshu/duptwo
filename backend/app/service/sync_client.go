package service

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"datauptwo/global"
)

// SyncClient 同步HTTP客户端（带代理支持）
type SyncClient struct {
	client   *http.Client
	baseURL  string
	apiKey   string
	proxyURL string
}

// NewSyncClient 创建同步客户端
func NewSyncClient(baseURL, apiKey string) *SyncClient {
	client := &SyncClient{
		baseURL: baseURL,
		apiKey:  apiKey,
	}

	// 配置代理
	proxyEnabled := global.CONF.Sync.Proxy.Enabled
	proxyURL := global.CONF.Sync.Proxy.URL

	if proxyEnabled && proxyURL != "" {
		client.proxyURL = proxyURL
	}

	// 创建HTTP客户端
	client.client = &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
			},
		},
	}

	// 设置代理
	if client.proxyURL != "" {
		proxy, err := url.Parse(client.proxyURL)
		if err == nil {
			client.client.Transport.(*http.Transport).Proxy = http.ProxyURL(proxy)
		}
	}

	return client
}

// Request 发送HTTP请求
func (c *SyncClient) Request(method, path string, body interface{}) ([]byte, int, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("序列化请求失败: %w", err)
		}
		reqBody = bytes.NewReader(jsonBytes)
	}

	req, err := http.NewRequest(method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, 0, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)

	// 发送请求
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("读取响应失败: %w", err)
	}

	return respBody, resp.StatusCode, nil
}

// Get 发送GET请求
func (c *SyncClient) Get(path string) ([]byte, int, error) {
	return c.Request("GET", path, nil)
}

// Post 发送POST请求
func (c *SyncClient) Post(path string, body interface{}) ([]byte, int, error) {
	return c.Request("POST", path, body)
}

// Put 发送PUT请求
func (c *SyncClient) Put(path string, body interface{}) ([]byte, int, error) {
	return c.Request("PUT", path, body)
}

// Delete 发送DELETE请求
func (c *SyncClient) Delete(path string) ([]byte, int, error) {
	return c.Request("DELETE", path, nil)
}

// ============ 同步客户端方法 ============

// SyncUploadRequest 上传记录同步请求结构
type SyncUploadRequest struct {
	StationCode string `json:"stationCode"`
	StationName string `json:"stationName"`
	Records     []SyncRecordData `json:"records"`
}

// SyncRecordData 同步记录数据
type SyncRecordData struct {
	SerialNo    string                 `json:"serialNo"`
	ProjectName string                 `json:"projectName"`
	DiskLabel   string                 `json:"diskLabel"`
	DestPath    string                 `json:"destPath"`
	FileSize    int64                  `json:"fileSize"`
	Uploader    string                 `json:"uploader"`
	Status      string                 `json:"status"`
	Remark      string                 `json:"remark"`
	Data        map[string]interface{} `json:"data"`
}

// SyncUploadResponse 上传记录同步响应
type SyncUploadResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		TotalRecords   int `json:"totalRecords"`
		SuccessCount   int `json:"successCount"`
		FailCount      int `json:"failCount"`
		ConflictCount  int `json:"conflictCount"`
	} `json:"data"`
}

// UploadRecords 上传记录到中心站点
func (c *SyncClient) UploadRecords(req *SyncUploadRequest) (*SyncUploadResponse, error) {
	respBody, statusCode, err := c.Post("/api/sync/upload-records", req)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，状态码: %d，响应: %s", statusCode, string(respBody))
	}

	var resp SyncUploadResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &resp, nil
}

// SyncStatusResponse 同步状态响应
type SyncStatusResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Enabled     bool   `json:"enabled"`
		Mode        string `json:"mode"`
		IsCenter    bool   `json:"isCenter"`
		StationID   string `json:"stationId"`
		StationName string `json:"stationName"`
	} `json:"data"`
}

// GetStatus 获取中心站点状态
func (c *SyncClient) GetStatus() (*SyncStatusResponse, error) {
	respBody, statusCode, err := c.Get("/api/sync/status")
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，状态码: %d，响应: %s", statusCode, string(respBody))
	}

	var resp SyncStatusResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &resp, nil
}

// PingResponse Ping响应
type PingResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Ping ping中心站点
func (c *SyncClient) Ping() (*PingResponse, error) {
	respBody, statusCode, err := c.Get("/health")
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("站点不可达，状态码: %d", statusCode)
	}

	var resp PingResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		// 健康检查返回的格式可能不同，尝试简单解析
		resp.Code = 200
		resp.Message = "ok"
	}

	return &resp, nil
}

// Close 关闭客户端
func (c *SyncClient) Close() {
	if c.client != nil && c.client.Transport != nil {
		if transport, ok := c.client.Transport.(*http.Transport); ok {
			transport.CloseIdleConnections()
		}
	}
}
