package business

import "time"

// SalesScript 场景话术
type SalesScript struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProductID uint      `gorm:"index;not null" json:"product_id"`
	Scenario  string    `gorm:"size:255" json:"scenario"`
	Content   string    `gorm:"type:text" json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (SalesScript) TableName() string {
	return "sales_scripts"
}
