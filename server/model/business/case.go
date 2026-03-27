package business

import "time"

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

func (Case) TableName() string {
	return "cases"
}
