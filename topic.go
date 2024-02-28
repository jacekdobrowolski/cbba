package main

import (
	"context"
	"log/slog"
	"sync"
)

type Topic[T any] struct {
	pub    chan T
	subs   []chan T
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
			for _, sub := range t.subs {
				select {
				case sub <- msg:
				default:
					t.logger.Warn("Channel does not listen")
					continue
				}
			}
			t.mu.Unlock()
		case <-ctx.Done():
			break out
		}
	}
}

func (t *Topic[T]) subscribe() chan T {
	t.mu.Lock()
	defer t.mu.Unlock()
	new := make(chan T, 10)
	t.subs = append(t.subs, new)
	return new
}
