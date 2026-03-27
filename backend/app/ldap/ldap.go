package ldap

import (
	"crypto/tls"
	"fmt"

	"github.com/go-ldap/ldap/v3"
)

// ADUserInfo AD用户信息
type ADUserInfo struct {
	DN          string
	Username    string
	Email       string
	DisplayName string
	Nickname    string
}

// Authenticate authenticates a user against LDAP
func Authenticate(server string, port int, useSSL bool, baseDN, bindDN, bindPassword, userFilter, username, password string) (*ADUserInfo, error) {
	address := fmt.Sprintf("%s:%d", server, port)

	var conn *ldap.Conn
	var err error
	if useSSL {
		conn, err = ldap.DialTLS("tcp", address, &tls.Config{InsecureSkipVerify: true})
	} else {
		conn, err = ldap.Dial("tcp", address)
	}
	if err != nil {
		return nil, fmt.Errorf("无法连接AD服务器: %v", err)
	}
	defer conn.Close()

	// Bind with admin credentials to search for the user
	if bindDN != "" && bindPassword != "" {
		if err := conn.Bind(bindDN, bindPassword); err != nil {
			return nil, fmt.Errorf("管理账号绑定失败: %v", err)
		}
	}

	// Search for the user DN
	filter := fmt.Sprintf(userFilter, ldap.EscapeFilter(username))
	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		1, 0, false,
		filter,
		[]string{"dn", "mail", "displayName", "sAMAccountName", "cn", "givenName", "sn"},
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("搜索用户失败: %v", err)
	}

	if len(sr.Entries) == 0 {
		return nil, fmt.Errorf("用户在AD中不存在")
	}
	if len(sr.Entries) > 1 {
		return nil, fmt.Errorf("搜索到多个用户，请检查AD配置")
	}

	userDN := sr.Entries[0].DN
	userEntry := sr.Entries[0]

	// Bind with the user's credentials to verify password
	if err := conn.Bind(userDN, password); err != nil {
		return nil, fmt.Errorf("密码验证失败")
	}

	info := &ADUserInfo{
		DN:          userDN,
		Username:    userEntry.GetAttributeValue("sAMAccountName"),
		Email:       userEntry.GetAttributeValue("mail"),
		DisplayName: userEntry.GetAttributeValue("displayName"),
		Nickname:    userEntry.GetAttributeValue("cn"),
	}

	if info.Username == "" {
		info.Username = username
	}
	if info.Nickname == "" {
		info.Nickname = info.DisplayName
		if info.Nickname == "" {
			info.Nickname = username
		}
	}

	return info, nil
}
