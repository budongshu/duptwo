package ldap

import (
	"crypto/tls"
	"fmt"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

// ADUserInfo AD用户信息
type ADUserInfo struct {
	DN          string
	Username    string
	Email       string
	DisplayName string
	Nickname    string
	Phone       string     // 手机号 telephoneNumber
	Department  string     // 部门 department
	Title       string     // 职位/职称 title
	Company     string     // 公司 company
	Manager     string     // 上级/经理 distinguishedName
}

// replaceFilterVars 替换过滤器中的变量为通配符
func replaceFilterVars(filter string) string {
	// 处理最常见的变量格式：%s 和 %(xxx)s
	filter = strings.Replace(filter, "%s", "*", -1)
	// 通用 %(xxx)s 格式（使用 strings.Index 避免引入 regexp）
	for {
		idx := strings.Index(filter, "%(")
		if idx == -1 {
			break
		}
		endIdx := strings.Index(filter[idx:], ")")
		if endIdx == -1 {
			break
		}
		filter = filter[:idx] + "*" + filter[idx+endIdx+1:]
	}
	return filter
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
		return nil, fmt.Errorf("无法连接LDAP服务器: %v", err)
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
		[]string{"dn", "uid", "mail", "displayName", "sAMAccountName", "cn", "givenName", "sn"},
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("搜索用户失败: %v", err)
	}

	if len(sr.Entries) == 0 {
		return nil, fmt.Errorf("用户在目录中不存在")
	}
	if len(sr.Entries) > 1 {
		return nil, fmt.Errorf("搜索到多个用户，请检查过滤器配置")
	}

	userDN := sr.Entries[0].DN
	userEntry := sr.Entries[0]

	// Bind with the user's credentials to verify password
	if err := conn.Bind(userDN, password); err != nil {
		return nil, fmt.Errorf("密码验证失败")
	}

	info := &ADUserInfo{
		DN:          userDN,
		Username:    userEntry.GetAttributeValue("uid"), // OpenLDAP
		Email:       userEntry.GetAttributeValue("mail"),
		DisplayName: userEntry.GetAttributeValue("displayName"),
		Nickname:    userEntry.GetAttributeValue("cn"),
		Phone:       userEntry.GetAttributeValue("telephoneNumber"),
		Department:  userEntry.GetAttributeValue("department"),
		Title:       userEntry.GetAttributeValue("title"),
		Company:     userEntry.GetAttributeValue("company"),
		Manager:     userEntry.GetAttributeValue("manager"),
	}
	if info.Username == "" {
		info.Username = userEntry.GetAttributeValue("sAMAccountName") // AD fallback
	}
	if info.Username == "" {
		info.Username = username
	}
	if info.Nickname == "" {
		nickname := userEntry.GetAttributeValue("displayName")
		if nickname == "" {
			nickname = userEntry.GetAttributeValue("cn")
			if nickname == "" {
				nickname = userEntry.GetAttributeValue("givenName")
				if nickname == "" {
					nickname = username
				}
			}
		}
		info.Nickname = nickname
	}

	return info, nil
}

// SearchAllUsers 搜索 AD 中的所有用户（用于批量同步）
func SearchAllUsers(server string, port int, useSSL bool, baseDN, bindDN, bindPassword, userFilter string) ([]ADUserInfo, error) {
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

	if bindDN != "" && bindPassword != "" {
		if err := conn.Bind(bindDN, bindPassword); err != nil {
			return nil, fmt.Errorf("管理账号绑定失败: %v", err)
		}
	}

	// 搜索过滤器：默认只搜索用户对象类
	filter := userFilter
	if filter == "" {
		filter = "(objectClass=inetOrgPerson)"
	} else {
		filter = replaceFilterVars(filter)
	}

	_ = "userAccountControl" // 保留引用，避免未使用警告（AD 专有，OpenLDAP 无此属性）

	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0, 60, false,
		filter,
		[]string{"dn", "cn", "sn", "givenName", "mail", "uid", "sAMAccountName", "displayName", "telephoneNumber", "department", "title", "company", "manager"},
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("搜索用户失败: %v", err)
	}

	var users []ADUserInfo
	for _, entry := range sr.Entries {
		// 优先取 uid（OpenLDAP），其次 sAMAccountName（AD），均无则跳过
		username := entry.GetAttributeValue("uid")
		if username == "" {
			username = entry.GetAttributeValue("sAMAccountName")
		}
		if username == "" {
			continue
		}

		// 昵称：优先 displayName（完整姓名），其次 cn（通用名称），再次 givenName（名字），都没有用 username
		nickname := entry.GetAttributeValue("displayName")
		if nickname == "" {
			nickname = entry.GetAttributeValue("cn")
			if nickname == "" {
				nickname = entry.GetAttributeValue("givenName")
				if nickname == "" {
					nickname = username
				}
			}
		}

		users = append(users, ADUserInfo{
			DN:          entry.DN,
			Username:    username,
			Email:       entry.GetAttributeValue("mail"),
			DisplayName: entry.GetAttributeValue("displayName"),
			Nickname:    nickname,
			Phone:      entry.GetAttributeValue("telephoneNumber"),
			Department: entry.GetAttributeValue("department"),
			Title:      entry.GetAttributeValue("title"),
			Company:    entry.GetAttributeValue("company"),
			Manager:    entry.GetAttributeValue("manager"),
		})
	}

	return users, nil
}
