package entity

import "time"

type Comment struct {
	ID        uint      `json:"id"`
	Message   string    `json:"message" gorm:"not null, type:varchar(255)"`
	PhotoID   uint      `json:"photo_id" gorm:"not null"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User      `json:"User"`
	Photo     Photo     `json:"Photo"`
}
