package agent_test

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"path"
	"syscall"
	"testing"
	"time"

	"github.com/ErdemOzgen/blackdagger/internal/agent"
	"github.com/ErdemOzgen/blackdagger/internal/persistence"
	"github.com/ErdemOzgen/blackdagger/internal/persistence/client"

	"github.com/ErdemOzgen/blackdagger/internal/config"
	"github.com/ErdemOzgen/blackdagger/internal/dag"
	"github.com/ErdemOzgen/blackdagger/internal/engine"
	"github.com/ErdemOzgen/blackdagger/internal/persistence/model"
	"github.com/ErdemOzgen/blackdagger/internal/scheduler"
	"github.com/ErdemOzgen/blackdagger/internal/utils"
	"github.com/stretchr/testify/require"
)

var testdataDir = path.Join(utils.MustGetwd(), "testdata")

func setupTest(t *testing.T) (string, engine.Engine, persistence.DataStoreFactory) {
	t.Helper()

	tmpDir := utils.MustTempDir("blackdagger_test")
	_ = os.Setenv("HOME", tmpDir)
	_ = config.LoadConfig(tmpDir)

	ds := client.NewDataStoreFactory(&config.Config{
		DataDir: path.Join(tmpDir, ".blackdagger", "data"),
	})

	e := engine.NewFactory(ds, nil).Create()

	return tmpDir, e, ds
}

func TestRunDAG(t *testing.T) {
	tmpDir, e, df := setupTest(t)
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	d := testLoadDAG(t, "run.yaml")
	a := agent.New(&agent.Config{DAG: d}, e, df)

	status, _ := e.GetLatestStatus(d)
	require.Equal(t, scheduler.StatusNone, status.Status)

	go func() {
		err := a.Run(context.Background())
		require.NoError(t, err)
	}()

	time.Sleep(100 * time.Millisecond)

	require.Eventually(t, func() bool {
		status, err := e.GetLatestStatus(d)
		require.NoError(t, err)
		return status.Status == scheduler.StatusSuccess
	}, time.Second*2, time.Millisecond*100)

	// check deletion of expired history files
	d.HistRetentionDays = 0
	a = agent.New(&agent.Config{DAG: d}, e, df)
	err := a.Run(context.Background())
	require.NoError(t, err)
	statusList := e.GetRecentHistory(d, 100)
	require.Equal(t, 1, len(statusList))
}

func TestCheckRunning(t *testing.T) {
	tmpDir, e, df := setupTest(t)
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	d := testLoadDAG(t, "is_running.yaml")
	a := agent.New(&agent.Config{DAG: d}, e, df)

	go func() {
		_ = a.Run(context.Background())
	}()

	time.Sleep(time.Millisecond * 30)

	status := a.Status()
	require.NotNil(t, status)
	require.Equal(t, status.Status, scheduler.StatusRunning)

	a = agent.New(&agent.Config{DAG: d}, e, df)
	err := a.Run(context.Background())
	require.Error(t, err)
	require.Contains(t, err.Error(), "is already running")
}

func TestDryRun(t *testing.T) {
	tmpDir, e, df := setupTest(t)
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	d := testLoadDAG(t, "dry.yaml")
	a := agent.New(&agent.Config{DAG: d, Dry: true}, e, df)

	err := a.Run(context.Background())
	require.NoError(t, err)

	status := a.Status()
	require.NoError(t, err)

	require.Equal(t, scheduler.StatusSuccess, status.Status)
}

func TestCancelDAG(t *testing.T) {
	tmpDir, e, df := setupTest(t)
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	for _, abort := range []func(*agent.Agent){
		func(a *agent.Agent) { a.Signal(syscall.SIGTERM) },
	} {
		d := testLoadDAG(t, "sleep.yaml")
		a := agent.New(&agent.Config{DAG: d}, e, df)

		go func() {
			_ = a.Run(context.Background())
		}()

		time.Sleep(time.Millisecond * 100)
		abort(a)
		time.Sleep(time.Millisecond * 500)
		status, err := e.GetLatestStatus(d)
		require.NoError(t, err)
		require.Equal(t, scheduler.StatusCancel, status.Status)
	}
}

func TestPreConditionInvalid(t *testing.T) {
	tmpDir, e, df := setupTest(t)
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()
	d := testLoadDAG(t, "multiple_steps.yaml")

	d.Preconditions = []*dag.Condition{
		{
			Condition: "`echo 1`",
			Expected:  "0",
		},
	}

	a := agent.New(&agent.Config{DAG: d}, e, df)
	err := a.Run(context.Background())
	require.Error(t, err)

	status := a.Status()

	require.Equal(t, scheduler.StatusCancel, status.Status)
	require.Equal(t, scheduler.NodeStatusNone, status.Nodes[0].Status)
	require.Equal(t, scheduler.NodeStatusNone, status.Nodes[1].Status)
}

