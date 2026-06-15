package executor

import (
	"context"
	"testing"

	"github.com/ErdemOzgen/blackdagger/internal/dag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnsibleExecutor(t *testing.T) {
	t.Run("BuildsArgsAndDefaults", func(t *testing.T) {
		step := dag.Step{
			Command: "site.yml",
			Args:    []string{"--syntax-check"},
			ExecutorConfig: dag.ExecutorConfig{
				Type: "ansible",
				Config: map[string]any{
					"inventory": "inventory.ini",
					"extraVars": map[string]any{
						"env": "prod",
					},
					"tags":   []string{"deploy", "app"},
					"check":  true,
					"diff":   true,
					"become": true,
					"forks":  10,
				},
			},
		}

		exec, err := newAnsible(context.Background(), step)
		require.NoError(t, err)

		ans, ok := exec.(*ansibleExecutor)
		require.True(t, ok)
		assert.Equal(t, "ansible-playbook", ans.cfg.Binary)
		assert.Equal(t, "site.yml", ans.args[0])
		assert.Contains(t, ans.args, "--check")
		assert.Contains(t, ans.args, "--diff")
		assert.Contains(t, ans.args, "--become")
		assert.Contains(t, ans.args, "--forks")
		assert.Contains(t, ans.args, "10")
		assert.Contains(t, ans.args, "-e")
		assert.Contains(t, ans.args, "env=prod")
		assert.Contains(t, ans.args, "--syntax-check")
	})

	t.Run("RequiresPlaybook", func(t *testing.T) {
		step := dag.Step{
			ExecutorConfig: dag.ExecutorConfig{Type: "ansible", Config: map[string]any{}},
		}
		_, err := newAnsible(context.Background(), step)
		require.Error(t, err)
		require.ErrorIs(t, err, errAnsiblePlaybookRequired)
	})
}
