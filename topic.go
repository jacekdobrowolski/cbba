package main

import (
	"context"
	"log/slog"
	"sync"
)

type Topic[T any] struct {
	pub    chan T
	subs   []func(T)
	mu     sync.Mutex
	logger *slog.Logger
}

func newTopic[T any](logger *slog.Logger) *Topic[T] {
	t := &Topic[T]{
		pub:    make(chan T),
		logger: logger,
	}
	return t
}

func (t *Topic[T]) start(ctx context.Context) {
out:
	for {
		select {
		case msg := <-t.pub:
			t.logger.Info("Topic recived message", "msg", msg)
			t.mu.Lock()
			for _, callback := range t.subs {
				callback := callback
				go callback(msg)
			}
			t.mu.Unlock()
		case <-ctx.Done():
			break out
		}
	}
}

func (t *Topic[T]) subscribe(callback func(T)) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.subs = append(t.subs, callback)
}
