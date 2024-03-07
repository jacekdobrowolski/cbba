package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/jacekdobrowolski/cbba/internal/cbba"
	"github.com/jacekdobrowolski/cbba/pkg/topic"
)

func main() {
	logger := slog.Default()

	tasks := cbba.NewTasksGrid(2, 2)
	agent0 := cbba.NewAgent(-1, 0, logger)
	agent1 := cbba.NewAgent(2, 1, logger)

	auctionTopic := topic.New[cbba.Task](logger)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	agent0.CBAAInit(ctx, auctionTopic)
	agent1.CBAAInit(ctx, auctionTopic)

	go auctionTopic.Start(ctx)

	for _, task := range tasks {
		auctionTopic.Publish(ctx, task)
	}
	time.Sleep(1 * time.Second) // letting agents finish
}
