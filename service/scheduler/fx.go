package scheduler

import (
	"context"

	"github.com/ErdemOzgen/blackdagger/internal/config"
	"github.com/ErdemOzgen/blackdagger/internal/engine"
	blackdaggerlogger "github.com/ErdemOzgen/blackdagger/internal/logger"
	"github.com/ErdemOzgen/blackdagger/service/scheduler/entry_reader"
	"github.com/ErdemOzgen/blackdagger/service/scheduler/scheduler"
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
	Logger      blackdaggerlogger.Logger
	EntryReader scheduler.EntryReader
}

func EntryReaderProvider(
	cfg *config.Config,
	engineFactory engine.Factory,
	jf entry_reader.JobFactory,
	logger blackdaggerlogger.Logger,
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
		WorkDir:       cfg.WorkDir,
		EngineFactory: engineFactory,
		Executable:    cfg.Executable,
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
