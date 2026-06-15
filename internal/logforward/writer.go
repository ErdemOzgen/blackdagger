package logforward

import (
	"bytes"
	"context"
	"strings"
	"sync"
)

type Writer struct {
	sink      Sink
	requestID string
	dagName   string
	stepName  string
	ctx       context.Context
	mu        sync.Mutex
	pending   bytes.Buffer
}

func NewWriter(
	ctx context.Context,
	sink Sink,
	requestID, dagName, stepName string,
) *Writer {
	return &Writer{
		ctx:       ctx,
		sink:      sink,
		requestID: requestID,
		dagName:   dagName,
		stepName:  stepName,
	}
}

func (w *Writer) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	_, _ = w.pending.Write(p)
	for {
		data := w.pending.String()
		i := strings.IndexByte(data, '\n')
		if i < 0 {
			break
		}
		line := data[:i]
		w.pending.Reset()
		_, _ = w.pending.WriteString(data[i+1:])
		w.forwardLine(line)
	}

	return len(p), nil
}

func (w *Writer) Flush() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.pending.Len() == 0 {
		return
	}
	line := w.pending.String()
	w.pending.Reset()
	w.forwardLine(line)
}

func (w *Writer) forwardLine(line string) {
	line = strings.TrimSpace(line)
	if line == "" || w.sink == nil {
		return
	}
	_ = w.sink.Forward(w.ctx, Record{
		RequestID: w.requestID,
		DAGName:   w.dagName,
		StepName:  w.stepName,
		Line:      line,
	})
}
