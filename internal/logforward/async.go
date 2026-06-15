package logforward

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrAsyncSinkClosed = errors.New("log forwarding sink is closed")
	ErrQueueFull       = errors.New("log forwarding queue is full")
)

type AsyncOptions struct {
	QueueSize      int
	MaxRetries     int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
	Logger         AsyncEventLogger
}

type AsyncEventLogger interface {
	Warn(msg string, tags ...any)
	Error(msg string, tags ...any)
}

type AsyncStats struct {
	Queued    uint64
	Forwarded uint64
	Dropped   uint64
	Failed    uint64
	Retried   uint64
}

type AsyncSnapshot struct {
	Queued        uint64 `json:"queued"`
	Forwarded     uint64 `json:"forwarded"`
	Dropped       uint64 `json:"dropped"`
	Failed        uint64 `json:"failed"`
	Retried       uint64 `json:"retried"`
	QueueDepth    int    `json:"queueDepth"`
	QueueCapacity int    `json:"queueCapacity"`
	Closed        bool   `json:"closed"`
}

type AsyncSink struct {
	sink  Sink
	opts  AsyncOptions
	queue chan Record
	done  chan struct{}

	mu     sync.Mutex
	closed bool

	queued    atomic.Uint64
	forwarded atomic.Uint64
	dropped   atomic.Uint64
	failed    atomic.Uint64
	retried   atomic.Uint64
}

func NewAsyncSink(sink Sink, opts AsyncOptions) *AsyncSink {
	if opts.QueueSize <= 0 {
		opts.QueueSize = 256
	}
	if opts.InitialBackoff <= 0 {
		opts.InitialBackoff = 100 * time.Millisecond
	}
	if opts.MaxBackoff <= 0 {
		opts.MaxBackoff = 2 * time.Second
	}
	if opts.MaxBackoff < opts.InitialBackoff {
		opts.MaxBackoff = opts.InitialBackoff
	}

	a := &AsyncSink{
		sink:  sink,
		opts:  opts,
		queue: make(chan Record, opts.QueueSize),
		done:  make(chan struct{}),
	}
	go a.worker()
	return a
}

func (a *AsyncSink) Forward(_ context.Context, rec Record) error {
	if a == nil || a.sink == nil {
		return nil
	}

	a.mu.Lock()
	defer a.mu.Unlock()
	if a.closed {
		return ErrAsyncSinkClosed
	}

	select {
	case a.queue <- rec:
		a.queued.Add(1)
		return nil
	default:
		a.dropped.Add(1)
		if a.opts.Logger != nil {
			a.opts.Logger.Warn(
				"Log forwarding record dropped",
				"event", "log_forwarding_drop",
				"reason", "queue_full",
				"requestID", rec.RequestID,
				"dag", rec.DAGName,
				"step", rec.StepName,
				"queueDepth", len(a.queue),
				"queueCapacity", cap(a.queue),
			)
		}
		return ErrQueueFull
	}
}

func (a *AsyncSink) Close(ctx context.Context) error {
	if a == nil {
		return nil
	}

	a.mu.Lock()
	if a.closed {
		a.mu.Unlock()
		return nil
	}
	a.closed = true
	close(a.queue)
	a.mu.Unlock()

	select {
	case <-a.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (a *AsyncSink) Stats() AsyncStats {
	if a == nil {
		return AsyncStats{}
	}
	return AsyncStats{
		Queued:    a.queued.Load(),
		Forwarded: a.forwarded.Load(),
		Dropped:   a.dropped.Load(),
		Failed:    a.failed.Load(),
		Retried:   a.retried.Load(),
	}
}

func (a *AsyncSink) Snapshot() AsyncSnapshot {
	if a == nil {
		return AsyncSnapshot{}
	}

	a.mu.Lock()
	closed := a.closed
	depth := len(a.queue)
	capacity := cap(a.queue)
	a.mu.Unlock()

	return AsyncSnapshot{
		Queued:        a.queued.Load(),
		Forwarded:     a.forwarded.Load(),
		Dropped:       a.dropped.Load(),
		Failed:        a.failed.Load(),
		Retried:       a.retried.Load(),
		QueueDepth:    depth,
		QueueCapacity: capacity,
		Closed:        closed,
	}
}

func (a *AsyncSink) worker() {
	defer close(a.done)
	for rec := range a.queue {
		a.deliver(rec)
	}
}

func (a *AsyncSink) deliver(rec Record) {
	attempts := a.opts.MaxRetries + 1
	backoff := a.opts.InitialBackoff

	for i := 0; i < attempts; i++ {
		err := a.sink.Forward(context.Background(), rec)
		if err == nil {
			a.forwarded.Add(1)
			return
		}

		if i == attempts-1 {
			a.failed.Add(1)
			if a.opts.Logger != nil {
				a.opts.Logger.Error(
					"Log forwarding delivery failed",
					"event", "log_forwarding_failed",
					"requestID", rec.RequestID,
					"dag", rec.DAGName,
					"step", rec.StepName,
					"attempt", i+1,
					"maxAttempts", attempts,
					"error", err,
				)
			}
			return
		}

		a.retried.Add(1)
		if a.opts.Logger != nil {
			a.opts.Logger.Warn(
				"Log forwarding delivery retry",
				"event", "log_forwarding_retry",
				"requestID", rec.RequestID,
				"dag", rec.DAGName,
				"step", rec.StepName,
				"attempt", i+1,
				"maxAttempts", attempts,
				"backoff", backoff,
				"error", err,
			)
		}
		time.Sleep(backoff)
		backoff = backoff * 2
		if backoff > a.opts.MaxBackoff {
			backoff = a.opts.MaxBackoff
		}
	}
}
