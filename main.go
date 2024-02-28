package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

type Position struct {
	x int
	y int
}

type Task Position

func newTasksGrid(width, height int) []Task {
	tasks := make([]Task, 0, width*height)
	for x := range width {
		for y := range height {
			tasks = append(tasks, Task{x, y})
		}
	}
	return tasks
}

func main() {
	logger := slog.Default()

	tasks := newTasksGrid(4, 4)
	fmt.Println("Tasks: ", tasks)
	rover0 := newRover(-1, -1, logger)
	rover1 := newRover(4, 4, logger)
	fmt.Println("Rovers: ", rover0, rover1)
	testTopic := newTopic[string](logger)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go testTopic.start(ctx)
	go rover0.CBAA(ctx, testTopic)
	rover1.CBAA(ctx, testTopic)
	time.Sleep(1 * time.Second)
}
