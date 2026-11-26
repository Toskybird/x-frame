package service

import (
	"context"
	"fmt"

	"github.com/x-frame/x-frame/internal/pkg/pagination"
)

// AdminService provides admin user management for the base scaffold.
type AdminService interface {
	ListUsers(ctx context.Context, page, pageSize int, filters UserListFilters) ([]User, int64, error)
	GetUser(ctx context.Context, id int64) (*User, error)
	CreateUser(ctx context.Context, input *CreateUserInput) (*User, error)
	UpdateUser(ctx context.Context, id int64, input *UpdateUserInput) (*User, error)
	DeleteUser(ctx context.Context, id int64) error
}

type CreateUserInput struct {
	Email       string
	Password    string
	Username    string
	Notes       string
	Balance     float64
	Concurrency int
}

type UpdateUserInput struct {
	Email       string
	Password    string
	Username    *string
	Notes       *string
	Balance     *float64
	Concurrency *int
	Status      string
}

type adminServiceImpl struct {
	userRepo UserRepository
}

// NewAdminService creates a scaffold AdminService (user CRUD only).
func NewAdminService(userRepo UserRepository) AdminService {
	return &adminServiceImpl{userRepo: userRepo}
}

func (s *adminServiceImpl) ListUsers(ctx context.Context, page, pageSize int, filters UserListFilters) ([]User, int64, error) {
	params := pagination.PaginationParams{Page: page, PageSize: pageSize}
	users, result, err := s.userRepo.ListWithFilters(ctx, params, filters)
	if err != nil {
		return nil, 0, err
	}
	return users, result.Total, nil
}

func (s *adminServiceImpl) GetUser(ctx context.Context, id int64) (*User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *adminServiceImpl) CreateUser(ctx context.Context, input *CreateUserInput) (*User, error) {
	if input == nil {
		return nil, fmt.Errorf("nil input")
	}
	user := &User{
		Email:        input.Email,
		Username:     input.Username,
		Notes:        input.Notes,
		Role:         RoleUser,
		Balance:      input.Balance,
		Concurrency:  input.Concurrency,
		Status:       StatusActive,
	}
	if err := user.SetPassword(input.Password); err != nil {
		return nil, err
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *adminServiceImpl) UpdateUser(ctx context.Context, id int64, input *UpdateUserInput) (*User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if input.Email != "" {
		user.Email = input.Email
	}
	if input.Password != "" {
		if err := user.SetPassword(input.Password); err != nil {
			return nil, err
		}
	}
	if input.Username != nil {
		user.Username = *input.Username
	}
	if input.Notes != nil {
		user.Notes = *input.Notes
	}
	if input.Balance != nil {
		user.Balance = *input.Balance
	}
	if input.Concurrency != nil {
		user.Concurrency = *input.Concurrency
	}
	if input.Status != "" {
		user.Status = input.Status
	}
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *adminServiceImpl) DeleteUser(ctx context.Context, id int64) error {
	return s.userRepo.Delete(ctx, id)
}
