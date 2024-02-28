package main

import (
	"fmt"
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

type Rover struct {
	id           int
	position     Position
	plannedTasks []Task
}

var roverIDSequence int

func newRover(x, y int) Rover {

	rover := Rover{
		id:       roverIDSequence,
		position: Position{x, y},
	}
	roverIDSequence++
	return rover
}

func main() {
	tasks := newTasksGrid(4, 4)
	fmt.Println("Tasks: ", tasks)
	rover0 := newRover(-1, -1)
	rover1 := newRover(4, 4)
	fmt.Println("Rovers: ", rover0, rover1)
}
