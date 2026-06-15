package logforward

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

type MonitorOptions struct {
	Host     string
	Port     int
	BasePath string
}

type Monitor struct {
	addr   string
	server *http.Server
}

func StartMonitor(sink *AsyncSink, opts MonitorOptions) (*Monitor, error) {
	if sink == nil {
		return nil, fmt.Errorf("log forwarding monitor requires async sink")
	}
	if strings.TrimSpace(opts.Host) == "" {
		opts.Host = "127.0.0.1"
	}
	if opts.Port < 0 {
		return nil, fmt.Errorf("log forwarding monitor port must be zero or greater")
	}
	if strings.TrimSpace(opts.BasePath) == "" {
		opts.BasePath = "/log-forwarding"
	}

	addr := fmt.Sprintf("%s:%d", opts.Host, opts.Port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	server := &http.Server{
		Handler:           NewMonitorHandler(sink, opts.BasePath),
		ReadHeaderTimeout: 5 * time.Second,
	}

	m := &Monitor{addr: ln.Addr().String(), server: server}
	go func() {
		_ = server.Serve(ln)
	}()

	return m, nil
}

func (m *Monitor) Address() string {
	if m == nil {
		return ""
	}
	return m.addr
}

func (m *Monitor) Shutdown(ctx context.Context) error {
	if m == nil || m.server == nil {
		return nil
	}
	return m.server.Shutdown(ctx)
}

func NewMonitorHandler(sink *AsyncSink, basePath string) http.Handler {
	mux := http.NewServeMux()
	base := normalizeMonitorBasePath(basePath)

	mux.HandleFunc(base+"/health", func(w http.ResponseWriter, _ *http.Request) {
		snapshot := sink.Snapshot()
		status := "ok"
		if snapshot.Closed {
			status = "closed"
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"status":        status,
			"queueDepth":    snapshot.QueueDepth,
			"queueCapacity": snapshot.QueueCapacity,
			"dropped":       snapshot.Dropped,
			"failed":        snapshot.Failed,
		})
	})

	mux.HandleFunc(base+"/metrics/prometheus", func(w http.ResponseWriter, _ *http.Request) {
		snapshot := sink.Snapshot()
		writePrometheus(w, http.StatusOK, snapshot)
	})

	mux.HandleFunc(base+"/metrics", func(w http.ResponseWriter, r *http.Request) {
		snapshot := sink.Snapshot()
		if wantsPrometheus(r) {
			writePrometheus(w, http.StatusOK, snapshot)
			return
		}
		writeJSON(w, http.StatusOK, snapshot)
	})

	return mux
}

func normalizeMonitorBasePath(basePath string) string {
	base := strings.TrimSpace(basePath)
	if base == "" {
		return "/log-forwarding"
	}
	if !strings.HasPrefix(base, "/") {
		base = "/" + base
	}
	return strings.TrimRight(base, "/")
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writePrometheus(w http.ResponseWriter, status int, snapshot AsyncSnapshot) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	w.WriteHeader(status)

	_, _ = fmt.Fprintf(w,
		"# HELP blackdagger_log_forwarding_queue_depth Current queued log records waiting to be forwarded.\n"+
			"# TYPE blackdagger_log_forwarding_queue_depth gauge\n"+
			"blackdagger_log_forwarding_queue_depth %d\n"+
			"# HELP blackdagger_log_forwarding_queue_capacity Maximum queue capacity for buffered log forwarding.\n"+
			"# TYPE blackdagger_log_forwarding_queue_capacity gauge\n"+
			"blackdagger_log_forwarding_queue_capacity %d\n"+
			"# HELP blackdagger_log_forwarding_sink_closed Whether the async forwarding sink has been closed (1 closed, 0 open).\n"+
			"# TYPE blackdagger_log_forwarding_sink_closed gauge\n"+
			"blackdagger_log_forwarding_sink_closed %d\n"+
			"# HELP blackdagger_log_forwarding_records_total Total log forwarding records by result.\n"+
			"# TYPE blackdagger_log_forwarding_records_total counter\n"+
			"blackdagger_log_forwarding_records_total{result=\"queued\"} %d\n"+
			"blackdagger_log_forwarding_records_total{result=\"forwarded\"} %d\n"+
			"blackdagger_log_forwarding_records_total{result=\"dropped\"} %d\n"+
			"blackdagger_log_forwarding_records_total{result=\"failed\"} %d\n"+
			"blackdagger_log_forwarding_records_total{result=\"retried\"} %d\n",
		snapshot.QueueDepth,
		snapshot.QueueCapacity,
		boolToProm(snapshot.Closed),
		snapshot.Queued,
		snapshot.Forwarded,
		snapshot.Dropped,
		snapshot.Failed,
		snapshot.Retried,
	)
}

func wantsPrometheus(r *http.Request) bool {
	if r == nil {
		return false
	}
	format := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("format")))
	if format == "prometheus" || format == "prom" {
		return true
	}
	accept := strings.ToLower(r.Header.Get("Accept"))
	return strings.Contains(accept, "application/openmetrics-text") ||
		strings.Contains(accept, "text/plain")
}

func boolToProm(v bool) int {
	if v {
		return 1
	}
	return 0
}
