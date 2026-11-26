package service

import (
	"context"
	"fmt"
	"time"

	"github.com/x-frame/x-frame/internal/config"
	"github.com/x-frame/x-frame/internal/pkg/pagination"
)

// DashboardStats is a minimal scaffold dashboard payload (no usage_log dependency).
type DashboardStats struct {
	TotalUsers  int64 `json:"total_users"`
	ActiveUsers int64 `json:"active_users"`
	AdminUsers  int64 `json:"admin_users"`
	GeneratedAt int64 `json:"generated_at"`
}

// DashboardService provides admin dashboard statistics for the base scaffold.
type DashboardService struct {
	userRepo UserRepository
}

// NewDashboardService creates a minimal DashboardService.
func NewDashboardService(userRepo UserRepository, _ *config.Config) *DashboardService {
	return &DashboardService{userRepo: userRepo}
}

// GetDashboardStats returns basic user counts.
func (s *DashboardService) GetDashboardStats(ctx context.Context) (*DashboardStats, error) {
	_, totalRes, err := s.userRepo.List(ctx, pagination.PaginationParams{Page: 1, PageSize: 1})
	if err != nil {
		return nil, fmt.Errorf("count users: %w", err)
	}
	_, activeRes, err := s.userRepo.ListWithFilters(ctx, pagination.PaginationParams{Page: 1, PageSize: 1}, UserListFilters{Status: StatusActive})
	if err != nil {
		return nil, fmt.Errorf("count active users: %w", err)
	}
	_, adminRes, err := s.userRepo.ListWithFilters(ctx, pagination.PaginationParams{Page: 1, PageSize: 1}, UserListFilters{Role: RoleAdmin})
	if err != nil {
		return nil, fmt.Errorf("count admin users: %w", err)
	}

	return &DashboardStats{
		TotalUsers:  totalRes.Total,
		ActiveUsers: activeRes.Total,
		AdminUsers:  adminRes.Total,
		GeneratedAt: time.Now().Unix(),
	}, nil
}
