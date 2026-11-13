package models

import "time"

// BackupPayload 全量备份载荷（v1 兼容；v2 将扩展服务与场景）
type BackupPayload struct {
	Version            int                `json:"version"`
	ExportedAt         string             `json:"exported_at"`
	TopologyLayers     []TopologyLayer    `json:"topology_layers"`
	SolutionCategories []SolutionCategory `json:"solution_categories"`
	SolutionVendors    []SolutionVendor   `json:"solution_vendors"`
	Users              []BackupUserRow    `json:"users"`
	Products           []Product          `json:"products"`
	SalesScripts       []SalesScript      `json:"sales_scripts"`
	Cases              []Case             `json:"cases"`
	AuthSettings       *AuthSettings      `json:"auth_settings"`

	// v2 新增
	Services          []Service         `json:"services"`
	ServiceCategories []ServiceCategory `json:"service_categories"`
	ScenarioGroups    []ScenarioGroup   `json:"scenario_groups"`
	Scenarios         []Scenario        `json:"scenarios"`
	ScenarioProducts  []ScenarioProduct `json:"scenario_products"`
	ScenarioServices  []ScenarioService `json:"scenario_services"`
}

// BackupUserRow 备份中的用户行（含密码哈希，仅用于恢复）
type BackupUserRow struct {
	ID           uint       `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"password_hash"`
	AuthProvider string     `json:"auth_provider"`
	Role         string     `json:"role"`
	IsActive     bool       `json:"is_active"`
	LastLogin    *time.Time `json:"last_login,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
