package business

import "time"

// SolutionCategory 解决方案分类
type SolutionCategory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Label     string    `gorm:"size:128;not null" json:"label"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (SolutionCategory) TableName() string {
	return "solution_categories"
}
