package scheduler

import (
	"context"
	"github.com/ErdemOzgen/blackdagger/internal/config"
	"github.com/ErdemOzgen/blackdagger/internal/engine"
	"github.com/ErdemOzgen/blackdagger/internal/logger"
	"github.com/ErdemOzgen/blackdagger/service/core/scheduler/entry_reader"
	"github.com/ErdemOzgen/blackdagger/service/core/scheduler/scheduler"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(EntryReaderProvider),
	fx.Provide(JobFactoryProvider),
	fx.Provide(New),
)

type Params struct {
	fx.In

	Config      *config.Config
	Logger      logger.Logger
	EntryReader scheduler.EntryReader
}

func EntryReaderProvider(
	cfg *config.Config,
	engineFactory engine.Factory,
	jf entry_reader.JobFactory,
	logger logger.Logger,
) scheduler.EntryReader {
	return entry_reader.New(entry_reader.Params{
		EngineFactory: engineFactory,
		// TODO: fix this
		DagsDir:    cfg.DAGs,
		JobFactory: jf,
		Logger:     logger,
	})
}

func JobFactoryProvider(cfg *config.Config, engineFactory engine.Factory) entry_reader.JobFactory {
	return &jobFactory{
		Command:       cfg.Command,
		WorkDir:       cfg.WorkDir,
		EngineFactory: engineFactory,
	}
}

func New(params Params) *scheduler.Scheduler {
	return scheduler.New(scheduler.Params{
		EntryReader: params.EntryReader,
		Logger:      params.Logger,
		// TODO: check this is used
		LogDir: params.Config.LogDir,
	})
}

func LifetimeHooks(lc fx.Lifecycle, a *scheduler.Scheduler) {
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) (err error) {
				return a.Start()
			},
			OnStop: func(_ context.Context) error {
				a.Stop()
				return nil
			},
		},
	)
}
