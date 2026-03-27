package auth

// LdapConfig LDAP 连接与搜索参数（存 JSONB）
type LdapConfig struct {
	Host          string `json:"host"`
	Port          int    `json:"port"`
	UseTLS        bool   `json:"use_tls"`
	BindDN        string `json:"bind_dn"`
	BindPassword  string `json:"bind_password"`
	UserSearchBase string `json:"user_search_base"`
	UserFilter    string `json:"user_filter"` // 包含 %s 作为用户名占位
}

// OidcConfig OIDC 客户端配置
type OidcConfig struct {
	Issuer       string `json:"issuer"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURL  string `json:"redirect_url"`
}
