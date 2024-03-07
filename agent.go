package main

import (
	"context"
	"log/slog"
	"math"
	"sync"
)

const MaxTaskScore = 100

type Agent struct {
	id       int
	position Position
	bundle   []TaskID
	path     []TaskID
	tasks    map[TaskID]Task
	mu       sync.RWMutex
	logger   *slog.Logger
}

var agentIDSequence int

func newAgent(x, y float64, logger *slog.Logger) *Agent {

	agent := &Agent{
		id:       agentIDSequence,
		position: Position{x, y},
		bundle:   make([]TaskID, 0),
		path:     make([]TaskID, 0),
		tasks:    make(map[TaskID]Task),
		logger:   logger,
	}
	logger.Info("Agent created", "agent_id", agent.id, "position", agent.position)
	agentIDSequence++
	return agent
}

func (a *Agent) CBAAInit(ctx context.Context, topic *Topic[Task]) {
	topic.subscribe(func(receivedTask Task) {
		a.logger.Info("Agent received message", "agent_id", a.id, "msg", receivedTask)
		a.mu.Lock()
		currentState, ok := a.tasks[receivedTask.id]
		if !ok || currentState.bid < receivedTask.bid {
			a.tasks[receivedTask.id] = receivedTask
			go a.updateBid(ctx, receivedTask, topic)

		}
		a.mu.Unlock()
		a.logger.Info("Agent tasks updated", "agent_id", a.id, "planned_tasks", a.tasks)
	})
	a.logger.Info("Agent ready", "agent_id", a.id)
}

func (a *Agent) updateBid(ctx context.Context, task Task, topic *Topic[Task]) {
	a.mu.RLock()
	taskReward := scoreDistance(a.position, task.position)
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

func scoreDistance(a, b Position) float64 {
	distance := math.Hypot(a.x-b.x, a.y-b.y)
	if distance < math.Nextafter(0, 1) {
		distance = math.Nextafter(0, 1)
	}
	score := 1 / distance
	if score > MaxTaskScore {
		score = MaxTaskScore
	}
	return score
}

func (a *Agent) scorePath(path []TaskID) float64 {
	previousPosition := a.position
	var score float64 = 0
	for _, taskID := range path {
		score += scoreDistance(previousPosition, a.tasks[taskID].position)
		previousPosition = a.tasks[taskID].position
	}
	return score
}

func (a *Agent) bestPlaceToInsertTaskIntoPath(task_id TaskID) int {
	argmax := 0
	var maxScore float64 = 0
	newPath := make([]TaskID, 0, len(a.path)+1)
	for i := 0; i <= len(a.path); i++ {
		newPath := append(newPath, a.path[:i]...)
		newPath = append(newPath, task_id)
		newPath = append(newPath, a.path[i:len(a.path)]...)
		insertionScore := a.scorePath(newPath)
		if insertionScore > maxScore {
			maxScore = insertionScore
			argmax = i
		}
	}
	return argmax
}
