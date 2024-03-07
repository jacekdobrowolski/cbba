package main

type TaskID uint
type Task struct {
	id              TaskID
	position        Position
	bid             float64
	highestBidderId int
}

var taskIDSequence TaskID

func newTask(position Position) Task {
	result := Task{
		id:              taskIDSequence,
		position:        position,
		bid:             0,
		highestBidderId: -1,
	}
	taskIDSequence++
	return result
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
