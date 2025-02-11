package agent

import (
	"testing"

	"github.com/jpotts18/tictactoe/game"
)

func TestGetStateKey(t *testing.T) {
	tests := []struct {
		name     string
		board    [][]int
		expected string
	}{
		{
			name: "empty board",
			board: [][]int{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			expected: "000000000",
		},
		{
			name: "board with X's and O's",
			board: [][]int{
				{1, 0, 2},
				{0, 1, 0},
				{2, 0, 1},
			},
			expected: "102010201",
		},
		{
			name: "full board",
			board: [][]int{
				{1, 2, 1},
				{2, 1, 2},
				{1, 2, 1},
			},
			expected: "121212121",
		},
	}

	agent := &BaseAgent{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := game.NewBoard()
			b.SetState(tt.board)
			
			result := agent.GetStateKey(b)
			
			if result != tt.expected {
				t.Errorf("GetStateKey() = %v, want %v", result, tt.expected)
			}
		})
	}
}
