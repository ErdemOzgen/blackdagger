package executor

import (
	"bytes"
	"context"
	"io"
	nethttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/ErdemOzgen/blackdagger/internal/dag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebhookExecutor(t *testing.T) {
	t.Run("ConfigDecodeAndEnvExpansion", func(t *testing.T) {
		t.Setenv("WEBHOOK_TEST_URL", "https://example.com/hook")
		t.Setenv("WEBHOOK_TEST_TOKEN", "token-value")

		step := dag.Step{
			ExecutorConfig: dag.ExecutorConfig{
				Type: "webhook",
				Config: map[string]any{
					"url":    "${WEBHOOK_TEST_URL}",
					"method": "post",
					"headers": map[string]any{
						"Authorization": "Bearer ${WEBHOOK_TEST_TOKEN}",
					},
					"body": "hello",
				},
			},
		}

		exec, err := newWebhook(context.Background(), step)
		require.NoError(t, err)

		wh, ok := exec.(*webhook)
		require.True(t, ok)
		assert.Equal(t, "https://example.com/hook", wh.cfg.URL)
		assert.Equal(t, "POST", wh.cfg.Method)
		assert.Equal(t, "Bearer token-value", wh.cfg.Headers["Authorization"])
		assert.Equal(t, "hello", wh.cfg.Body)
	})

	t.Run("RunSuccess", func(t *testing.T) {
		t.Parallel()

		var gotBody, gotMethod, gotToken, gotQuery string
		srv := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			gotBody = string(body)
			gotMethod = r.Method
			gotToken = r.Header.Get("X-Token")
			gotQuery = r.URL.Query().Get("task")
			w.WriteHeader(nethttp.StatusOK)
			_, _ = w.Write([]byte("ok"))
		}))
		defer srv.Close()

		step := dag.Step{
			ExecutorConfig: dag.ExecutorConfig{
				Type: "webhook",
				Config: map[string]any{
					"url":    srv.URL,
					"method": "POST",
					"headers": map[string]any{
						"X-Token": "abc",
					},
					"query": map[string]any{
						"task": "deploy",
					},
					"body": "payload",
				},
			},
		}

		exec, err := newWebhook(context.Background(), step)
		require.NoError(t, err)

		buf := new(bytes.Buffer)
		exec.SetStdout(buf)
		require.NoError(t, exec.Run())
		assert.Equal(t, "POST", gotMethod)
		assert.Equal(t, "abc", gotToken)
		assert.Equal(t, "deploy", gotQuery)
		assert.Equal(t, "payload", gotBody)
		assert.Contains(t, buf.String(), "200 OK")
		assert.Contains(t, buf.String(), "ok")
	})

	t.Run("RunSilentOnlyWritesBody", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, _ *nethttp.Request) {
			w.WriteHeader(nethttp.StatusOK)
			_, _ = w.Write([]byte("silent-body"))
		}))
		defer srv.Close()

		step := dag.Step{
			ExecutorConfig: dag.ExecutorConfig{
				Type: "webhook",
				Config: map[string]any{
					"url":    srv.URL,
					"silent": true,
				},
			},
		}

		exec, err := newWebhook(context.Background(), step)
		require.NoError(t, err)

		buf := new(bytes.Buffer)
		exec.SetStdout(buf)
		require.NoError(t, exec.Run())
		assert.Equal(t, "silent-body", buf.String())
	})

	t.Run("RunWithCustomSuccessCodes", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, _ *nethttp.Request) {
			w.WriteHeader(nethttp.StatusAccepted)
			_, _ = w.Write([]byte("accepted"))
		}))
		defer srv.Close()

		step := dag.Step{
			ExecutorConfig: dag.ExecutorConfig{
				Type: "webhook",
				Config: map[string]any{
					"url":                srv.URL,
					"successStatusCodes": []int{202},
				},
			},
		}

		exec, err := newWebhook(context.Background(), step)
		require.NoError(t, err)

		buf := new(bytes.Buffer)
		exec.SetStdout(buf)
		require.NoError(t, exec.Run())
	})

	t.Run("RunReturnsErrorOnFailureStatus", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, _ *nethttp.Request) {
			w.WriteHeader(nethttp.StatusInternalServerError)
			_, _ = w.Write([]byte("failed"))
		}))
		defer srv.Close()

		step := dag.Step{
			ExecutorConfig: dag.ExecutorConfig{
				Type: "webhook",
				Config: map[string]any{
					"url": srv.URL,
				},
			},
		}

		exec, err := newWebhook(context.Background(), step)
		require.NoError(t, err)

		buf := new(bytes.Buffer)
		exec.SetStdout(buf)

		err = exec.Run()
		require.Error(t, err)
		require.ErrorIs(t, err, errWebhookStatusCode)
	})

	t.Run("URLRequired", func(t *testing.T) {
		execStep := dag.Step{ExecutorConfig: dag.ExecutorConfig{Type: "webhook"}}
		_, err := newWebhook(context.Background(), execStep)
		require.Error(t, err)
		require.ErrorIs(t, err, errWebhookURLRequired)
	})
}
