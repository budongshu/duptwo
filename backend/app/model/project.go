package model

import (
	"time"
)

// OnSiteStation 驻场点信息
type OnSiteStation struct {
	Location string `json:"location"` // 场地名称
	Person  string `json:"person"`  // 驻场人员
	Phone   string `json:"phone"`   // 联系方式
}

// Project 项目
type Project struct {
	BaseModel
	Name           string          `json:"name" gorm:"size:128;not null;index"`          // 项目名称
	Code          string          `json:"code" gorm:"size:64;uniqueIndex"`             // 项目编码
	Description   string          `json:"description" gorm:"size:512"`                  // 项目描述
	Status        string          `json:"status" gorm:"size:16;default:active"`          // 状态: active/inactive
	Stage         string          `json:"stage" gorm:"size:32;default:planning"`        // 阶段: planning/deploying/running/paused/archived
	Sort          int             `json:"sort" gorm:"default:0"`                         // 排序
	ProjectPerson string          `json:"projectPerson" gorm:"size:128"`                 // 项目人员
	OpsPerson     string          `json:"opsPerson" gorm:"size:128"`                    // 运维人员
	OpsStaffPerson string         `json:"opsStaffPerson" gorm:"size:128"`               // 运营人员
	Solution      string          `json:"solution" gorm:"size:1024"`                    // 解决方案描述
	SolutionPerson string         `json:"solutionPerson" gorm:"size:128"`               // 解决方案人员
	CompanyAddr   string          `json:"companyAddr" gorm:"size:256"`                  // 公司地点
	ProjectPeriod string          `json:"projectPeriod" gorm:"size:64"`                 // 项目周期
	OnsiteStations string         `json:"onsiteStations" gorm:"type:text"`             // 驻场点信息(JSON)
	IsDeleted     bool            `json:"-" gorm:"default:false"`                        // 软删除标记
	CreatedAt     time.Time       `json:"createdAt"`
	UpdatedAt     time.Time       `json:"updatedAt"`
}

func (Project) TableName() string {
	return "projects"
}
