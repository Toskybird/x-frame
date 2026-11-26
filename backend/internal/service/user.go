package service

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int64
	Email        string
	Username     string
	Notes        string
	PasswordHash string
	Role         string
	Balance      float64
	Concurrency  int
	Status       string
	TokenVersion int64
	CreatedAt    time.Time
	UpdatedAt    time.Time

	// TOTP 双因素认证字段
	TotpSecretEncrypted *string
	TotpEnabled         bool
	TotpEnabledAt       *time.Time
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) IsActive() bool {
	return u.Status == StatusActive
}

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) == nil
}
