package config

import (
	"os"
	"path/filepath"

	"github.com/ErdemOzgen/blackdagger/internal/util"
)

type resolver struct {
	configDir       string
	dagsDir         string
	suspendFlagsDir string
	dataDir         string
	logsDir         string
	adminLogsDir    string
	baseConfigFile  string
}

type XDGConfig struct {
	DataHome   string
	ConfigHome string
}

func newResolver(appHomeEnv, legacyPath string, xdg XDGConfig) resolver {
	var (
		r           resolver
		useXDGRules = true
	)

	// For backward compatibility.
	// If the environment variable is set, use it.
	// Use the legacy ~/.blackdagger directory if it exists.
	if v := os.Getenv(appHomeEnv); v != "" {
		r.configDir = v
		useXDGRules = false
	} else if util.FileExists(legacyPath) {
		r.configDir = legacyPath
		useXDGRules = false
	} else {
		r.configDir = filepath.Join(xdg.ConfigHome, appName)
	}

	if useXDGRules {
		r.dataDir = filepath.Join(xdg.DataHome, appName, "history")
		r.logsDir = filepath.Join(xdg.DataHome, appName, "logs")
		r.baseConfigFile = filepath.Join(xdg.ConfigHome, appName, "base.yaml")
		r.adminLogsDir = filepath.Join(xdg.DataHome, appName, "logs", "admin")
		r.suspendFlagsDir = filepath.Join(xdg.DataHome, appName, "suspend")
		r.dagsDir = filepath.Join(xdg.ConfigHome, appName, "dags")
	} else {
		r.dataDir = filepath.Join(r.configDir, "data")
		r.logsDir = filepath.Join(r.configDir, "logs")
		r.baseConfigFile = filepath.Join(r.configDir, "base.yaml")
		r.adminLogsDir = filepath.Join(r.configDir, "logs", "admin")
		r.suspendFlagsDir = filepath.Join(r.configDir, "suspend")
		r.dagsDir = filepath.Join(r.configDir, "dags")
	}

	return r
}
