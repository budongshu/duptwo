package model

import "time"

// SecuritySettings 安全设置（全局唯一，只有一条记录）
type SecuritySettings struct {
	ID uint `json:"id" gorm:"primarykey"`

	// 登录验证码
	CaptchaEnabled bool `json:"captchaEnabled" gorm:"default:true"`     // 是否启用登录验证码
	CaptchaMinLen  int  `json:"captchaMinLen" gorm:"default:3"`         // 验证码最少错误次数后启用

	// 不活跃用户自动禁用
	InactiveAutoDisable     bool `json:"inactiveAutoDisable" gorm:"default:false"`     // 是否启用不活跃自动禁用
	InactiveDaysThreshold   int  `json:"inactiveDaysThreshold" gorm:"default:90"`      // 不活跃天数阈值

	// 用户级别登录限制
	UserLoginMaxAttempts    int  `json:"userLoginMaxAttempts" gorm:"default:5"`       // 用户最大连续失败次数
	UserLoginLockMinutes    int  `json:"userLoginLockMinutes" gorm:"default:30"`       // 用户被锁定时长（分钟）

	// IP级别登录限制
	IPLoginMaxAttempts      int  `json:"ipLoginMaxAttempts" gorm:"default:20"`         // IP最大连续失败次数
	IPLoginLockMinutes      int  `json:"ipLoginLockMinutes" gorm:"default:60"`         // IP被锁定时长（分钟）

	// IP白名单（逗号分隔，支持CIDR）
	IPWhitelist             string `json:"ipWhitelist" gorm:"size:2048"`               // IP白名单，白名单IP不受限制

	// IP黑名单（逗号分隔，支持CIDR）
	IPBlacklist             string `json:"ipBlacklist" gorm:"size:2048"`               // IP黑名单，黑名单IP禁止登录

	// 密码安全
	PasswordExpiryDays      int  `json:"passwordExpiryDays" gorm:"default:0"`          // 密码过期天数（0=永不过期）
	PasswordMinLength       int  `json:"passwordMinLength" gorm:"default:6"`           // 密码最小长度
	PasswordRequireUppercase bool `json:"passwordRequireUppercase" gorm:"default:false"` // 必须包含大写字母
	PasswordRequireLowercase bool `json:"passwordRequireLowercase" gorm:"default:false"` // 必须包含小写字母
	PasswordRequireDigit    bool  `json:"passwordRequireDigit" gorm:"default:false"`     // 必须包含数字
	PasswordRequireSpecial  bool  `json:"passwordRequireSpecial" gorm:"default:false"`   // 必须包含特殊字符

	// 会话配置
	SessionTimeoutHours     int  `json:"sessionTimeoutHours" gorm:"default:24"`        // 会话超时时长（小时）

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (SecuritySettings) TableName() string {
	return "security_settings"
}

// LoginLockout 登录锁定记录
type LoginLockout struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	Target    string    `json:"target" gorm:"size:128;index"` // 用户名 或 IP地址
	Type      string    `json:"type" gorm:"size:16;index"`     // user / ip
	FailCount int       `json:"failCount" gorm:"default:1"`   // 连续失败次数
	LockedAt  *time.Time `json:"lockedAt"`                    // 锁定时间（被锁后才记录）
	Locked    bool      `json:"locked" gorm:"default:false"` // 是否已锁定
	UnlockedAt *time.Time `json:"unlockedAt"`                 // 解锁时间（自动解锁或手动解锁时更新）
	Reason    string    `json:"reason" gorm:"size:256"`       // 锁定原因
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (LoginLockout) TableName() string {
	return "login_lockouts"
}
