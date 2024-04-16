package core

import (
	"github.com/ErdemOzgen/blackdagger/service/core/scheduler"
	"go.uber.org/fx"
)

func NewScheduler(topLevelModule fx.Option) *fx.App {
	return fx.New(
		topLevelModule,
		scheduler.Module,
		fx.Invoke(scheduler.LifetimeHooks),
		fx.NopLogger,
	)
}
