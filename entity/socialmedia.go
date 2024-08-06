package entity

import "time"

type SocialMedia struct {
	ID             uint      `json:"id"`
	Name           string    `json:"name" gorm:"not null, type:varchar(255)"`
	SocialMediaURL string    `json:"social_media_url" gorm:"not null"`
	UserID         uint      `json:"user_id" gorm:"not null"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	User           User      `json:"User"`
}
