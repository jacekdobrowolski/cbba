package topic

import (
	"context"
	"log/slog"
	"sync"
)

type Topic[T any] struct {
	pub    chan T
	subs   []func(T)
	mu     sync.RWMutex
	logger *slog.Logger
}

func New[T any](logger *slog.Logger) *Topic[T] {
	t := &Topic[T]{
		pub:    make(chan T),
		logger: logger,
	}
	return t
}

func (t *Topic[T]) Start(ctx context.Context) {
out:
	for {
		select {
		case msg := <-t.pub:
			go func() {
				t.logger.Info("Topic received message", "msg", msg)
				t.mu.RLock()
				defer t.mu.RUnlock()
				for _, callback := range t.subs {
					callback := callback
					go callback(msg)
				}
			}()
		case <-ctx.Done():
			break out
		}
	}
}

func (t *Topic[T]) Subscribe(callback func(T)) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.subs = append(t.subs, callback)
}

func (t *Topic[T]) Publish(ctx context.Context, msg T) {
	select {
	case t.pub <- msg:
	case <-ctx.Done():
	}
}
