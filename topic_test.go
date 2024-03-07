package main

import (
	"context"
	"io"
	"log/slog"
	"sync"
	"testing"
)

func TestTopic_HappyCase(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError}))
	topic := newTopic[int](logger)
	var sub1, sub2 int
	var wg sync.WaitGroup
	wg.Add(2)
	topic.subscribe(func(msg int) {
		defer wg.Done()
		sub1 = msg
	})
	topic.subscribe(func(msg int) {
		defer wg.Done()
		sub2 = msg
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go topic.start(ctx)
	msg := 1
	topic.pub <- msg
	wg.Wait()
	if sub1 == 0 || sub2 == 0 {
		t.Errorf("expected %d got %d and %d", msg, sub1, sub2)
	}
}
