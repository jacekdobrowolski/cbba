package cbba

import (
	"math"
	"testing"
)

func TestAgent_scoreDistance(t *testing.T) {
	epsilon := 2*math.Nextafter(1, 2) - 1
	data := []struct {
		name     string
		a        Position
		b        Position
		expected float64
	}{
		{"zero case", Position{0, 0}, Position{0, 0}, MaxTaskScore},
		{"one in y", Position{0, 0}, Position{0, 1}, 1},
		{"one in x", Position{0, 0}, Position{1, 0}, 1},
		{"1 diagonal", Position{0, 0}, Position{1, 1}, 1 / math.Sqrt(2)},
		{"1 diagonal switched args", Position{1, 1}, Position{0, 0}, 1 / math.Sqrt(2)},
		{"1 diagonal negative", Position{0, 0}, Position{-1, -1}, 1 / math.Sqrt(2)},
		{"sqrt34 comb0", Position{-1, -3}, Position{2, 2}, 1 / math.Sqrt(34)},
		{"sqrt34 comb1", Position{-1, 2}, Position{2, -3}, 1 / math.Sqrt(34)},
		{"sqrt34 comb2", Position{2, -3}, Position{-1, 2}, 1 / math.Sqrt(34)},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			result := scoreDistance(d.a, d.b)
			if math.Abs(result-d.expected) > epsilon {
				t.Errorf("Expected %.15f, got %.15f", d.expected, result)
			}
		})
	}
}

func TestAgent_scorePath(t *testing.T) {
	type fields struct {
		position Position
		tasks    map[TaskID]Task
	}
	type args struct {
		path []TaskID
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{"empty case", fields{Position{0, 0}, map[TaskID]Task{}}, args{[]TaskID{}}, 0},
		{"empty path", fields{Position{0, 0}, map[TaskID]Task{0: NewTask(Position{0, 1})}}, args{[]TaskID{0}}, 1},
		{"one in path new task", fields{Position{0, 0}, map[TaskID]Task{0: NewTask(Position{0, 1}), 1: NewTask(Position{0, 2})}}, args{[]TaskID{0, 1}}, 2},
		{"loop", fields{Position{1, 1}, map[TaskID]Task{0: NewTask(Position{-1, -1}), 1: NewTask(Position{1, 1})}}, args{[]TaskID{0, 1}}, 1 / math.Sqrt(2)},
		{"starting on top of a task", fields{Position{1, 1}, map[TaskID]Task{0: NewTask(Position{-1, -1}), 1: NewTask(Position{1, 1})}}, args{[]TaskID{1, 0}}, MaxTaskScore + 1.0/(2*math.Sqrt(2))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Agent{
				position: tt.fields.position,
				tasks:    tt.fields.tasks,
			}
			if got := a.scorePath(tt.args.path); got != tt.want {
				t.Errorf("Agent.scorePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAgent_bestPlaceToInsertTaskIntoPath(t *testing.T) {
	type fields struct {
		position Position
		path     []TaskID
		tasks    map[TaskID]Task
	}
	type args struct {
		task_id TaskID
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{
			"at the end",
			fields{Position{0, -1},
				[]TaskID{0, 1},
				map[TaskID]Task{
					0: NewTask(Position{0, 0}),
					1: NewTask(Position{0, 1}),
					2: NewTask(Position{0, 2})}},
			args{2},
			2,
		},
		{
			"in the middle",
			fields{Position{0, -1},
				[]TaskID{0, 2},
				map[TaskID]Task{
					0: NewTask(Position{0, 0}),
					1: NewTask(Position{0, 1}),
					2: NewTask(Position{0, 2})}},
			args{1},
			1,
		},
		{
			"at the beginning",
			fields{Position{0, -1},
				[]TaskID{1, 2},
				map[TaskID]Task{
					0: NewTask(Position{0, 0}),
					1: NewTask(Position{0, 1}),
					2: NewTask(Position{0, 2})}},
			args{0},
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Agent{
				position: tt.fields.position,
				path:     tt.fields.path,
				tasks:    tt.fields.tasks,
			}
			if got := a.bestPlaceToInsertTaskIntoPath(tt.args.task_id); got != tt.want {
				t.Errorf("Agent.bestPlaceToInsertTaskIntoPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
