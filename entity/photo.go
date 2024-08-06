package entity

import "time"

type Photo struct {
	ID        uint   `json:"id"`
	Title     string `json:"name" gorm:"not null, type:varchar(255)"`
	Caption   string `json:"caption" gorm:"not null, type:varchar(255)"`
	PhotoURL  string `json:"photo_url" gorm:"not null"`
	UserID    uint
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User      `json:"User"`
}
