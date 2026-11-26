package dto

import "time"

type User struct {
	ID          int64     `json:"id"`
	Email       string    `json:"email"`
	Username    string    `json:"username"`
	Role        string    `json:"role"`
	Balance     float64   `json:"balance"`
	Concurrency int       `json:"concurrency"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AdminUser 是管理员接口使用的 user DTO。
type AdminUser struct {
	User
	Notes string `json:"notes"`
}
