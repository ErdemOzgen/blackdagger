package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ErdemOzgen/blackdagger/internal/dag"
	"github.com/ErdemOzgen/blackdagger/internal/scheduler"
	"github.com/ErdemOzgen/blackdagger/internal/utils"
)

type StatusResponse struct {
	Status *Status `json:"status"`
}

type Pid int

const PidNotRunning Pid = -1

func (p Pid) String() string {
	if p == PidNotRunning {
		return ""
	}
	return fmt.Sprintf("%d", p)
}

func (p Pid) IsRunning() bool {
	return p != PidNotRunning
}

type Status struct {
	RequestId  string                    `json:"RequestId"`
	Name       string                    `json:"Name"`
	Status     scheduler.SchedulerStatus `json:"Status"`
	StatusText string                    `json:"StatusText"`
	Pid        Pid                       `json:"Pid"`
	Nodes      []*Node                   `json:"Nodes"`
	OnExit     *Node                     `json:"OnExit"`
	OnSuccess  *Node                     `json:"OnSuccess"`
	OnFailure  *Node                     `json:"OnFailure"`
	OnCancel   *Node                     `json:"OnCancel"`
	StartedAt  string                    `json:"StartedAt"`
	FinishedAt string                    `json:"FinishedAt"`
	Log        string                    `json:"Log"`
	Params     string                    `json:"Params"`
}

type StatusFile struct {
	File   string
	Status *Status
}

func StatusFromJson(s string) (*Status, error) {
	status := &Status{}
	err := json.Unmarshal([]byte(s), status)
	if err != nil {
		return nil, err
	}
	return status, err
}

func NewStatusDefault(d *dag.DAG) *Status {
	return NewStatus(d, nil, scheduler.SchedulerStatus_None, int(PidNotRunning), nil, nil)
}

func NewStatus(
	d *dag.DAG,
	nodes []*scheduler.Node,
	status scheduler.SchedulerStatus,
	pid int,
	startTime, endTime *time.Time,
) *Status {
	var onExit, onSuccess, onFailure, onCancel *Node
	onExit = nodeOrNil(d.HandlerOn.Exit)
	onSuccess = nodeOrNil(d.HandlerOn.Success)
	onFailure = nodeOrNil(d.HandlerOn.Failure)
	onCancel = nodeOrNil(d.HandlerOn.Cancel)
	return &Status{
		Name:       d.Name,
		Status:     status,
		StatusText: status.String(),
		Pid:        Pid(pid),
		Nodes:      nodesOrSteps(nodes, d.Steps),
		OnExit:     onExit,
		OnSuccess:  onSuccess,
		OnFailure:  onFailure,
		OnCancel:   onCancel,
		StartedAt:  formatTime(startTime),
		FinishedAt: formatTime(endTime),
		Params:     strings.Join(d.Params, " "),
	}
}

func nodesOrSteps(nodes []*scheduler.Node, steps []*dag.Step) []*Node {
	if len(nodes) != 0 {
		return FromNodes(nodes)
	}
	return FromSteps(steps)
}

func formatTime(val *time.Time) string {
	if val == nil {
		return ""
	}
	return utils.FormatTime(*val)
}

func (st *Status) CorrectRunningStatus() {
	if st.Status == scheduler.SchedulerStatus_Running {
		st.Status = scheduler.SchedulerStatus_Error
		st.StatusText = st.Status.String()
	}
}

func (st *Status) ToJson() ([]byte, error) {
	js, err := json.Marshal(st)
	if err != nil {
		return []byte{}, err
	}
	return js, nil
}
