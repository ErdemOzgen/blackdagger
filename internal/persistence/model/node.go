package model

import (
	"fmt"

	"github.com/ErdemOzgen/blackdagger/internal/dag"
	"github.com/ErdemOzgen/blackdagger/internal/scheduler"
	"github.com/ErdemOzgen/blackdagger/internal/utils"
)

type Node struct {
	*dag.Step  `json:"Step"`
	Log        string               `json:"Log"`
	StartedAt  string               `json:"StartedAt"`
	FinishedAt string               `json:"FinishedAt"`
	Status     scheduler.NodeStatus `json:"Status"`
	RetryCount int                  `json:"RetryCount"`
	DoneCount  int                  `json:"DoneCount"`
	Error      string               `json:"Error"`
	StatusText string               `json:"StatusText"`
}

func (n *Node) ToNode() *scheduler.Node {
	startedAt, _ := utils.ParseTime(n.StartedAt)
	finishedAt, _ := utils.ParseTime(n.FinishedAt)
	return &scheduler.Node{
		Step: n.Step,
		NodeState: scheduler.NodeState{
			Status:     n.Status,
			Log:        n.Log,
			StartedAt:  startedAt,
			FinishedAt: finishedAt,
			RetryCount: n.RetryCount,
			DoneCount:  n.DoneCount,
			Error:      errFromText(n.Error),
		},
	}
}

func FromNode(n scheduler.NodeState, step *dag.Step) *Node {
	return &Node{
		Step:       step,
		Log:        n.Log,
		StartedAt:  utils.FormatTime(n.StartedAt),
		FinishedAt: utils.FormatTime(n.FinishedAt),
		Status:     n.Status,
		StatusText: n.Status.String(),
		RetryCount: n.RetryCount,
		DoneCount:  n.DoneCount,
		Error:      errText(n.Error),
	}
}

func errFromText(err string) error {
	if err == "" {
		return nil
	}
	return fmt.Errorf(err)
}

func errText(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func FromNodes(nodes []*scheduler.Node) []*Node {
	var ret []*Node
	for _, n := range nodes {
		ret = append(ret, FromNode(n.State(), n.Step))
	}
	return ret
}

func FromSteps(steps []*dag.Step) []*Node {
	var ret []*Node
	for _, s := range steps {
		ret = append(ret, nodeOrNil(s))
	}
	return ret
}

func NewNode(step *dag.Step) *Node {
	return &Node{
		Step:       step,
		StartedAt:  "-",
		FinishedAt: "-",
		Status:     scheduler.NodeStatus_None,
		StatusText: scheduler.NodeStatus_None.String(),
	}
}

func nodeOrNil(s *dag.Step) *Node {
	if s == nil {
		return nil
	}
	return NewNode(s)
}
