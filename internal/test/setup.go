package test

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/ErdemOzgen/blackdagger/internal/client"
	"github.com/ErdemOzgen/blackdagger/internal/config"
	"github.com/ErdemOzgen/blackdagger/internal/logger"
	"github.com/ErdemOzgen/blackdagger/internal/persistence"
	dsclient "github.com/ErdemOzgen/blackdagger/internal/persistence/client"
	"github.com/ErdemOzgen/blackdagger/internal/util"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

type Setup struct {
	Config *config.Config
	Logger logger.Logger

	homeDir string
}

func (t Setup) Cleanup() {
	_ = os.RemoveAll(t.homeDir)
}

func (t Setup) DataStore() persistence.DataStores {
	return dsclient.NewDataStores(
		t.Config.DAGs,
		t.Config.DataDir,
		t.Config.SuspendFlagsDir,
		dsclient.DataStoreOptions{
			LatestStatusToday: t.Config.LatestStatusToday,
		},
	)
}

func (t Setup) Client() client.Client {
	return client.New(
		t.DataStore(), t.Config.Executable, t.Config.WorkDir, logger.Default,
	)
}

var (
	lock sync.Mutex
)

func SetupTest(t *testing.T) Setup {
	lock.Lock()
	defer lock.Unlock()

	tmpDir := util.MustTempDir("blackdagger_test")
	err := os.Setenv("HOME", tmpDir)
	require.NoError(t, err)

	configDir := filepath.Join(tmpDir, "config")
	viper.AddConfigPath(configDir)
	viper.SetConfigType("yaml")
	viper.SetConfigName("admin")

	cfg, err := config.Load()
	require.NoError(t, err)

	cfg.DAGs = filepath.Join(tmpDir, "dags")
	cfg.WorkDir = tmpDir
	cfg.BaseConfig = filepath.Join(tmpDir, "config", "base.yaml")
	cfg.DataDir = filepath.Join(tmpDir, "data")
	cfg.LogDir = filepath.Join(tmpDir, "log")
	cfg.AdminLogsDir = filepath.Join(tmpDir, "log", "admin")

	// Set the executable path to the test binary.
	cfg.Executable = filepath.Join(util.MustGetwd(), "../../bin/blackdagger")

	// Set environment variables.
	// This is required for some tests that run the executable
	_ = os.Setenv("BLACKDAGGER_DAGS_DIR", cfg.DAGs)
	_ = os.Setenv("BLACKDAGGER_WORK_DIR", cfg.WorkDir)
	_ = os.Setenv("BLACKDAGGER_BASE_CONFIG", cfg.BaseConfig)
	_ = os.Setenv("BLACKDAGGER_LOG_DIR", cfg.LogDir)
	_ = os.Setenv("BLACKDAGGER_DATA_DIR", cfg.DataDir)
	_ = os.Setenv("BLACKDAGGER_SUSPEND_FLAGS_DIR", cfg.SuspendFlagsDir)
	_ = os.Setenv("BLACKDAGGER_ADMIN_LOG_DIR", cfg.AdminLogsDir)

	return Setup{
		Config: cfg,
		Logger: NewLogger(),

		homeDir: tmpDir,
	}
}

func SetupForDir(t *testing.T, dir string) Setup {
	lock.Lock()
	defer lock.Unlock()

	tmpDir := util.MustTempDir("blackdagger_test")
	err := os.Setenv("HOME", tmpDir)
	require.NoError(t, err)

	viper.AddConfigPath(dir)
	viper.SetConfigType("yaml")
	viper.SetConfigName("admin")

	cfg, err := config.Load()
	require.NoError(t, err)

	return Setup{
		Config: cfg,
		Logger: NewLogger(),

		homeDir: tmpDir,
	}
}

func NewLogger() logger.Logger {
	return logger.NewLogger(logger.NewLoggerArgs{
		Debug:  true,
		Format: "text",
	})
}
