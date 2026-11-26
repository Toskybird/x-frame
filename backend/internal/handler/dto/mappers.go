package dto

import "github.com/x-frame/x-frame/internal/service"

func UserFromServiceShallow(u *service.User) *User {
	if u == nil {
		return nil
	}
	return &User{
		ID:          u.ID,
		Email:       u.Email,
		Username:    u.Username,
		Role:        u.Role,
		Balance:     u.Balance,
		Concurrency: u.Concurrency,
		Status:      u.Status,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

func UserFromService(u *service.User) *User {
	return UserFromServiceShallow(u)
}

// UserFromServiceAdmin converts a service User to DTO for admin users.
func UserFromServiceAdmin(u *service.User) *AdminUser {
	if u == nil {
		return nil
	}
	base := UserFromService(u)
	if base == nil {
		return nil
	}
	return &AdminUser{
		User:  *base,
		Notes: u.Notes,
	}
}
