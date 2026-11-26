//go:build wireinject
// +build wireinject

package main

import (
	"net/http"

	"github.com/x-frame/x-frame/internal/config"
	"github.com/x-frame/x-frame/internal/handler"
	"github.com/x-frame/x-frame/internal/repository"
	"github.com/x-frame/x-frame/internal/server"
	"github.com/x-frame/x-frame/internal/server/middleware"
	"github.com/x-frame/x-frame/internal/service"

	"github.com/google/wire"
)

func initializeApplication(buildInfo handler.BuildInfo) (*Application, error) {
	wire.Build(
		config.ProviderSet,
		repository.ProviderSet,
		service.ProviderSet,
		middleware.ProviderSet,
		handler.ProviderSet,
		server.ProviderSet,
		provideServiceBuildInfo,
		provideCleanup,
		wire.Struct(new(Application), "Server", "Cleanup"),
	)
	return nil, nil
}

// silence unused import in wireinject mode when regenerating
var _ = http.Server{}
