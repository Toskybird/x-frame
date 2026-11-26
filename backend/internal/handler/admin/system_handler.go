package admin

import (
	"github.com/x-frame/x-frame/internal/pkg/response"
	"github.com/x-frame/x-frame/internal/service"

	"github.com/gin-gonic/gin"
)

// SystemHandler handles system-related operations
type SystemHandler struct {
	lockSvc *service.SystemOperationLockService
	version string
}

// NewSystemHandler creates a new SystemHandler
func NewSystemHandler(lockSvc *service.SystemOperationLockService, version string) *SystemHandler {
	return &SystemHandler{
		lockSvc: lockSvc,
		version: version,
	}
}

// GetVersion returns the current version
// GET /api/v1/admin/system/version
func (h *SystemHandler) GetVersion(c *gin.Context) {
	_ = h.lockSvc
	response.Success(c, gin.H{
		"version": h.version,
	})
}
