package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

type Rover struct {
	id           int
	position     Position
	plannedTasks []Task
	logger       *slog.Logger
}

var roverIDSequence int

func newRover(x, y int, logger *slog.Logger) Rover {

	rover := Rover{
		id:       roverIDSequence,
		position: Position{x, y},
		logger:   logger,
	}
	logger.Info("Rover created", "agent_id", rover.id)
	roverIDSequence++
	return rover
}

func (r *Rover) CBAA(ctx context.Context, topic *Topic[string]) {
	recive := topic.subscribe()
	go func() {
	out:
		for {
			select {
			case msg := <-recive:
				r.logger.Info("Agent recived message", "agent_id", r.id, "msg", msg)
			case <-ctx.Done():
				break out

			}
		}
	}()

	time.Sleep(500 * time.Millisecond)
	topic.pub <- fmt.Sprintf("Hello from %d", r.id)
	r.logger.Info("Message sent", "agent_id", r.id)
}
