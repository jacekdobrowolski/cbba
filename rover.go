package main

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
)

type Rover struct {
	id           int
	position     Position
	plannedTasks map[string]Task
	mu           sync.RWMutex
	logger       *slog.Logger
}

var roverIDSequence int

func newRover(x, y float64, logger *slog.Logger) *Rover {

	rover := &Rover{
		id:           roverIDSequence,
		position:     Position{x, y},
		plannedTasks: make(map[string]Task),
		logger:       logger,
	}
	logger.Info("Rover created", "agent_id", rover.id, "position", rover.position)
	roverIDSequence++
	return rover
}

func (r *Rover) CBAAInit(ctx context.Context, topic *Topic[Task]) {
	recive := topic.subscribe()
	go func() {
	out:
		for {
			select {
			case msg := <-recive:
				go func() {
					r.logger.Info("Agent recived message", "agent_id", r.id, "msg", msg)
					r.mu.Lock()
					currentState, ok := r.plannedTasks[fmt.Sprintf("%+v", msg.position)]
					if !ok || currentState.bid < msg.bid {
						r.plannedTasks[fmt.Sprintf("%+v", msg.position)] = msg
						go r.updateBid(ctx, msg, topic)

					}
					r.mu.Unlock()
					r.logger.Info("Agent tasks updated", "agent_id", r.id, "planned_tasks", r.plannedTasks)
				}()
			case <-ctx.Done():
				break out

			}
		}
	}()
	r.logger.Info("Agent ready", "agent_id", r.id)
}

func (r *Rover) updateBid(ctx context.Context, task Task, topic *Topic[Task]) {
	r.mu.RLock()
	taskReward := 1 / distance(r.position, task.position)
	if taskReward > task.bid {
		task.bid = taskReward
		task.highestBidderId = r.id
	}
	r.mu.RUnlock()
	select {
	case topic.pub <- task:
	case <-ctx.Done():
	}
}
