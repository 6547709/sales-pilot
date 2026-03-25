package models

import (
	"time"

	"gorm.io/datatypes"
)

// User 系统用户（本地 / LDAP / OIDC 统一落库或映射）
type User struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Username     string     `gorm:"uniqueIndex;size:128" json:"username"`
	Email        string     `gorm:"size:255" json:"email"`
	PasswordHash string     `gorm:"size:255" json:"-"`
	AuthProvider string     `gorm:"size:32;default:local" json:"auth_provider"` // local, ldap, oidc
	Role         string     `gorm:"size:32;default:user" json:"role"`           // admin, user
	IsActive     bool       `gorm:"default:true" json:"is_active"`
	LastLogin    *time.Time `json:"last_login,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// Product 销售赋能产品内容
type Product struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	Name               string         `gorm:"size:255;not null" json:"name"`
	Category           string         `gorm:"size:128" json:"category"`
	Description        string         `gorm:"type:text" json:"description"`
	Highlights         datatypes.JSON `gorm:"type:jsonb" json:"highlights"`           // []string，3 大亮点
	TargetPersonas     datatypes.JSON `gorm:"type:jsonb" json:"target_personas"`
	TriggerEvents      string         `gorm:"type:text" json:"trigger_events"`
	DiscoveryQuestions datatypes.JSON `gorm:"type:jsonb" json:"discovery_questions"` // 黄金三问
	CompetitorAnalysis string         `gorm:"type:text" json:"competitor_analysis"`
	ROIMetrics         string         `gorm:"type:text" json:"roi_metrics"`
	UpdatedAt          time.Time      `json:"updated_at"`
	CreatedAt          time.Time      `json:"created_at"`
}

// SalesScript 场景话术
type SalesScript struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProductID uint      `gorm:"index;not null" json:"product_id"`
	Scenario  string    `gorm:"size:255" json:"scenario"`
	Content   string    `gorm:"type:text" json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Case 客户案例
type Case struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ProductID      uint      `gorm:"index;not null" json:"product_id"`
	ClientName     string    `gorm:"size:255" json:"client_name"`
	PainPoints     string    `gorm:"type:text" json:"pain_points"`
	Solution       string    `gorm:"type:text" json:"solution"`
	ValueDelivered string    `gorm:"type:text" json:"value_delivered"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// AuthSettings 管理员可配置的认证参数（单例 id=1）
type AuthSettings struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	LdapEnabled bool           `json:"ldap_enabled"`
	LdapConfig  datatypes.JSON `gorm:"type:jsonb" json:"ldap_config"`
	OidcEnabled bool           `json:"oidc_enabled"`
	OidcConfig  datatypes.JSON `gorm:"type:jsonb" json:"oidc_config"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
