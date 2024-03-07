package main

import (
	"context"
	"log/slog"
	"time"
)

type Position struct {
	x float64
	y float64
}

func main() {
	logger := slog.Default()

	tasks := newTasksGrid(2, 2)
	agent0 := newAgent(-1, 0, logger)
	agent1 := newAgent(2, 1, logger)

	auctionTopic := newTopic[Task](logger)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	agent0.CBAAInit(ctx, auctionTopic)
	agent1.CBAAInit(ctx, auctionTopic)

	go auctionTopic.start(ctx)

	for _, task := range tasks {
		auctionTopic.pub <- task
	}
	time.Sleep(1 * time.Second) // letting agents finish
}
