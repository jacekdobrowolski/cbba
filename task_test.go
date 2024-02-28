package main

import (
	"math"
	"testing"
)

func TestDistance(t *testing.T) {
	epsilon := 2*math.Nextafter(1, 2) - 1
	data := []struct {
		name     string
		a        Position
		b        Position
		expected float64
	}{
		{"zero case", Position{0, 0}, Position{0, 0}, 0},
		{"one in y", Position{0, 0}, Position{0, 1}, 1},
		{"one in x", Position{0, 0}, Position{1, 0}, 1},
		{"1 diagonal", Position{0, 0}, Position{1, 1}, math.Sqrt(2)},
		{"1 diagonal switched args", Position{1, 1}, Position{0, 0}, math.Sqrt(2)},
		{"1 diagonal negative", Position{0, 0}, Position{-1, -1}, math.Sqrt(2)},
		{"sqrt34 comb0", Position{-1, -3}, Position{2, 2}, math.Sqrt(34)},
		{"sqrt34 comb1", Position{-1, 2}, Position{2, -3}, math.Sqrt(34)},
		{"sqrt34 comb2", Position{2, -3}, Position{-1, 2}, math.Sqrt(34)},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			result := distance(d.a, d.b)
			if math.Abs(result-d.expected) > epsilon {
				t.Errorf("Expected %.15f, got %.15f", d.expected, result)
			}
		})
	}
}
