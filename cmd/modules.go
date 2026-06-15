package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ErdemOzgen/blackdagger/internal/client"
	"github.com/ErdemOzgen/blackdagger/internal/config"
	"github.com/ErdemOzgen/blackdagger/internal/logforward"
	"github.com/ErdemOzgen/blackdagger/internal/logger"
	"github.com/ErdemOzgen/blackdagger/internal/persistence"
	dsclient "github.com/ErdemOzgen/blackdagger/internal/persistence/client"
)

func newClient(cfg *config.Config, ds persistence.DataStores, lg logger.Logger) client.Client {
	return client.New(ds, cfg.Executable, cfg.WorkDir, lg)
}

func newDataStores(cfg *config.Config) persistence.DataStores {
	return dsclient.NewDataStores(
		cfg.DAGs,
		cfg.DataDir,
		cfg.SuspendFlagsDir,
		dsclient.DataStoreOptions{
			LatestStatusToday: cfg.LatestStatusToday,
		},
	)
}

type managedLogForwardingSink struct {
	sink    logforward.Sink
	monitor *logforward.Monitor
}

func (m *managedLogForwardingSink) Forward(ctx context.Context, rec logforward.Record) error {
	if m == nil || m.sink == nil {
		return nil
	}
	return m.sink.Forward(ctx, rec)
}

func (m *managedLogForwardingSink) Close(ctx context.Context) error {
	if m == nil {
		return nil
	}

	var closeErr error
	if closer, ok := m.sink.(logforward.Closer); ok {
		if err := closer.Close(ctx); err != nil {
			closeErr = err
		}
	}
	if m.monitor != nil {
		if err := m.monitor.Shutdown(ctx); err != nil && closeErr == nil {
			closeErr = err
		}
	}

	return closeErr
}

func newLogForwarding(cfg *config.Config, lg logger.Logger) (logforward.Sink, bool, error) {
	if cfg.LogForwarding == nil || !cfg.LogForwarding.Enabled {
		return nil, false, nil
	}

	sinkType := strings.ToLower(strings.TrimSpace(cfg.LogForwarding.SinkType))
	if sinkType == "" {
		sinkType = "http"
	}

	switch sinkType {
	case "http":
		url := strings.TrimSpace(cfg.LogForwarding.HTTPURL)
		if url == "" {
			return nil, false, fmt.Errorf("log forwarding httpURL is required")
		}
		timeout := time.Duration(cfg.LogForwarding.TimeoutSec) * time.Second
		httpSink := logforward.NewHTTPSink(
			url,
			timeout,
			cfg.LogForwarding.Headers,
		)
		asyncOptions := logforward.AsyncOptions{
			QueueSize:      cfg.LogForwarding.QueueSize,
			MaxRetries:     cfg.LogForwarding.MaxRetries,
			InitialBackoff: time.Duration(cfg.LogForwarding.InitialBackoffMS) * time.Millisecond,
			MaxBackoff:     time.Duration(cfg.LogForwarding.MaxBackoffMS) * time.Millisecond,
		}
		if lg != nil {
			asyncOptions.Logger = lg.WithGroup("log_forwarding")
		}
		asyncSink := logforward.NewAsyncSink(httpSink, asyncOptions)

		var sink logforward.Sink = asyncSink
		if cfg.LogForwarding.MonitorEnabled {
			monitor, err := logforward.StartMonitor(asyncSink, logforward.MonitorOptions{
				Host:     cfg.LogForwarding.MonitorHost,
				Port:     cfg.LogForwarding.MonitorPort,
				BasePath: cfg.LogForwarding.MonitorBasePath,
			})
			if err != nil {
				return nil, false, fmt.Errorf("failed to start log forwarding monitor: %w", err)
			}
			if lg != nil {
				lg.Info(
					"Log forwarding monitor started",
					"address", monitor.Address(),
					"basePath", cfg.LogForwarding.MonitorBasePath,
				)
			}
			sink = &managedLogForwardingSink{sink: asyncSink, monitor: monitor}
		}

		return sink, cfg.LogForwarding.IncludeStepOutput, nil
	default:
		return nil, false, fmt.Errorf("unsupported log forwarding sink type: %s", sinkType)
	}
}
