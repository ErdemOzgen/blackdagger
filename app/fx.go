package app

import (
	"os"

	"github.com/ErdemOzgen/blackdagger/internal/config"
	"github.com/ErdemOzgen/blackdagger/internal/engine"
	"github.com/ErdemOzgen/blackdagger/internal/logger"
	"github.com/ErdemOzgen/blackdagger/internal/persistence/client"
	"github.com/ErdemOzgen/blackdagger/service/frontend"
	"go.uber.org/fx"
)

var (
	TopLevelModule = fx.Options(
		fx.Provide(ConfigProvider),
		fx.Provide(engine.NewFactory),
		fx.Provide(logger.NewSlogLogger),
		fx.Provide(client.NewDataStoreFactory),
	)
	cfgInstance *config.Config
)

func ConfigProvider() *config.Config {
	if cfgInstance != nil {
		return cfgInstance
	}
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	if err := config.LoadConfig(home); err != nil {
		panic(err)
	}
	cfgInstance = config.Get()
	return cfgInstance
}

func NewFrontendService() *fx.App {
	return fx.New(
		TopLevelModule,
		frontend.Module,
		fx.Invoke(frontend.LifetimeHooks),
		fx.NopLogger,
	)
}
