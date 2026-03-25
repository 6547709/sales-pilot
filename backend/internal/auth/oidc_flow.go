package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// OidcStateStore 简易内存 state 存储（生产可换 Redis）
type OidcStateStore struct {
	mu    sync.Mutex
	items map[string]time.Time
}

func NewOidcStateStore() *OidcStateStore {
	s := &OidcStateStore{items: map[string]time.Time{}}
	go s.gc()
	return s
}

func (s *OidcStateStore) gc() {
	t := time.NewTicker(2 * time.Minute)
	for range t.C {
		s.mu.Lock()
		now := time.Now()
		for k, exp := range s.items {
			if now.After(exp) {
				delete(s.items, k)
			}
		}
		s.mu.Unlock()
	}
}

// Put 写入 state，约 10 分钟有效
func (s *OidcStateStore) Put(state string) {
	s.mu.Lock()
	s.items[state] = time.Now().Add(10 * time.Minute)
	s.mu.Unlock()
}

// Take 取出并删除，不存在或过期返回 false
func (s *OidcStateStore) Take(state string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	exp, ok := s.items[state]
	if !ok || time.Now().After(exp) {
		return false
	}
	delete(s.items, state)
	return true
}

// RandomState 生成随机 state
func RandomState() (string, error) {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// OidcEndpoints 根据配置构造 OAuth2
func OidcEndpoints(ctx context.Context, cfg OidcConfig) (*oauth2.Config, *oidc.Provider, error) {
	provider, err := oidc.NewProvider(ctx, cfg.Issuer)
	if err != nil {
		return nil, nil, err
	}
	oauthCfg := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}
	return oauthCfg, provider, nil
}

// OidcUserInfo OIDC 回调后解析的用户信息
type OidcUserInfo struct {
	Sub           string
	Email         string
	PreferredName string
}

// ExchangeAndVerify 用授权码换 token 并校验 id_token，返回 claims
func ExchangeAndVerify(ctx context.Context, oauthCfg *oauth2.Config, provider *oidc.Provider, code string) (*OidcUserInfo, error) {
	oauth2Token, err := oauthCfg.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("oauth2 exchange: %w", err)
	}
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("无 id_token")
	}
	verifier := provider.Verifier(&oidc.Config{ClientID: oauthCfg.ClientID})
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, err
	}
	var claims map[string]any
	if err := idToken.Claims(&claims); err != nil {
		return nil, err
	}
	ui := &OidcUserInfo{Sub: idToken.Subject}
	if v, ok := claims["email"].(string); ok {
		ui.Email = v
	}
	if v, ok := claims["preferred_username"].(string); ok {
		ui.PreferredName = v
	}
	if ui.PreferredName == "" && ui.Email != "" {
		ui.PreferredName = ui.Email
	}
	return ui, nil
}

// ParseOidcConfig 从 JSON 字节解析
func ParseOidcConfig(b []byte) (OidcConfig, error) {
	var c OidcConfig
	if len(b) == 0 {
		return c, fmt.Errorf("空配置")
	}
	if err := json.Unmarshal(b, &c); err != nil {
		return c, err
	}
	return c, nil
}

// ParseLdapConfig 从 JSON 字节解析
func ParseLdapConfig(b []byte) (LdapConfig, error) {
	var c LdapConfig
	if len(b) == 0 {
		return c, fmt.Errorf("空配置")
	}
	if err := json.Unmarshal(b, &c); err != nil {
		return c, err
	}
	return c, nil
}
