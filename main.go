package main

import (
	"context"
	"log/slog"
	"math"
	"time"
)

type Position struct {
	x float64
	y float64
}

type Task struct {
	position        Position
	bid             float64
	highestBidderId int
}

func newTask(position Position) Task {
	return Task{
		position:        position,
		bid:             0,
		highestBidderId: -1,
	}
}

func newTasksGrid(width, height int) []Task {
	tasks := make([]Task, 0, width*height)
	for x := range width {
		for y := range height {
			tasks = append(tasks, newTask(Position{float64(x), float64(y)}))
		}
	}
	return tasks
}

func distance(a, b Position) float64 {
	return math.Hypot(a.x-b.x, a.y-b.y)
}

func main() {
	logger := slog.Default()

	tasks := newTasksGrid(2, 2)
	rover0 := newRover(-1, 0, logger)
	rover1 := newRover(2, 1, logger)

	auctionTopic := newTopic[Task](logger)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rover0.CBAAInit(ctx, auctionTopic)
	rover1.CBAAInit(ctx, auctionTopic)

	go auctionTopic.start(ctx)

	for _, task := range tasks {
		auctionTopic.pub <- task
	}
	time.Sleep(1 * time.Second) // lettting agents finish
}
