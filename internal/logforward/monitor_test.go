package logforward

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMonitorHandlerHealthAndMetrics(t *testing.T) {
	sink := NewAsyncSink(&flakySink{}, AsyncOptions{QueueSize: 8})
	defer func() { _ = sink.Close(context.Background()) }()

	err := sink.Forward(context.Background(), Record{Line: "line-1"})
	require.NoError(t, err)
	require.NoError(t, sink.Close(context.Background()))

	srv := httptest.NewServer(NewMonitorHandler(sink, "/log-forwarding"))
	defer srv.Close()

	healthRes, err := http.Get(srv.URL + "/log-forwarding/health")
	require.NoError(t, err)
	defer healthRes.Body.Close()
	require.Equal(t, http.StatusOK, healthRes.StatusCode)

	var health map[string]any
	require.NoError(t, json.NewDecoder(healthRes.Body).Decode(&health))
	require.Equal(t, "closed", health["status"])

	metricsRes, err := http.Get(srv.URL + "/log-forwarding/metrics")
	require.NoError(t, err)
	defer metricsRes.Body.Close()
	require.Equal(t, http.StatusOK, metricsRes.StatusCode)

	var metrics AsyncSnapshot
	require.NoError(t, json.NewDecoder(metricsRes.Body).Decode(&metrics))
	require.Equal(t, uint64(1), metrics.Queued)
	require.Equal(t, uint64(1), metrics.Forwarded)
}

func TestMonitorHandlerPrometheusMetrics(t *testing.T) {
	sink := NewAsyncSink(&flakySink{}, AsyncOptions{QueueSize: 8})
	defer func() { _ = sink.Close(context.Background()) }()

	err := sink.Forward(context.Background(), Record{Line: "line-1"})
	require.NoError(t, err)
	require.NoError(t, sink.Close(context.Background()))

	srv := httptest.NewServer(NewMonitorHandler(sink, "/log-forwarding"))
	defer srv.Close()

	promRes, err := http.Get(srv.URL + "/log-forwarding/metrics/prometheus")
	require.NoError(t, err)
	defer promRes.Body.Close()
	require.Equal(t, http.StatusOK, promRes.StatusCode)
	require.Contains(t, promRes.Header.Get("Content-Type"), "text/plain")

	body, err := io.ReadAll(promRes.Body)
	require.NoError(t, err)
	text := string(body)
	require.Contains(t, text, "blackdagger_log_forwarding_queue_depth")
	require.Contains(t, text, "blackdagger_log_forwarding_records_total{result=\"queued\"} 1")
	require.Contains(t, text, "blackdagger_log_forwarding_records_total{result=\"forwarded\"} 1")
}

func TestMonitorHandlerMetricsContentNegotiation(t *testing.T) {
	sink := NewAsyncSink(&flakySink{}, AsyncOptions{QueueSize: 8})
	defer func() { _ = sink.Close(context.Background()) }()

	err := sink.Forward(context.Background(), Record{Line: "line-1"})
	require.NoError(t, err)
	require.NoError(t, sink.Close(context.Background()))

	srv := httptest.NewServer(NewMonitorHandler(sink, "/log-forwarding"))
	defer srv.Close()

	req, err := http.NewRequest(http.MethodGet, srv.URL+"/log-forwarding/metrics?format=prometheus", nil)
	require.NoError(t, err)
	req.Header.Set("Accept", "text/plain")

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Contains(t, res.Header.Get("Content-Type"), "text/plain")

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.True(t, strings.Contains(string(body), "blackdagger_log_forwarding_records_total"))
}
