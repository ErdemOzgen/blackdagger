package dag

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Load(t *testing.T) {
	tests := []struct {
		name             string
		file             string
		expectedError    string
		expectedLocation string
	}{
		{
			name:             "WithExt",
			file:             filepath.Join(testdataDir, "loader_test.yaml"),
			expectedLocation: filepath.Join(testdataDir, "loader_test.yaml"),
		},
		{
			name:             "WithoutExt",
			file:             filepath.Join(testdataDir, "loader_test"),
			expectedLocation: filepath.Join(testdataDir, "loader_test.yaml"),
		},
		{
			name:          "InvalidPath",
			file:          filepath.Join(testdataDir, "not_existing_file.yaml"),
			expectedError: "no such file or directory",
		},
		{
			name:          "InvalidDAG",
			file:          filepath.Join(testdataDir, "err_decode.yaml"),
			expectedError: "has invalid keys: invalidkey",
		},
		{
			name:          "InvalidYAML",
			file:          filepath.Join(testdataDir, "err_parse.yaml"),
			expectedError: "cannot unmarshal",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dg, err := Load("", tt.file, "")
			if tt.expectedError != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedError)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedLocation, dg.Location)
			}
		})
	}
}

func Test_LoadMetadata(t *testing.T) {
	t.Run("Metadata", func(t *testing.T) {
		dg, err := LoadMetadata(filepath.Join(testdataDir, "default.yaml"))
		require.NoError(t, err)

		require.Equal(t, dg.Name, "default")
		// Check if steps are empty since we are loading metadata only
		require.True(t, len(dg.Steps) == 0)
	})
}

func Test_loadBaseConfig(t *testing.T) {
	t.Run("LoadBaseConfigFile", func(t *testing.T) {
		dg, err := loadBaseConfig(filepath.Join(testdataDir, "base.yaml"), buildOpts{})
		require.NotNil(t, dg)
		require.NoError(t, err)
	})
}

func Test_LoadDefaultConfig(t *testing.T) {
	t.Run("DefaultConfigWithoutBaseConfig", func(t *testing.T) {
		file := filepath.Join(testdataDir, "default.yaml")
		dg, err := Load("", file, "")

		require.NoError(t, err)

		// Check if the default values are set correctly
		assert.Equal(t, "", dg.LogDir)
		assert.Equal(t, file, dg.Location)
		assert.Equal(t, "default", dg.Name)
		assert.Equal(t, time.Second*60, dg.MaxCleanUpTime)
		assert.Equal(t, 30, dg.HistRetentionDays)

		// Check if the steps are loaded correctly
		require.Len(t, dg.Steps, 1)
		assert.Equal(t, "1", dg.Steps[0].Name, "1")
		assert.Equal(t, "true", dg.Steps[0].Command, "true")
		assert.Equal(t, filepath.Dir(file), dg.Steps[0].Dir)
	})
}

const (
	testDAG = `
name: test DAG
steps:
  - name: "1"
    command: "true"
`
)

func Test_LoadYAML(t *testing.T) {
	t.Run("ValidYAMLData", func(t *testing.T) {
		ret, err := loadYAML([]byte(testDAG), buildOpts{})
		require.NoError(t, err)
		require.Equal(t, ret.Name, "test DAG")

		step := ret.Steps[0]
		require.Equal(t, step.Name, "1")
		require.Equal(t, step.Command, "true")
	})
	t.Run("InvalidYAMLData", func(t *testing.T) {
		_, err := loadYAML([]byte(`invalidyaml`), buildOpts{})
		require.Error(t, err)
	})
	t.Run("ImportsRequireFilePath", func(t *testing.T) {
		_, err := loadYAML([]byte(`
imports:
  - child
steps:
  - name: "1"
    command: "true"
`), buildOpts{})
		require.Error(t, err)
		require.ErrorIs(t, err, errImportsRequirePath)
	})
}

func Test_Load_WithImports(t *testing.T) {
	t.Run("MergesImportedSteps", func(t *testing.T) {
		file := filepath.Join(testdataDir, "imports", "parent.yaml")

		dg, err := Load("", file, "")
		require.NoError(t, err)

		require.Len(t, dg.Steps, 2)
		assert.Equal(t, "imported-step", dg.Steps[0].Name)
		assert.Equal(t, "local-step", dg.Steps[1].Name)
	})

	t.Run("MergesNestedImports", func(t *testing.T) {
		file := filepath.Join(testdataDir, "imports", "nested_parent.yaml")

		dg, err := Load("", file, "")
		require.NoError(t, err)

		require.Len(t, dg.Steps, 3)
		assert.Equal(t, "nested-grandchild-step", dg.Steps[0].Name)
		assert.Equal(t, "nested-child-step", dg.Steps[1].Name)
		assert.Equal(t, "nested-local-step", dg.Steps[2].Name)
	})

	t.Run("DetectsImportCycle", func(t *testing.T) {
		file := filepath.Join(testdataDir, "imports", "cycle_a.yaml")

		_, err := Load("", file, "")
		require.Error(t, err)
		require.Contains(t, err.Error(), errCircularImport.Error())
	})

	t.Run("DetectsDuplicateImportedStepName", func(t *testing.T) {
		file := filepath.Join(testdataDir, "imports", "duplicate_parent.yaml")

		_, err := Load("", file, "")
		require.Error(t, err)
		require.ErrorIs(t, err, errDuplicateStepName)
	})

	t.Run("BlocksPathTraversalOutsideRoot", func(t *testing.T) {
		rootDir := t.TempDir()
		importsDir := filepath.Join(rootDir, "imports")
		outsideDir := filepath.Join(rootDir, "outside")
		require.NoError(t, os.MkdirAll(importsDir, 0o755))
		require.NoError(t, os.MkdirAll(outsideDir, 0o755))

		mainFile := filepath.Join(importsDir, "main.yaml")
		require.NoError(t, os.WriteFile(mainFile, []byte(`
imports:
  - ../outside/shared
steps:
  - name: main
    command: echo main
`), 0o644))

		outsideFile := filepath.Join(outsideDir, "shared.yaml")
		require.NoError(t, os.WriteFile(outsideFile, []byte(`
steps:
  - name: outside
    command: echo outside
`), 0o644))

		_, err := Load("", mainFile, "")
		require.Error(t, err)
		require.ErrorIs(t, err, errImportPathOutsideRoot)
	})

	t.Run("BlocksSymlinkEscapeOutsideRoot", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("symlink test is skipped on windows")
		}

		rootDir := t.TempDir()
		importsDir := filepath.Join(rootDir, "imports")
		outsideDir := filepath.Join(rootDir, "outside")
		require.NoError(t, os.MkdirAll(importsDir, 0o755))
		require.NoError(t, os.MkdirAll(outsideDir, 0o755))

		outsideFile := filepath.Join(outsideDir, "shared.yaml")
		require.NoError(t, os.WriteFile(outsideFile, []byte(`
steps:
  - name: outside
    command: echo outside
`), 0o644))

		linkPath := filepath.Join(importsDir, "linked.yaml")
		require.NoError(t, os.Symlink(outsideFile, linkPath))

		mainFile := filepath.Join(importsDir, "main.yaml")
		require.NoError(t, os.WriteFile(mainFile, []byte(`
imports:
  - ./linked
steps:
  - name: main
    command: echo main
`), 0o644))

		_, err := Load("", mainFile, "")
		require.Error(t, err)
		require.ErrorIs(t, err, errImportPathOutsideRoot)
	})
}
