package model

import "time"

// Personnel 人员
type Personnel struct {
	BaseModel
	Name            string    `json:"name" gorm:"size:64;not null;index"`              // 姓名
	Phone           string    `json:"phone" gorm:"size:32"`                             // 手机号
	Email           string    `json:"email" gorm:"size:128"`                            // 邮箱
	Company         string    `json:"company" gorm:"size:128"`                          // 所属公司
	Position        string    `json:"position" gorm:"size:64"`                          // 职位
	WorkExperience  string    `json:"workExperience" gorm:"size:64"`                    // 工作经验/工龄
	EntryDate       string    `json:"entryDate" gorm:"size:32"`                         // 入项时间
	ProjectStartDate string   `json:"projectStartDate" gorm:"size:32"`                  // 立项时间
	OnProjectStatus string    `json:"onProjectStatus" gorm:"size:16;default:离项"`      // 是否在项: 在项/离项
	Salary          string    `json:"salary" gorm:"size:32"`                            // 薪资
	Location        string    `json:"location" gorm:"size:128"`                         // 人员驻场地点
	Remark          string    `json:"remark" gorm:"size:256"`                           // 备注
	Status          string    `json:"status" gorm:"size:16;default:active"`             // 状态: active/inactive
	Sort            int       `json:"sort" gorm:"default:0"`                             // 排序
	IsDeleted       bool      `json:"-" gorm:"default:false"`                           // 软删除标记
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func (Personnel) TableName() string {
	return "personnels"
}
