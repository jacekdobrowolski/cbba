package main

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
)

type Agent struct {
	id           int
	position     Position
	plannedTasks map[string]Task
	mu           sync.RWMutex
	logger       *slog.Logger
}

var agentIDSequence int

func newAgent(x, y float64, logger *slog.Logger) *Agent {

	agent := &Agent{
		id:           agentIDSequence,
		position:     Position{x, y},
		plannedTasks: make(map[string]Task),
		logger:       logger,
	}
	logger.Info("Agent created", "agent_id", agent.id, "position", agent.position)
	agentIDSequence++
	return agent
}

func (a *Agent) CBAAInit(ctx context.Context, topic *Topic[Task]) {
	recive := topic.subscribe()
	go func() {
	out:
		for {
			select {
			case msg := <-recive:
				go func() {
					a.logger.Info("Agent recived message", "agent_id", a.id, "msg", msg)
					a.mu.Lock()
					currentState, ok := a.plannedTasks[fmt.Sprintf("%+v", msg.position)]
					if !ok || currentState.bid < msg.bid {
						a.plannedTasks[fmt.Sprintf("%+v", msg.position)] = msg
						go a.updateBid(ctx, msg, topic)

					}
					a.mu.Unlock()
					a.logger.Info("Agent tasks updated", "agent_id", a.id, "planned_tasks", a.plannedTasks)
				}()
			case <-ctx.Done():
				break out

			}
		}
	}()
	a.logger.Info("Agent ready", "agent_id", a.id)
}

func (a *Agent) updateBid(ctx context.Context, task Task, topic *Topic[Task]) {
	a.mu.RLock()
	taskReward := 1 / distance(a.position, task.position)
	if taskReward > task.bid {
		task.bid = taskReward
		task.highestBidderId = a.id
	}
	a.mu.RUnlock()
	select {
	case topic.pub <- task:
	case <-ctx.Done():
	}
}
