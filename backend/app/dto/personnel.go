package dto

import "time"

// ============ 人员相关 ============

// PersonnelCreateReq 创建人员请求
type PersonnelCreateReq struct {
	Name            string `json:"name" validate:"required,max=64"`
	Phone           string `json:"phone" validate:"max=32"`
	Email           string `json:"email" validate:"max=128"`
	Company         string `json:"company" validate:"max=128"`
	Position        string `json:"position" validate:"max=64"`
	WorkExperience  string `json:"workExperience" validate:"max=64"`
	EntryDate       string `json:"entryDate" validate:"max=32"`
	ProjectStartDate string `json:"projectStartDate" validate:"max=32"`
	OnProjectStatus string `json:"onProjectStatus" validate:"max=16"`
	Salary          string `json:"salary" validate:"max=32"`
	Location        string `json:"location" validate:"max=128"`
	Remark          string `json:"remark" validate:"max=256"`
	Status          string `json:"status"`
	Sort            int    `json:"sort"`
}

// PersonnelUpdateReq 更新人员请求
type PersonnelUpdateReq struct {
	ID              uint   `json:"id" validate:"required"`
	Name            string `json:"name" validate:"required,max=64"`
	Phone           string `json:"phone" validate:"max=32"`
	Email           string `json:"email" validate:"max=128"`
	Company         string `json:"company" validate:"max=128"`
	Position        string `json:"position" validate:"max=64"`
	WorkExperience  string `json:"workExperience" validate:"max=64"`
	EntryDate       string `json:"entryDate" validate:"max=32"`
	ProjectStartDate string `json:"projectStartDate" validate:"max=32"`
	OnProjectStatus string `json:"onProjectStatus" validate:"max=16"`
	Salary          string `json:"salary" validate:"max=32"`
	Location        string `json:"location" validate:"max=128"`
	Remark          string `json:"remark" validate:"max=256"`
	Status          string `json:"status"`
	Sort            int    `json:"sort"`
}

// PersonnelListReq 人员列表请求
type PersonnelListReq struct {
	Page      int    `form:"page" validate:"min=1"`
	PageSize  int    `form:"pageSize" validate:"min=1,max=100"`
	Keyword   string `form:"keyword"`
	Status    string `form:"status"`
	OnProject string `form:"onProject"`
}

// PersonnelResp 人员响应
type PersonnelResp struct {
	ID               uint      `json:"id"`
	Name             string    `json:"name"`
	Phone            string    `json:"phone"`
	Email            string    `json:"email"`
	Company          string    `json:"company"`
	Position         string    `json:"position"`
	WorkExperience   string    `json:"workExperience"`
	EntryDate        string    `json:"entryDate"`
	ProjectStartDate string    `json:"projectStartDate"`
	OnProjectStatus  string    `json:"onProjectStatus"`
	Salary           string    `json:"salary"`
	Location         string    `json:"location"`
	Remark           string    `json:"remark"`
	Status           string    `json:"status"`
	Sort             int       `json:"sort"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}
