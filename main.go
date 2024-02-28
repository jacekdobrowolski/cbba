package main

import (
	"context"
	"fmt"
	"sync"
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

func (r *Rover) CBAA(ctx context.Context, topic *Topic[string]) {
	recive := topic.subscribe()
	go func() {
	out:
		for {
			select {
			case msg := <-recive:
				fmt.Printf("Agent %d recived message: %s\n", r.id, msg)
			case <-ctx.Done():
				break out

			}
		}
	}()

	time.Sleep(500 * time.Millisecond)
	topic.pub <- fmt.Sprintf("Hello from %d", r.id)
	fmt.Println("Message sent by agent ", r.id)
}

func main() {
	tasks := newTasksGrid(4, 4)
	fmt.Println("Tasks: ", tasks)
	rover0 := newRover(-1, -1)
	rover1 := newRover(4, 4)
	fmt.Println("Rovers: ", rover0, rover1)

	testTopic := newTopic[string]()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go testTopic.start(ctx, 1)
	go rover0.CBAA(ctx, testTopic)
	rover1.CBAA(ctx, testTopic)
	time.Sleep(1 * time.Second)
}

type Topic[T any] struct {
	pub  chan T
	subs []chan T
	mu   sync.Mutex
}

func newTopic[T any]() *Topic[T] {
	t := &Topic[T]{
		pub: make(chan T),
	}
	return t
}

func (t *Topic[T]) start(ctx context.Context, timeout time.Duration) {
	fmt.Println("Topic started")
out:
	for {
		select {
		case msg := <-t.pub:
			fmt.Printf("Topic recived message: %v \n", msg)
			t.mu.Lock()
			for _, sub := range t.subs {
				select {
				case sub <- msg:
				default:
					fmt.Println("channel does not listen")
					continue
				}
			}
			t.mu.Unlock()
			fmt.Println("Message sent to all subs")
		case <-ctx.Done():
			break out
		}
	}
}

func (t *Topic[T]) subscribe() chan T {
	t.mu.Lock()
	defer t.mu.Unlock()
	new := make(chan T)
	t.subs = append(t.subs, new)
	return new
}
