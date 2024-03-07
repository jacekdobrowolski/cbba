package main

import (
	"reflect"
	"testing"
)

func Test_newTasksGrid(t *testing.T) {
	type args struct {
		width  int
		height int
	}
	tests := []struct {
		name string
		args args
		want []Task
	}{
		{"simple 1x2 grid", args{1, 2}, []Task{newTask(Position{0, 0}), newTask(Position{0, 1})}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newTasksGrid(tt.args.width, tt.args.height)
			for i, task := range got {
				if !reflect.DeepEqual(task.position, tt.want[i].position) {
					t.Errorf("newTasksGrid() = %v, want %v", task.position, tt.want[i].position)
				}
			}
		})
	}
}
