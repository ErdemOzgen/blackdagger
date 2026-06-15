package executor

import (
	"context"
	"testing"

	"github.com/ErdemOzgen/blackdagger/internal/dag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTerraformExecutor(t *testing.T) {
	t.Run("BuildsArgsAndDefaults", func(t *testing.T) {
		t.Setenv("TF_WORKDIR", "/tmp/tfdir")

		step := dag.Step{
			Command: "apply",
			Args:    []string{"-parallelism=2"},
			ExecutorConfig: dag.ExecutorConfig{
				Type: "terraform",
				Config: map[string]any{
					"workingDir":  "${TF_WORKDIR}",
					"init":        true,
					"initArgs":    []string{"-upgrade"},
					"varFiles":    []string{"env.tfvars"},
					"vars":        map[string]any{"region": "us-east-1"},
					"autoApprove": true,
				},
			},
		}

		exec, err := newTerraform(context.Background(), step)
		require.NoError(t, err)

		tf, ok := exec.(*terraformExecutor)
		require.True(t, ok)
		assert.Equal(t, "terraform", tf.cfg.Binary)
		assert.Equal(t, []string{"-chdir=/tmp/tfdir", "init", "-upgrade"}, tf.initArgs)
		assert.Contains(t, tf.mainArgs, "apply")
		assert.Contains(t, tf.mainArgs, "-auto-approve")
		assert.Contains(t, tf.mainArgs, "-var-file=env.tfvars")
		assert.Contains(t, tf.mainArgs, "-var")
		assert.Contains(t, tf.mainArgs, "region=us-east-1")
	})

	t.Run("RequiresSubcommand", func(t *testing.T) {
		step := dag.Step{
			ExecutorConfig: dag.ExecutorConfig{Type: "terraform", Config: map[string]any{}},
		}
		_, err := newTerraform(context.Background(), step)
		require.Error(t, err)
		require.ErrorIs(t, err, errTerraformSubcommandRequired)
	})
}
