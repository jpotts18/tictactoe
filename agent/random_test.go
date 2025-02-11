package agent

import (
	"testing"

	"github.com/jpotts18/tictactoe/game"
)

func TestNewRandomAgent(t *testing.T) {
	agent := NewRandomAgent(1)
	if agent == nil {
		t.Error("NewRandomAgent() returned nil")
	}
}

func TestRandomAgentGetMove(t *testing.T) {
	tests := []struct {
		name          string
		boardSetup    []int // sequence of moves to set up the board
		player        int   // player making the move
		validMoves    []int // expected valid moves
	}{
		{
			name:       "empty board",
			boardSetup: []int{},
			player:     1,
			validMoves: []int{0, 1, 2, 3, 4, 5, 6, 7, 8},
		},
		{
			name:       "partially filled board",
			boardSetup: []int{0, 4, 8}, // X in corners and center
			player:     2,
			validMoves: []int{1, 2, 3, 5, 6, 7},
		},
		{
			name:       "nearly full board",
			boardSetup: []int{0, 1, 2, 3, 4, 5, 6, 7},
			player:     1,
			validMoves: []int{8},
		},
		{
			name:       "full board",
			boardSetup: []int{0, 1, 2, 3, 4, 5, 6, 7, 8},
			player:     1,
			validMoves: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent := NewRandomAgent(1	)
			g := game.NewTicTacToe()

			// Set up the board
			for _, move := range tt.boardSetup {
				err := g.MakeMove(move)
				if err != nil {
					t.Fatalf("Failed to set up board: %v", err)
				}
			}

			// Test multiple times to ensure randomness
			moveFrequency := make(map[int]int)
			numTests := 1000

			for i := 0; i < numTests; i++ {
				move := agent.GetMove(g.GetBoard(), tt.player)
				moveFrequency[move]++

				// Verify move is valid
				found := false
				for _, validMove := range tt.validMoves {
					if move == validMove {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("GetMove() returned invalid move %v, valid moves are %v", 
						move, tt.validMoves)
				}
			}

			// Verify all valid moves were used (testing randomness)
			for _, validMove := range tt.validMoves {
				if count := moveFrequency[validMove]; count == 0 {
					t.Errorf("Move %v was never chosen in %v attempts", validMove, numTests)
				}
			}

			// Check for reasonable distribution
			expectedFreq := float64(numTests) / float64(len(tt.validMoves))
			tolerance := expectedFreq * 0.5 // Allow 50% deviation

			for move, freq := range moveFrequency {
				if float64(freq) < expectedFreq-tolerance || 
				   float64(freq) > expectedFreq+tolerance {
					t.Errorf("Move %v frequency %v is outside expected range [%v, %v]",
						move, freq, expectedFreq-tolerance, expectedFreq+tolerance)
				}
			}
		})
	}
} 
