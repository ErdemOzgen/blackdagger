package local

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

const validDAGSpec = `name: test-dag
steps:
  - name: step1
    command: echo hello
`

func TestDAGStoreGetSpecWithYMLExtension(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	file := filepath.Join(dir, "example.yml")
	require.NoError(t, os.WriteFile(file, []byte(validDAGSpec), 0644))

	store := NewDAGStore(&NewDAGStoreArgs{Dir: dir})

	spec, err := store.GetSpec("example")
	require.NoError(t, err)
	require.Equal(t, validDAGSpec, spec)
}

func TestDAGStoreWithSymlinkedRootDirectory(t *testing.T) {
	t.Parallel()

	if runtime.GOOS == "windows" {
		t.Skip("symlink tests are not stable on windows CI")
	}

	realDagsDir := t.TempDir()
	require.NoError(t, os.WriteFile(
		filepath.Join(realDagsDir, "job.yaml"),
		[]byte(validDAGSpec),
		0644,
	))

	linkParent := t.TempDir()
	symlinkedDagsDir := filepath.Join(linkParent, "dags")
	require.NoError(t, os.Symlink(realDagsDir, symlinkedDagsDir))

	store := NewDAGStore(&NewDAGStoreArgs{Dir: symlinkedDagsDir})

	dags, errs, err := store.List()
	require.NoError(t, err)
	require.Empty(t, errs)
	require.Len(t, dags, 1)
	require.Equal(t, "test-dag", dags[0].Name)

	spec, err := store.GetSpec("job")
	require.NoError(t, err)
	require.Equal(t, validDAGSpec, spec)
}

func TestDAGStoreUpdateSpecWithSymlinkedFile(t *testing.T) {
	t.Parallel()

	if runtime.GOOS == "windows" {
		t.Skip("symlink tests are not stable on windows CI")
	}

	dagsDir := t.TempDir()
	targetDir := t.TempDir()
	targetFile := filepath.Join(targetDir, "job.yaml")
	require.NoError(t, os.WriteFile(targetFile, []byte(validDAGSpec), 0644))

	symlinkedFile := filepath.Join(dagsDir, "job.yaml")
	require.NoError(t, os.Symlink(targetFile, symlinkedFile))

	store := NewDAGStore(&NewDAGStoreArgs{Dir: dagsDir})

	updatedSpec := `name: test-dag
steps:
  - name: step1
    command: echo updated
`
	require.NoError(t, store.UpdateSpec("job", []byte(updatedSpec)))

	content, err := os.ReadFile(targetFile)
	require.NoError(t, err)
	require.Equal(t, updatedSpec, string(content))
}
