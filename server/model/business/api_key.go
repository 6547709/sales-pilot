package business

import "time"

// APIKey API 密钥（用于 AI/MCP 机器访问）
type APIKey struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	Name       string     `gorm:"size:128;not null" json:"name"`
	KeyHash    string     `gorm:"size:255;not null" json:"-"` // 不返回给前端
	IsActive   bool       `gorm:"default:true" json:"is_active"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	CreatedBy  uint       `gorm:"index" json:"created_by"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

func (APIKey) TableName() string {
	return "api_keys"
}
