package auth

import (
	"fmt"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

// TryLDAP 使用服务账号搜索用户 DN，再以用户凭据绑定验证密码
func TryLDAP(cfg LdapConfig, username, password string) error {
	if cfg.Host == "" || cfg.UserSearchBase == "" {
		return fmt.Errorf("ldap 未正确配置")
	}
	port := cfg.Port
	if port == 0 {
		port = 389
	}
	addr := fmt.Sprintf("%s:%d", cfg.Host, port)

	var conn *ldap.Conn
	var err error
	if cfg.UseTLS {
		conn, err = ldap.DialTLS("tcp", addr, nil)
	} else {
		conn, err = ldap.Dial("tcp", addr)
	}
	if err != nil {
		return err
	}
	defer conn.Close()

	if cfg.BindDN != "" {
		if err := conn.Bind(cfg.BindDN, cfg.BindPassword); err != nil {
			return fmt.Errorf("ldap 绑定失败: %w", err)
		}
	}

	filter := cfg.UserFilter
	if filter == "" {
		filter = "(uid=%s)"
	}
	if strings.Contains(filter, "%s") {
		filter = fmt.Sprintf(filter, ldap.EscapeFilter(username))
	} else {
		filter = strings.ReplaceAll(filter, "{{username}}", ldap.EscapeFilter(username))
	}

	req := ldap.NewSearchRequest(
		cfg.UserSearchBase,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{"dn"},
		nil,
	)
	sr, err := conn.Search(req)
	if err != nil {
		return err
	}
	if len(sr.Entries) != 1 {
		return fmt.Errorf("用户未找到或匹配多条")
	}
	userDN := sr.Entries[0].DN

	// 以用户身份二次绑定验证密码
	if err := conn.Bind(userDN, password); err != nil {
		return fmt.Errorf("ldap 用户认证失败")
	}
	return nil
}
