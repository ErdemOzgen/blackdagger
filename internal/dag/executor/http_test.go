package executor

import (
	"context"
	"testing"

	"github.com/ErdemOzgen/blackdagger/internal/dag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPExecutor(t *testing.T) {
	t.Run("URLRequired", func(t *testing.T) {
		step := dag.Step{Command: "GET"}

		_, err := newHTTP(context.Background(), step)
		require.Error(t, err)
		require.ErrorIs(t, err, errHTTPURLRequired)
	})

	t.Run("ExpandsMethodURLAndQueryFromEnv", func(t *testing.T) {
		t.Setenv("HTTP_TEST_METHOD", "post")
		t.Setenv("HTTP_TEST_URL", "https://example.com/hook")
		t.Setenv("HTTP_TEST_QUERY", "deploy")

		step := dag.Step{
			Command: "${HTTP_TEST_METHOD}",
			Args:    []string{"${HTTP_TEST_URL}"},
			ExecutorConfig: dag.ExecutorConfig{
				Type: "http",
				Config: map[string]any{
					"query": map[string]any{
						"task": "${HTTP_TEST_QUERY}",
					},
				},
			},
		}

		exec, err := newHTTP(context.Background(), step)
		require.NoError(t, err)

		h, ok := exec.(*http)
		require.True(t, ok)
		assert.Equal(t, "post", h.method)
		assert.Equal(t, "https://example.com/hook", h.url)
		assert.Equal(t, "deploy", h.cfg.Query["task"])
	})
}
