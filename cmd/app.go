package cmd

import (
	"github.com/ErdemOzgen/blackdagger/internal/config"
	"github.com/ErdemOzgen/blackdagger/internal/engine"
	"github.com/ErdemOzgen/blackdagger/internal/logger"
	"github.com/ErdemOzgen/blackdagger/internal/persistence/client"
	"github.com/ErdemOzgen/blackdagger/service/frontend"
	"go.uber.org/fx"
)

var (
	topLevelModule = fx.Options(
		fx.Provide(config.Get),
		fx.Provide(engine.NewFactory),
		fx.Provide(logger.NewSlogLogger),
		fx.Provide(client.NewDataStoreFactory),
	)
)

func newFrontend() *fx.App {
	return fx.New(
		topLevelModule,
		frontend.Module,
		fx.Invoke(frontend.LifetimeHooks),
		fx.NopLogger,
	)
}
