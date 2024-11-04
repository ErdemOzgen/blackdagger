package client

import (
	"path/filepath"

	"github.com/ErdemOzgen/blackdagger/internal/dag"
	"github.com/ErdemOzgen/blackdagger/internal/frontend/gen/restapi/operations/dags"
	"github.com/ErdemOzgen/blackdagger/internal/persistence"
	"github.com/ErdemOzgen/blackdagger/internal/persistence/model"
)

type Client interface {
	CreateDAG(id string) (string, error)
	GetDAGSpec(id string) (string, error)
	Grep(pattern string) ([]*persistence.GrepResult, []string, error)
	Rename(oldID, newID string) error
	Stop(workflow *dag.DAG) error
	StartAsync(workflow *dag.DAG, opts StartOptions)
	Start(workflow *dag.DAG, opts StartOptions) error
	Restart(workflow *dag.DAG, opts RestartOptions) error
	Retry(workflow *dag.DAG, requestID string) error
	GetCurrentStatus(workflow *dag.DAG) (*model.Status, error)
	GetStatusByRequestID(workflow *dag.DAG, requestID string) (*model.Status, error)
	GetLatestStatus(workflow *dag.DAG) (*model.Status, error)
	GetRecentHistory(workflow *dag.DAG, n int) []*model.StatusFile
	UpdateStatus(workflow *dag.DAG, status *model.Status) error
	UpdateDAG(id string, spec string) error
	DeleteDAG(id, loc string) error
	GetAllStatus() (statuses []*DAGStatus, errs []string, err error)
	GetAllStatusPagination(params dags.ListDagsParams) ([]*DAGStatus, *DagListPaginationSummaryResult, error)
	GetStatus(dagLocation string) (*DAGStatus, error)
	IsSuspended(id string) bool
	ToggleSuspend(id string, suspend bool) error
	GetTagList() ([]string, []string, error)
}

type StartOptions struct {
	Params string
	Quiet  bool
}

type RestartOptions struct {
	Quiet bool
}

type DAGStatus struct {
	File      string
	Dir       string
	DAG       *dag.DAG
	Status    *model.Status
	Suspended bool
	Error     error
	ErrorT    *string
}

type DagListPaginationSummaryResult struct {
	PageCount int
	ErrorList []string
}

func newDAGStatus(
	workflow *dag.DAG, s *model.Status, suspended bool, err error,
) *DAGStatus {
	ret := &DAGStatus{
		File:      filepath.Base(workflow.Location),
		Dir:       filepath.Dir(workflow.Location),
		DAG:       workflow,
		Status:    s,
		Suspended: suspended,
		Error:     err,
	}
	if err != nil {
		errT := err.Error()
		ret.ErrorT = &errT
	}
	return ret
}
