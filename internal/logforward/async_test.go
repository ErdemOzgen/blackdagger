package logforward

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type flakySink struct {
	mu       sync.Mutex
	fails    int
	received int
}

type captureLogger struct {
	mu     sync.Mutex
	warns  []string
	errors []string
}

func (l *captureLogger) Warn(msg string, _ ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.warns = append(l.warns, msg)
}

func (l *captureLogger) Error(msg string, _ ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.errors = append(l.errors, msg)
}

type blockingSink struct {
	started chan struct{}
	release chan struct{}
}

func (b *blockingSink) Forward(_ context.Context, _ Record) error {
	select {
	case b.started <- struct{}{}:
	default:
	}
	<-b.release
	return nil
}

func (f *flakySink) Forward(_ context.Context, _ Record) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.received++
	if f.fails > 0 {
		f.fails--
		return errors.New("temporary failure")
	}
	return nil
}

func TestAsyncSinkRetryAndClose(t *testing.T) {
	sink := &flakySink{fails: 2}
	async := NewAsyncSink(sink, AsyncOptions{
		QueueSize:      8,
		MaxRetries:     3,
		InitialBackoff: 1 * time.Millisecond,
		MaxBackoff:     2 * time.Millisecond,
	})

	err := async.Forward(context.Background(), Record{Line: "hello"})
	require.NoError(t, err)
	require.NoError(t, async.Close(context.Background()))

	stats := async.Stats()
	require.Equal(t, uint64(1), stats.Queued)
	require.Equal(t, uint64(1), stats.Forwarded)
	require.Equal(t, uint64(2), stats.Retried)
	require.Equal(t, uint64(0), stats.Failed)
}

func TestAsyncSinkQueueFull(t *testing.T) {
	sink := &blockingSink{
		started: make(chan struct{}, 1),
		release: make(chan struct{}),
	}
	async := NewAsyncSink(sink, AsyncOptions{QueueSize: 1})
	defer func() { _ = async.Close(context.Background()) }()
	defer close(sink.release)

	err := async.Forward(context.Background(), Record{Line: "a"})
	require.NoError(t, err)
	<-sink.started

	err = async.Forward(context.Background(), Record{Line: "b"})
	require.NoError(t, err)

	err = async.Forward(context.Background(), Record{Line: "c"})
	require.ErrorIs(t, err, ErrQueueFull)
}

func TestAsyncSinkLogsRetryAndFailure(t *testing.T) {
	sink := &flakySink{fails: 10}
	logCapture := &captureLogger{}
	async := NewAsyncSink(sink, AsyncOptions{
		QueueSize:      8,
		MaxRetries:     1,
		InitialBackoff: 1 * time.Millisecond,
		MaxBackoff:     1 * time.Millisecond,
		Logger:         logCapture,
	})

	err := async.Forward(context.Background(), Record{Line: "hello"})
	require.NoError(t, err)
	require.NoError(t, async.Close(context.Background()))

	require.NotEmpty(t, logCapture.warns)
	require.NotEmpty(t, logCapture.errors)
}

func TestAsyncSinkCloseRejectsNewWrites(t *testing.T) {
	sink := &flakySink{}
	async := NewAsyncSink(sink, AsyncOptions{})
	require.NoError(t, async.Close(context.Background()))
	err := async.Forward(context.Background(), Record{Line: "late"})
	require.ErrorIs(t, err, ErrAsyncSinkClosed)
}
