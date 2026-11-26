package handler

import (
	"github.com/x-frame/x-frame/internal/handler/admin"
)

// AdminHandlers contains admin-related HTTP handlers kept in the scaffold.
type AdminHandlers struct {
	Dashboard     *admin.DashboardHandler
	User          *admin.UserHandler
	Announcement  *admin.AnnouncementHandler
	Backup        *admin.BackupHandler
	Setting       *admin.SettingHandler
	System        *admin.SystemHandler
	UserAttribute *admin.UserAttributeHandler
}

// Handlers contains all HTTP handlers
type Handlers struct {
	Auth         *AuthHandler
	User         *UserHandler
	Announcement *AnnouncementHandler
	Admin        *AdminHandlers
	Setting      *SettingHandler
	Totp         *TotpHandler
}

// BuildInfo contains build-time information
type BuildInfo struct {
	Version   string
	BuildType string // "source" for manual builds, "release" for CI builds
}
