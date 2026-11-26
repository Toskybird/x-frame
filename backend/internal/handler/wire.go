package handler

import (
	"github.com/x-frame/x-frame/internal/handler/admin"
	"github.com/x-frame/x-frame/internal/service"

	"github.com/google/wire"
)

// ProvideAdminHandlers creates the AdminHandlers struct
func ProvideAdminHandlers(
	dashboardHandler *admin.DashboardHandler,
	userHandler *admin.UserHandler,
	announcementHandler *admin.AnnouncementHandler,
	backupHandler *admin.BackupHandler,
	settingHandler *admin.SettingHandler,
	systemHandler *admin.SystemHandler,
	userAttributeHandler *admin.UserAttributeHandler,
) *AdminHandlers {
	return &AdminHandlers{
		Dashboard:     dashboardHandler,
		User:          userHandler,
		Announcement:  announcementHandler,
		Backup:        backupHandler,
		Setting:       settingHandler,
		System:        systemHandler,
		UserAttribute: userAttributeHandler,
	}
}

// ProvideSystemHandler creates admin.SystemHandler without UpdateService
func ProvideSystemHandler(lockService *service.SystemOperationLockService, buildInfo BuildInfo) *admin.SystemHandler {
	return admin.NewSystemHandler(lockService, buildInfo.Version)
}

// ProvideSettingHandler creates SettingHandler with version from BuildInfo
func ProvideSettingHandler(settingService *service.SettingService, buildInfo BuildInfo) *SettingHandler {
	return NewSettingHandler(settingService, buildInfo.Version)
}

// ProvideHandlers creates the Handlers struct
func ProvideHandlers(
	authHandler *AuthHandler,
	userHandler *UserHandler,
	announcementHandler *AnnouncementHandler,
	adminHandlers *AdminHandlers,
	settingHandler *SettingHandler,
	totpHandler *TotpHandler,
	_ *service.IdempotencyCoordinator,
	_ *service.IdempotencyCleanupService,
) *Handlers {
	return &Handlers{
		Auth:         authHandler,
		User:         userHandler,
		Announcement: announcementHandler,
		Admin:        adminHandlers,
		Setting:      settingHandler,
		Totp:         totpHandler,
	}
}

// ProviderSet is the Wire provider set for scaffold handlers
var ProviderSet = wire.NewSet(
	NewAuthHandler,
	NewUserHandler,
	NewAnnouncementHandler,
	NewTotpHandler,
	ProvideSettingHandler,

	admin.NewDashboardHandler,
	admin.NewUserHandler,
	admin.NewAnnouncementHandler,
	admin.NewBackupHandler,
	admin.NewSettingHandler,
	ProvideSystemHandler,
	admin.NewUserAttributeHandler,

	ProvideAdminHandlers,
	ProvideHandlers,
)
