package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ErdemOzgen/blackdagger/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolver(t *testing.T) {
	t.Run("App home directory", func(t *testing.T) {
		tmpDir := util.MustTempDir("test")
		defer os.RemoveAll(tmpDir)

		os.Setenv("TEST_APP_HOME", filepath.Join(tmpDir, "blackdagger"))
		r := newResolver("TEST_APP_HOME", filepath.Join(tmpDir, ".blackdagger"), XDGConfig{})

		assert.Equal(t, r, resolver{
			configDir:       filepath.Join(tmpDir, "blackdagger"),
			dagsDir:         filepath.Join(tmpDir, "blackdagger", "dags"),
			suspendFlagsDir: filepath.Join(tmpDir, "blackdagger", "suspend"),
			dataDir:         filepath.Join(tmpDir, "blackdagger", "data"),
			logsDir:         filepath.Join(tmpDir, "blackdagger", "logs"),
			adminLogsDir:    filepath.Join(tmpDir, "blackdagger", "logs/admin"),
			baseConfigFile:  filepath.Join(tmpDir, "blackdagger", "base.yaml"),
		})
	})
	t.Run("Legacy home directory", func(t *testing.T) {
		tmpDir := util.MustTempDir("test")
		defer os.RemoveAll(tmpDir)

		legacyPath := filepath.Join(tmpDir, ".blackdagger")
		err := os.MkdirAll(legacyPath, os.ModePerm)
		require.NoError(t, err)

		r := newResolver("UNSET_APP_HOME", legacyPath, XDGConfig{})

		assert.Equal(t, r, resolver{
			configDir:       filepath.Join(tmpDir, ".blackdagger"),
			dagsDir:         filepath.Join(tmpDir, ".blackdagger", "dags"),
			suspendFlagsDir: filepath.Join(tmpDir, ".blackdagger", "suspend"),
			dataDir:         filepath.Join(tmpDir, ".blackdagger", "data"),
			logsDir:         filepath.Join(tmpDir, ".blackdagger", "logs"),
			adminLogsDir:    filepath.Join(tmpDir, ".blackdagger", "logs", "admin"),
			baseConfigFile:  filepath.Join(tmpDir, ".blackdagger", "base.yaml"),
		})
	})
	t.Run("XDG_CONFIG_HOME", func(t *testing.T) {
		r := newResolver("UNSET_APP_HOME", ".test", XDGConfig{
			DataHome:   "/home/user/.local/share",
			ConfigHome: "/home/user/.config",
		})
		assert.Equal(t, r, resolver{
			configDir:       "/home/user/.config/blackdagger",
			dagsDir:         "/home/user/.config/blackdagger/dags",
			suspendFlagsDir: "/home/user/.local/share/blackdagger/suspend",
			dataDir:         "/home/user/.local/share/blackdagger/history",
			logsDir:         "/home/user/.local/share/blackdagger/logs",
			adminLogsDir:    "/home/user/.local/share/blackdagger/logs/admin",
			baseConfigFile:  "/home/user/.config/blackdagger/base.yaml",
		})
	})
}
