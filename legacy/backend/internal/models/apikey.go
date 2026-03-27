package models

import "time"

// AdminAPIKey 管理员用于自动化/MCP 的长期密钥（仅存 SHA256 摘要，可设过期或永久）
type AdminAPIKey struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	UserID     uint       `gorm:"index;not null" json:"user_id"` // 创建者
	Name       string     `gorm:"size:128;not null" json:"name"`
	KeyHash    string     `gorm:"uniqueIndex;size:64;not null" json:"-"` // hex(sha256(完整密钥))
	Prefix     string     `gorm:"size:32;not null" json:"prefix"`          // 前缀展示，如 sp_a1b2c3d4
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`                   // nil 表示永不过期
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	IsActive   bool       `gorm:"default:true" json:"is_active"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

func (AdminAPIKey) TableName() string { return "admin_api_keys" }
