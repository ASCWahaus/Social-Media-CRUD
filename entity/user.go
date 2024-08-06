package entity

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"not null"`
	Age       uint      `json:"age" gorm:"not null"`
	Email     string    `json:"email" gorm:"unique;type:varchar(255)"`
	Username  string    `json:"username" gorm:"not null"`
	Password  string    `gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at"`
	Photos    []Photo
}
