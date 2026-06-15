package logforward

import "context"

// Sink receives log lines associated with a DAG execution and step.
type Sink interface {
	Forward(ctx context.Context, rec Record) error
}

// Closer can be implemented by sinks that need graceful shutdown.
type Closer interface {
	Close(ctx context.Context) error
}

// Record is a normalized log event forwarded to external systems.
type Record struct {
	RequestID string `json:"requestId"`
	DAGName   string `json:"dagName"`
	StepName  string `json:"stepName"`
	Line      string `json:"line"`
}