func TestPreConditionValid(t *testing.T) {
	tmpDir, e, df := setupTest(t)
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()
	d := testLoadDAG(t, "with_params.yaml")

	d.Preconditions = []*dag.Condition{
		{
			Condition: "`echo 1`",
			Expected:  "1",
		},
	}

	a := agent.New(&agent.Config{DAG: d}, e, df)
	err := a.Run(context.Background())
	require.NoError(t, err)

	status := a.Status()
	require.Equal(t, scheduler.StatusSuccess, status.Status)
	for _, s := range status.Nodes {
		require.Equal(t, scheduler.NodeStatusSuccess, s.Status)
	}
}

func TestStartError(t *testing.T) {
	tmpDir, e, df := setupTest(t)
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()
	d := testLoadDAG(t, "error.yaml")

	a := agent.New(&agent.Config{DAG: d}, e, df)
	err := a.Run(context.Background())
	require.Error(t, err)

	status := a.Status()
	require.Equal(t, scheduler.StatusError, status.Status)
}

func TestOnExit(t *testing.T) {
	tmpDir, e, df := setupTest(t)
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	d := testLoadDAG(t, "on_exit.yaml")
	a := agent.New(&agent.Config{DAG: d}, e, df)
	err := a.Run(context.Background())
	require.NoError(t, err)

	status := a.Status()
	require.Equal(t, scheduler.StatusSuccess, status.Status)
	for _, s := range status.Nodes {
		require.Equal(t, scheduler.NodeStatusSuccess, s.Status)
	}
	require.Equal(t, scheduler.NodeStatusSuccess, status.OnExit.Status)
}

func TestRetry(t *testing.T) {
	tmpDir, e, df := setupTest(t)
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	d := testLoadDAG(t, "retry.yaml")

	a := agent.New(&agent.Config{DAG: d}, e, df)
	err := a.Run(context.Background())
	require.Error(t, err)

	status := a.Status()
	require.Equal(t, scheduler.StatusError, status.Status)

	for _, n := range status.Nodes {
		n.CmdWithArgs = "true"
	}

	a = agent.New(&agent.Config{DAG: d, RetryTarget: status}, e, df)
	err = a.Run(context.Background())
	require.NoError(t, err)

	status = a.Status()
	require.Equal(t, scheduler.StatusSuccess, status.Status)

	for _, n := range status.Nodes {
		if n.Status != scheduler.NodeStatusSuccess &&
			n.Status != scheduler.NodeStatusSkipped {
			t.Errorf("invalid status: %s", n.Status.String())
		}
	}
}

func TestHandleHTTP(t *testing.T) {
	tmpDir, e, df := setupTest(t)
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	d := testLoadDAG(t, "handle_http.yaml")
	a := agent.New(&agent.Config{DAG: d}, e, df)

	go func() {
		err := a.Run(context.Background())
		require.NoError(t, err)
	}()

	timer := time.NewTimer(time.Millisecond * 50)
	defer timer.Stop()
	<-timer.C

	var mockResponseWriter = mockResponseWriter{}

	// status
	req := &http.Request{
		Method: "GET",
		URL: &url.URL{
			Path: "/status",
		},
	}

	a.HandleHTTP(&mockResponseWriter, req)
	require.Equal(t, http.StatusOK, mockResponseWriter.status)

	status, err := model.StatusFromJson(mockResponseWriter.body)
	require.NoError(t, err)
	require.Equal(t, scheduler.StatusRunning, status.Status)

	// invalid path
	req = &http.Request{
		Method: "GET",
		URL: &url.URL{
			Path: "/invalid-path",
		},
	}
	a.HandleHTTP(&mockResponseWriter, req)
	require.Equal(t, http.StatusNotFound, mockResponseWriter.status)

	// cancel
	req = &http.Request{
		Method: "POST",
		URL: &url.URL{
			Path: "/stop",
		},
	}
	a.HandleHTTP(&mockResponseWriter, req)
	require.Equal(t, http.StatusOK, mockResponseWriter.status)
	require.Equal(t, "OK", mockResponseWriter.body)

	timer2 := time.NewTimer(time.Millisecond * 50)
	defer timer2.Stop()
	<-timer2.C

	status = a.Status()
	require.Equal(t, status.Status, scheduler.StatusCancel)
}

type mockResponseWriter struct {
	status int
	body   string
	header *http.Header
}

var _ http.ResponseWriter = (*mockResponseWriter)(nil)

func (h *mockResponseWriter) Header() http.Header {
	if h.header == nil {
		h.header = &http.Header{}
	}
	return *h.header
}

func (h *mockResponseWriter) Write(body []byte) (int, error) {
	h.body = string(body)
	return 0, nil
}

func (h *mockResponseWriter) WriteHeader(statusCode int) {
	h.status = statusCode
}

func testLoadDAG(t *testing.T, name string) *dag.DAG {
	file := path.Join(testdataDir, name)
	cl := &dag.Loader{}
	d, err := cl.Load(file, "")
	require.NoError(t, err)
	return d
}
