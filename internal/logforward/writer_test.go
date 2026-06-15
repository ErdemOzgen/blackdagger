package logforward

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type memSink struct {
	mu      sync.Mutex
	records []Record
}

func (m *memSink) Forward(_ context.Context, rec Record) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.records = append(m.records, rec)
	return nil
}

func TestWriter(t *testing.T) {
	sink := &memSink{}
	w := NewWriter(context.Background(), sink, "req-1", "dag-1", "step-1")

	_, _ = w.Write([]byte("line1\nline2\npartial"))
	w.Flush()

	assert.Len(t, sink.records, 3)
	assert.Equal(t, "line1", sink.records[0].Line)
	assert.Equal(t, "line2", sink.records[1].Line)
	assert.Equal(t, "partial", sink.records[2].Line)
	assert.Equal(t, "req-1", sink.records[0].RequestID)
	assert.Equal(t, "dag-1", sink.records[0].DAGName)
	assert.Equal(t, "step-1", sink.records[0].StepName)
}
