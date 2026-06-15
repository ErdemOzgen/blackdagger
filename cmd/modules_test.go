package cmd

import (
	"context"
	"testing"

	"github.com/ErdemOzgen/blackdagger/internal/config"
	"github.com/ErdemOzgen/blackdagger/internal/logger"
	"github.com/stretchr/testify/require"
)

func TestNewLogForwarding(t *testing.T) {
	testLogger := logger.NewLogger(logger.NewLoggerArgs{Format: "text", Quiet: true})

	t.Run("Disabled", func(t *testing.T) {
		cfg := &config.Config{}
		sink, includeOutput, err := newLogForwarding(cfg, testLogger)
		require.NoError(t, err)
		require.Nil(t, sink)
		require.False(t, includeOutput)
	})

	t.Run("EnabledHTTPMissingURL", func(t *testing.T) {
		cfg := &config.Config{LogForwarding: &config.LogForwardingConfig{
			Enabled:  true,
			SinkType: "http",
		}}
		_, _, err := newLogForwarding(cfg, testLogger)
		require.Error(t, err)
	})

	t.Run("EnabledHTTPSuccess", func(t *testing.T) {
		cfg := &config.Config{LogForwarding: &config.LogForwardingConfig{
			Enabled:           true,
			SinkType:          "http",
			HTTPURL:           "https://logs.example.com/ingest",
			TimeoutSec:        3,
			IncludeStepOutput: true,
			QueueSize:         512,
			MaxRetries:        5,
			InitialBackoffMS:  50,
			MaxBackoffMS:      300,
			MonitorEnabled:    true,
			MonitorHost:       "127.0.0.1",
			MonitorPort:       0,
			MonitorBasePath:   "/log-forwarding",
		}}
		sink, includeOutput, err := newLogForwarding(cfg, testLogger)
		require.NoError(t, err)
		require.NotNil(t, sink)
		require.True(t, includeOutput)
		closer, ok := sink.(interface{ Close(context.Context) error })
		require.True(t, ok)
		require.NoError(t, closer.Close(context.Background()))
	})
}
