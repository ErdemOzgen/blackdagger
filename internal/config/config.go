package config

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/ErdemOzgen/blackdagger/internal/utils"
	"github.com/spf13/viper"
)

type Config struct {
	Host               string
	Port               int
	DAGs               string
	Command            string
	WorkDir            string
	IsBasicAuth        bool
	BasicAuthUsername  string
	BasicAuthPassword  string
	LogEncodingCharset string
	LogDir             string
	DataDir            string
	SuspendFlagsDir    string
	AdminLogsDir       string
	BaseConfig         string
	NavbarColor        string
	NavbarTitle        string
	Env                sync.Map
	TLS                *TLS
	IsAuthToken        bool
	AuthToken          string
}

type TLS struct {
	CertFile string
	KeyFile  string
}

var (
	instance *Config
	mu       sync.RWMutex
	isLoaded atomic.Bool
)

func Get() *Config {
	if !isLoaded.Load() {
		home, _ := os.UserHomeDir()
		if err := LoadConfig(home); err != nil {
			panic(err)
		}
	}
	return getConfig()
}

func getConfig() *Config {
	mu.Lock()
	defer mu.Unlock()
	return instance
}

func setConfig(cfg *Config) {
	mu.Lock()
	defer mu.Unlock()
	isLoaded.Swap(true)
	instance = cfg
}

var (
	loadConfigLock sync.Mutex
)

func LoadConfig(userHomeDir string) error {
	loadConfigLock.Lock()
	defer loadConfigLock.Unlock()

	appHome := appHomeDir(userHomeDir)

	viper.SetEnvPrefix("blackdagger")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	_ = viper.BindEnv("command", "BLACKDAGGER_EXECUTABLE")
	_ = viper.BindEnv("dags", "BLACKDAGGER_DAGS_DIR")
	_ = viper.BindEnv("workDir", "BLACKDAGGER_WORK_DIR")
	_ = viper.BindEnv("isBasicAuth", "BLACKDAGGER_IS_BASICAUTH")
	_ = viper.BindEnv("basicAuthUsername", "BLACKDAGGER_BASICAUTH_USERNAME")
	_ = viper.BindEnv("basicAuthPassword", "BLACKDAGGER_BASICAUTH_PASSWORD")
	_ = viper.BindEnv("logEncodingCharset", "BLACKDAGGER_LOG_ENCODING_CHARSET")
	_ = viper.BindEnv("baseConfig", "BLACKDAGGER_BASE_CONFIG")
	_ = viper.BindEnv("logDir", "BLACKDAGGER_LOG_DIR")
	_ = viper.BindEnv("dataDir", "BLACKDAGGER_DATA_DIR")
	_ = viper.BindEnv("suspendFlagsDir", "BLACKDAGGER_SUSPEND_FLAGS_DIR")
	_ = viper.BindEnv("adminLogsDir", "BLACKDAGGER_ADMIN_LOG_DIR")
	_ = viper.BindEnv("navbarColor", "BLACKDAGGER_NAVBAR_COLOR")
	_ = viper.BindEnv("navbarTitle", "BLACKDAGGER_NAVBAR_TITLE")
	_ = viper.BindEnv("tls.certFile", "BLACKDAGGER_CERT_FILE")
	_ = viper.BindEnv("tls.keyFile", "BLACKDAGGER_KEY_FILE")
	_ = viper.BindEnv("isAuthToken", "BLACKDAGGER_IS_AUTHTOKEN")
	_ = viper.BindEnv("authToken", "BLACKDAGGER_AUTHTOKEN")
	command := "blackdagger"
	if ex, err := os.Executable(); err == nil {
		command = ex
	}

	viper.SetDefault("host", "127.0.0.1")
	viper.SetDefault("port", "8080")
	viper.SetDefault("command", command)
	viper.SetDefault("dags", path.Join(appHome, "dags"))
	viper.SetDefault("workDir", "")
	viper.SetDefault("isBasicAuth", "0")
	viper.SetDefault("basicAuthUsername", "")
	viper.SetDefault("basicAuthPassword", "")
	viper.SetDefault("logEncodingCharset", "")
	viper.SetDefault("baseConfig", path.Join(appHome, "config.yaml"))
	viper.SetDefault("logDir", path.Join(appHome, "logs"))
	viper.SetDefault("dataDir", path.Join(appHome, "data"))
	viper.SetDefault("suspendFlagsDir", path.Join(appHome, "suspend"))
	viper.SetDefault("adminLogsDir", path.Join(appHome, "logs", "admin"))
	viper.SetDefault("navbarColor", "")
	viper.SetDefault("navbarTitle", "Dagu")
	viper.SetDefault("isAuthToken", "0")
	viper.SetDefault("authToken", "0")

	viper.AutomaticEnv()

	_ = viper.ReadInConfig()

	cfg := &Config{}
	err := viper.Unmarshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to unmarshal cfg file: %w", err)
	}
	loadLegacyEnvs(cfg)
	loadEnvs(cfg)

	setConfig(cfg)

	return nil
}

func (cfg *Config) GetAPIBaseURL() string {
	return "/api/v1"
}

func loadEnvs(cfg *Config) {
	cfg.Env.Range(func(k, v interface{}) bool {
		_ = os.Setenv(k.(string), v.(string))
		return true
	})
	for k, v := range utils.DefaultEnv() {
		if _, ok := cfg.Env.Load(k); !ok {
			cfg.Env.Store(k, v)
		}
	}
}

func loadLegacyEnvs(cfg *Config) {
	// For backward compatibility.
	cfg.NavbarColor = getEnv("BLACKDAGGER__ADMIN_NAVBAR_COLOR", cfg.NavbarColor)
	cfg.NavbarTitle = getEnv("BLACKDAGGER__ADMIN_NAVBAR_TITLE", cfg.NavbarTitle)
	cfg.Port = getEnvI("BLACKDAGGER__ADMIN_PORT", cfg.Port)
	cfg.Host = getEnv("BLACKDAGGER__ADMIN_HOST", cfg.Host)
	cfg.DataDir = getEnv("BLACKDAGGER__DATA", cfg.DataDir)
	cfg.LogDir = getEnv("BLACKDAGGER__DATA", cfg.LogDir)
	cfg.SuspendFlagsDir = getEnv("BLACKDAGGER__SUSPEND_FLAGS_DIR", cfg.SuspendFlagsDir)
	cfg.BaseConfig = getEnv("BLACKDAGGER__SUSPEND_FLAGS_DIR", cfg.BaseConfig)
	cfg.AdminLogsDir = getEnv("BLACKDAGGER__ADMIN_LOGS_DIR", cfg.AdminLogsDir)
}

func getEnv(env, def string) string {
	v := os.Getenv(env)
	if v == "" {
		return def
	}
	return v
}

func parseInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func getEnvI(env string, def int) int {
	v := os.Getenv(env)
	if v == "" {
		return def
	}
	return parseInt(v)
}

const (
	appHomeEnv     = "BLACKDAGGER_HOME"
	appHomeDefault = ".blackdagger"
)

func appHomeDir(userHomeDir string) string {
	appDir := os.Getenv(appHomeEnv)
	if appDir == "" {
		return path.Join(userHomeDir, appHomeDefault)
	}
	return appDir
}
