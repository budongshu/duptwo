package dto

import "time"

// ============ 项目相关 ============

// OnSiteStation 驻场点信息
type OnSiteStation struct {
	Location string `json:"location"` // 场地名称
	Person  string `json:"person"`  // 驻场人员
	Phone   string `json:"phone"`   // 联系方式
}

// ProjectCreateReq 创建项目请求
type ProjectCreateReq struct {
	Name           string           `json:"name" validate:"required,max=128"`
	Code          string           `json:"code" validate:"required,max=64"`
	Description   string           `json:"description" validate:"max=512"`
	Status        string           `json:"status"`
	Stage         string           `json:"stage"`
	Sort          int              `json:"sort"`
	ProjectPerson string           `json:"projectPerson"`    // 项目人员
	OpsPerson     string           `json:"opsPerson"`       // 运维人员
	OpsStaffPerson string          `json:"opsStaffPerson"`  // 运营人员
	Solution      string           `json:"solution"`        // 解决方案描述
	SolutionPerson string          `json:"solutionPerson"`  // 解决方案人员
	CompanyAddr   string           `json:"companyAddr"`    // 公司地点
	ProjectPeriod string           `json:"projectPeriod"`  // 项目周期
	OnsiteStations []OnSiteStation `json:"onsiteStations"` // 驻场点列表
}

// ProjectUpdateReq 更新项目请求
type ProjectUpdateReq struct {
	ID            uint             `json:"id" validate:"required"`
	Name           string           `json:"name" validate:"required,max=128"`
	Code          string           `json:"code" validate:"required,max=64"`
	Description   string           `json:"description" validate:"max=512"`
	Status        string           `json:"status"`
	Stage         string           `json:"stage"`
	Sort          int              `json:"sort"`
	ProjectPerson string           `json:"projectPerson"`    // 项目人员
	OpsPerson     string           `json:"opsPerson"`       // 运维人员
	OpsStaffPerson string          `json:"opsStaffPerson"`  // 运营人员
	Solution      string           `json:"solution"`        // 解决方案描述
	SolutionPerson string          `json:"solutionPerson"`  // 解决方案人员
	CompanyAddr   string           `json:"companyAddr"`    // 公司地点
	ProjectPeriod string           `json:"projectPeriod"`   // 项目周期
	OnsiteStations []OnSiteStation `json:"onsiteStations"`  // 驻场点列表
}

// ProjectListReq 项目列表请求
type ProjectListReq struct {
	Page     int    `form:"page" validate:"min=1"`
	PageSize int    `form:"pageSize" validate:"min=1,max=100"`
	Keyword  string `form:"keyword"`
	Status   string `form:"status"`
	Stage    string `form:"stage"`
}

// ProjectResp 项目响应
type ProjectResp struct {
	ID            uint             `json:"id"`
	Name           string           `json:"name"`
	Code          string           `json:"code"`
	Description   string           `json:"description"`
	Status        string           `json:"status"`
	Stage         string           `json:"stage"`
	Sort          int              `json:"sort"`
	ProjectPerson string           `json:"projectPerson"`    // 项目人员
	OpsPerson     string           `json:"opsPerson"`       // 运维人员
	OpsStaffPerson string          `json:"opsStaffPerson"`  // 运营人员
	Solution      string           `json:"solution"`        // 解决方案描述
	SolutionPerson string          `json:"solutionPerson"`  // 解决方案人员
	CompanyAddr   string           `json:"companyAddr"`     // 公司地点
	ProjectPeriod string           `json:"projectPeriod"`   // 项目周期
	OnsiteStations []OnSiteStation `json:"onsiteStations"`   // 驻场点列表
	RecordCount   int64            `json:"recordCount"`     // 该项目的上传记录数量
	TotalDataSize int64            `json:"totalDataSize"`    // 该项目的上传数据总量(字节)
	CreatedAt     time.Time        `json:"createdAt"`
	UpdatedAt     time.Time        `json:"updatedAt"`
}

// ProjectSimpleResp 简单项目响应（用于下拉选择）
type ProjectSimpleResp struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}
