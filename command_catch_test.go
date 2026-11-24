package main

import "testing"

func TestCatchChance(t *testing.T) {
	tests := []struct {
		name          string
		baseExperience int
		wantMax       float64
	}{
		{"no base experience is always catchable", 0, 1.0},
		{"very rare pokemon is hard to catch", 400, 0.2},
		{"common pokemon is easy to catch", 10, 0.9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := catchChance(tt.baseExperience)

			if got <= 0 || got > 1 {
				t.Fatalf("catchChance(%d) = %v, want a value in (0, 1]", tt.baseExperience, got)
			}

			if got > tt.wantMax {
				t.Errorf("catchChance(%d) = %v, want at most %v", tt.baseExperience, got, tt.wantMax)
			}
		})
	}
}

func TestCatchChanceDecreasesWithBaseExperience(t *testing.T) {
	low := catchChance(10)
	high := catchChance(200)

	if low <= high {
		t.Errorf("expected a higher chance for a lower base experience, got %v vs %v", low, high)
	}
}
