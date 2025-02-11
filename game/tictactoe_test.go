package game

import (
	"strings"
	"testing"
)

func TestNewTicTacToe(t *testing.T) {
	game := NewTicTacToe()

	if game.currentPlayer != 1 {
		t.Errorf("NewTicTacToe() currentPlayer = %v, want 1", game.currentPlayer)
	}
	if game.isGameOver {
		t.Error("NewTicTacToe() isGameOver = true, want false")
	}
	if game.winner != 0 {
		t.Errorf("NewTicTacToe() winner = %v, want 0", game.winner)
	}
}

func TestGameFlow(t *testing.T) {
	tests := []struct {
		name          string
		moves         []int
		wantErr       []bool
		wantPlayers   []int  // expected player after each move
		wantGameOver  bool
		wantWinner    int
	}{
		{
			name:          "alternating moves",
			moves:         []int{4, 0, 8},
			wantErr:       []bool{false, false, false},
			wantPlayers:   []int{2, 1, 2},
			wantGameOver:  false,
			wantWinner:    0,
		},
		{
			name:          "attempt move after game over",
			moves:         []int{0, 3, 1, 4, 2, 8}, // Player 1 wins, then invalid move
			wantErr:       []bool{false, false, false, false, false, true},
			wantPlayers:   []int{2, 1, 2, 1, 1, 1}, // Player stays at 1 after win
			wantGameOver:  true,
			wantWinner:    1,
		},
		{
			name:          "draw game",
			moves:         []int{0, 1, 2, 4, 3, 6, 5, 8, 7},
			wantErr:       []bool{false, false, false, false, false, false, false, false, false},
			wantPlayers:   []int{2, 1, 2, 1, 2, 1, 2, 1, 1}, // Last player stays at 1 after draw
			wantGameOver:  true,
			wantWinner:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := NewTicTacToe()

			// Execute moves and check state after each
			for i, move := range tt.moves {
				err := game.MakeMove(move)
				if (err != nil) != tt.wantErr[i] {
					t.Errorf("Move %d: MakeMove() error = %v, wantErr %v", 
						i, err, tt.wantErr[i])
				}
				if game.GetCurrentPlayer() != tt.wantPlayers[i] {
					t.Errorf("Move %d: currentPlayer = %v, want %v", 
						i, game.GetCurrentPlayer(), tt.wantPlayers[i])
				}
			}

			// Check final game state
			if game.IsGameOver() != tt.wantGameOver {
				t.Errorf("IsGameOver() = %v, want %v", 
					game.IsGameOver(), tt.wantGameOver)
			}
			if game.GetWinner() != tt.wantWinner {
				t.Errorf("GetWinner() = %v, want %v", 
					game.GetWinner(), tt.wantWinner)
			}
		})
	}
}

func TestGameDisplay(t *testing.T) {
	tests := []struct {
		name     string
		moves    []int
		contains []string
	}{
		{
			name:     "new game",
			moves:    []int{},
			contains: []string{
				"Current player: 1",
				"_|_|_\n",
				"_|_|_\n",
				"_|_|_",
			},
		},
		{
			name:     "mid-game",
			moves:    []int{4, 0},
			contains: []string{
				"Current player: 1",
				"2|_|_\n",
				"_|1|_\n",
				"_|_|_",
			},
		},
		{
			name:     "winning game",
			moves:    []int{0, 3, 1, 4, 2}, // Player 1 wins with top row
			contains: []string{
				"Game Over: Player 1 wins!",
				"1|1|1\n",
				"2|2|_\n",
				"_|_|_",
			},
		},
		{
			name:     "draw game",
			moves:    []int{0, 1, 2, 4, 3, 6, 5, 8, 7},
			contains: []string{
				"Game Over: Draw!",
				"1|2|1\n",
				"1|2|1\n",
				"2|1|2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := NewTicTacToe()
			
			// Execute moves
			for _, move := range tt.moves {
				game.MakeMove(move)
			}

			got := game.String()
			
			// Check if output contains expected strings
			for _, want := range tt.contains {
				if !strings.Contains(got, want) {
					t.Errorf("String() = \n%v\ndoes not contain \n%v", got, want)
				}
			}
		})
	}
}

// Helper function to check if a string contains another string
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[0:len(s)][0:len(substr)] == substr
} 
