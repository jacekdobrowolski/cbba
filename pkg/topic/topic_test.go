package topic

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"testing"
)

func TestTopic_HappyCase(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError}))
	topic := New[int](logger)
	var sub1, sub2 int
	var wg sync.WaitGroup
	wg.Add(2)
	topic.Subscribe(func(msg int) {
		defer wg.Done()
		sub1 = msg
	})
	topic.Subscribe(func(msg int) {
		defer wg.Done()
		sub2 = msg
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go topic.Start(ctx)
	msg := 1
	topic.Publish(ctx, msg)
	wg.Wait()
	if sub1 == 0 || sub2 == 0 {
		t.Errorf("expected %d got %d and %d", msg, sub1, sub2)
	}
}

func BenchmarkTopic(b *testing.B) {
	numMessages := 100
	for n := 10; n <= 10_000; n *= 10 {
		n := n
		b.Run(fmt.Sprintf("topic-%d-subs-%d-msgs", n, numMessages), func(b *testing.B) {

			logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError}))
			topic := New[int64](logger)
			sub := make([]int64, n)
			b.StartTimer()
			muxs := make([]sync.Mutex, n)
			var wg sync.WaitGroup
			wg.Add(n * numMessages)
			for i := range n {
				i := i
				topic.Subscribe(func(msg int64) {
					defer wg.Done()
					muxs[i].Lock()
					defer muxs[i].Unlock()
					sub[i] += msg
				})
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			go topic.Start(ctx)
			var msg int64 = 1

			for range numMessages {
				topic.Publish(ctx, msg)
			}
			wg.Wait()
			b.StopTimer()

			var sum int64 = 0
			for i := range n {
				sum += sub[i]
			}
			if int64(numMessages*n) != sum {
				b.Errorf("expected %d got %d", numMessages*n, sum)
			}
		})
	}
}
