package model

import (
	"errors"
	"fmt"

	"github.com/ErdemOzgen/blackdagger/internal/dag"
	"github.com/ErdemOzgen/blackdagger/internal/dag/scheduler"
	"github.com/ErdemOzgen/blackdagger/internal/util"
)

func FromSteps(steps []dag.Step) []*Node {
	var ret []*Node
	for _, s := range steps {
		ret = append(ret, NewNode(s))
	}
	return ret
}

func FromNodes(nodes []scheduler.NodeData) []*Node {
	var ret []*Node
	for _, node := range nodes {
		ret = append(ret, FromNode(node))
	}
	return ret
}

func FromNode(node scheduler.NodeData) *Node {
	return &Node{
		Step:       node.Step,
		Log:        node.State.Log,
		StartedAt:  util.FormatTime(node.State.StartedAt),
		FinishedAt: util.FormatTime(node.State.FinishedAt),
		Status:     node.State.Status,
		StatusText: node.State.Status.String(),
		RetryCount: node.State.RetryCount,
		DoneCount:  node.State.DoneCount,
		Error:      errText(node.State.Error),
	}
}

type Node struct {
	Step       dag.Step             `json:"Step"`
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
	startedAt, _ := util.ParseTime(n.StartedAt)
	finishedAt, _ := util.ParseTime(n.FinishedAt)
	return scheduler.NewNode(n.Step, scheduler.NodeState{
		Status:     n.Status,
		Log:        n.Log,
		StartedAt:  startedAt,
		FinishedAt: finishedAt,
		RetryCount: n.RetryCount,
		DoneCount:  n.DoneCount,
		Error:      errFromText(n.Error),
	})
}

func NewNode(step dag.Step) *Node {
	return &Node{
		Step:       step,
		StartedAt:  "-",
		FinishedAt: "-",
		Status:     scheduler.NodeStatusNone,
		StatusText: scheduler.NodeStatusNone.String(),
	}
}

var errNodeProcessing = errors.New("node processing error")

func errFromText(err string) error {
	if err == "" {
		return nil
	}
	return fmt.Errorf("%w: %s", errNodeProcessing, err)
}

func errText(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
