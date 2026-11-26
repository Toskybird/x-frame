package routes

import (
	"github.com/x-frame/x-frame/internal/handler"
	"github.com/x-frame/x-frame/internal/server/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterAdminRoutes 注册管理员路由（脚手架：dashboard/users/announcements/settings/system/backup）
func RegisterAdminRoutes(
	v1 *gin.RouterGroup,
	h *handler.Handlers,
	adminAuth middleware.AdminAuthMiddleware,
) {
	admin := v1.Group("/admin")
	admin.Use(gin.HandlerFunc(adminAuth))
	{
		registerDashboardRoutes(admin, h)
		registerUserManagementRoutes(admin, h)
		registerAnnouncementRoutes(admin, h)
		registerSettingsRoutes(admin, h)
		registerBackupRoutes(admin, h)
		registerSystemRoutes(admin, h)
		registerUserAttributeRoutes(admin, h)
	}
}

func registerDashboardRoutes(admin *gin.RouterGroup, h *handler.Handlers) {
	dashboard := admin.Group("/dashboard")
	{
		dashboard.GET("/stats", h.Admin.Dashboard.GetStats)
	}
}

func registerUserManagementRoutes(admin *gin.RouterGroup, h *handler.Handlers) {
	users := admin.Group("/users")
	{
		users.GET("", h.Admin.User.List)
		users.GET("/:id", h.Admin.User.GetByID)
		users.POST("", h.Admin.User.Create)
		users.PUT("/:id", h.Admin.User.Update)
		users.DELETE("/:id", h.Admin.User.Delete)

		users.GET("/:id/attributes", h.Admin.UserAttribute.GetUserAttributes)
		users.PUT("/:id/attributes", h.Admin.UserAttribute.UpdateUserAttributes)
	}
}

func registerAnnouncementRoutes(admin *gin.RouterGroup, h *handler.Handlers) {
	announcements := admin.Group("/announcements")
	{
		announcements.GET("", h.Admin.Announcement.List)
		announcements.GET("/:id", h.Admin.Announcement.GetByID)
		announcements.POST("", h.Admin.Announcement.Create)
		announcements.PUT("/:id", h.Admin.Announcement.Update)
		announcements.DELETE("/:id", h.Admin.Announcement.Delete)
	}
}

func registerSettingsRoutes(admin *gin.RouterGroup, h *handler.Handlers) {
	settings := admin.Group("/settings")
	{
		settings.GET("", h.Admin.Setting.GetSettings)
		settings.PUT("", h.Admin.Setting.UpdateSettings)
	}
}

func registerBackupRoutes(admin *gin.RouterGroup, h *handler.Handlers) {
	backup := admin.Group("/backup")
	{
		backup.GET("/s3-config", h.Admin.Backup.GetS3Config)
		backup.PUT("/s3-config", h.Admin.Backup.UpdateS3Config)
		backup.POST("/s3-config/test", h.Admin.Backup.TestS3Connection)
		backup.GET("/schedule", h.Admin.Backup.GetSchedule)
		backup.PUT("/schedule", h.Admin.Backup.UpdateSchedule)
		backup.GET("", h.Admin.Backup.ListBackups)
		backup.GET("/:id", h.Admin.Backup.GetBackup)
		backup.POST("", h.Admin.Backup.CreateBackup)
		backup.POST("/:id/restore", h.Admin.Backup.RestoreBackup)
		backup.DELETE("/:id", h.Admin.Backup.DeleteBackup)
		backup.GET("/:id/download-url", h.Admin.Backup.GetDownloadURL)
	}
}

func registerSystemRoutes(admin *gin.RouterGroup, h *handler.Handlers) {
	system := admin.Group("/system")
	{
		system.GET("/version", h.Admin.System.GetVersion)
	}
}

func registerUserAttributeRoutes(admin *gin.RouterGroup, h *handler.Handlers) {
	attrs := admin.Group("/user-attributes")
	{
		attrs.GET("/definitions", h.Admin.UserAttribute.ListDefinitions)
		attrs.POST("/definitions", h.Admin.UserAttribute.CreateDefinition)
		attrs.PUT("/definitions/:id", h.Admin.UserAttribute.UpdateDefinition)
		attrs.DELETE("/definitions/:id", h.Admin.UserAttribute.DeleteDefinition)
	}
}
