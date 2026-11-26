package admin

import (
	"time"

	"github.com/x-frame/x-frame/internal/pkg/response"
	"github.com/x-frame/x-frame/internal/service"

	"github.com/gin-gonic/gin"
)

// DashboardHandler handles admin dashboard statistics
type DashboardHandler struct {
	dashboardService *service.DashboardService
	startTime        time.Time
}

// NewDashboardHandler creates a new admin dashboard handler
func NewDashboardHandler(dashboardService *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
		startTime:        time.Now(),
	}
}

// GetStats handles getting dashboard statistics
// GET /api/v1/admin/dashboard/stats
func (h *DashboardHandler) GetStats(c *gin.Context) {
	stats, err := h.dashboardService.GetDashboardStats(c.Request.Context())
	if err != nil {
		response.Error(c, 500, "Failed to get dashboard statistics")
		return
	}
	uptime := int64(time.Since(h.startTime).Seconds())
	response.Success(c, gin.H{
		"total_users":  stats.TotalUsers,
		"active_users": stats.ActiveUsers,
		"admin_users":  stats.AdminUsers,
		"generated_at": stats.GeneratedAt,
		"uptime":       uptime,
	})
}
