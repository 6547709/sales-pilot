package config

import (
	"os"
	"strings"
)

// Config 运行时配置
type Config struct {
	DatabaseURL    string
	JWTSecret      string
	AppBaseURL     string // 后端对外基址，用于 OIDC redirect
	FrontendOrigin string
	GinMode        string
}

// Load 从环境变量加载
func Load() Config {
	return Config{
		DatabaseURL:    getenv("DATABASE_URL", "postgres://sales:sales@localhost:5432/sales_pilot?sslmode=disable"),
		JWTSecret:      getenv("JWT_SECRET", "dev-change-me-in-production"),
		AppBaseURL:     strings.TrimRight(getenv("APP_BASE_URL", "http://localhost:8080"), "/"),
		FrontendOrigin: getenv("FRONTEND_ORIGIN", "http://localhost:3000"),
		GinMode:        getenv("GIN_MODE", "debug"),
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
