package service

import (
	"fmt"

	"datauptwo/app/ldap"
	"datauptwo/global"
)

// ADService AD/LDAP 认证服务
type ADService struct {
}

// NewADService 创建AD服务
func NewADService() *ADService {
	return &ADService{}
}

// Authenticate 使用LDAP认证用户
func (s *ADService) Authenticate(username, password string) (*ldap.ADUserInfo, error) {
	cfg := global.CONF.AD
	if !cfg.Enabled {
		return nil, fmt.Errorf("AD认证未启用")
	}

	if username == "" || password == "" {
		return nil, fmt.Errorf("用户名和密码不能为空")
	}

	port := cfg.Port
	if port == 0 {
		port = 389
	}

	return ldap.Authenticate(
		cfg.Server,
		port,
		cfg.UseSSL,
		cfg.BaseDN,
		cfg.BindDN,
		cfg.BindPassword,
		cfg.UserFilter,
		username,
		password,
	)
}

// IsEnabled AD认证是否启用
func (s *ADService) IsEnabled() bool {
	return global.CONF.AD.Enabled
}
