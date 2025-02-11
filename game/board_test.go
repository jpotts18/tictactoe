package game

import (
	"testing"
)

func TestNewBoard(t *testing.T) {
	b := NewBoard()
	state := b.GetState()

	// Check if board is initialized with zeros
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if state[i][j] != 0 {
				t.Errorf("NewBoard() cell [%d][%d] = %d; want 0", i, j, state[i][j])
			}
		}
	}
}

func TestSetCell(t *testing.T) {
	tests := []struct {
		name     string
		pos      int
		value    int
		want     bool
		wantCell int
	}{
		{
			name:     "valid move center",
			pos:      4,
			value:    1,
			want:     true,
			wantCell: 1,
		},
		{
			name:     "invalid position negative",
			pos:      -1,
			value:    1,
			want:     false,
			wantCell: 0,
		},
		{
			name:     "invalid position too large",
			pos:      9,
			value:    1,
			want:     false,
			wantCell: 0,
		},
		{
			name:     "position already occupied",
			pos:      4,
			value:    2,
			want:     false,
			wantCell: 1, // remains unchanged from previous test
		},
	}

	b := NewBoard()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := b.SetCell(tt.pos, tt.value)
			if got != tt.want {
				t.Errorf("SetCell(%v, %v) = %v, want %v", tt.pos, tt.value, got, tt.want)
			}

			if tt.pos >= 0 && tt.pos < 9 {
				if cell := b.GetCell(tt.pos); cell != tt.wantCell {
					t.Errorf("After SetCell(%v, %v), GetCell(%v) = %v, want %v", 
						tt.pos, tt.value, tt.pos, cell, tt.wantCell)
				}
			}
		})
	}
}

func TestGetCell(t *testing.T) {
	b := NewBoard()
	b.SetCell(4, 1) // Set center cell

	tests := []struct {
		name string
		pos  int
		want int
	}{
		{
			name: "get empty cell",
			pos:  0,
			want: 0,
		},
		{
			name: "get occupied cell",
			pos:  4,
			want: 1,
		},
		{
			name: "get last cell",
			pos:  8,
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := b.GetCell(tt.pos); got != tt.want {
				t.Errorf("GetCell(%v) = %v, want %v", tt.pos, got, tt.want)
			}
		})
	}
}

func TestGetEmptyCells(t *testing.T) {
	tests := []struct {
		name      string
		positions []int // positions to fill
		values    []int // values to fill with
		want      []int // expected empty positions
	}{
		{
			name:      "empty board",
			positions: []int{},
			values:    []int{},
			want:      []int{0, 1, 2, 3, 4, 5, 6, 7, 8},
		},
		{
			name:      "some cells filled",
			positions: []int{0, 4, 8},
			values:    []int{1, 2, 1},
			want:      []int{1, 2, 3, 5, 6, 7},
		},
		{
			name:      "full board",
			positions: []int{0, 1, 2, 3, 4, 5, 6, 7, 8},
			values:    []int{1, 2, 1, 2, 1, 2, 1, 2, 1},
			want:      []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBoard()
			
			// Fill board according to test case
			for i := range tt.positions {
				b.SetCell(tt.positions[i], tt.values[i])
			}

			got := b.GetEmptyCells()
			
			if len(got) != len(tt.want) {
				t.Errorf("GetEmptyCells() returned %v cells, want %v cells", 
					len(got), len(tt.want))
			}

			for i, cell := range tt.want {
				if got[i] != cell {
					t.Errorf("GetEmptyCells()[%d] = %v, want %v", i, got[i], cell)
				}
			}
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name      string
		positions []int
		values    []int
		want      string
	}{
		{
			name:      "empty board",
			positions: []int{},
			values:    []int{},
			want:      "_|_|_\n-----\n_|_|_\n-----\n_|_|_",
		},
		{
			name:      "some moves made",
			positions: []int{0, 4, 8},
			values:    []int{1, 2, 1},
			want:      "1|_|_\n-----\n_|2|_\n-----\n_|_|1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBoard()
			
			// Fill board according to test case
			for i := range tt.positions {
				b.SetCell(tt.positions[i], tt.values[i])
			}

			if got := b.String(); got != tt.want {
				t.Errorf("String() = \n%v\nwant\n%v", got, tt.want)
			}
		})
	}
}

func TestSetState(t *testing.T) {
	tests := []struct {
		name  string
		state [][]int
		want  [3][3]int
	}{
		{
			name: "empty state",
			state: [][]int{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			want: [3][3]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
		},
		{
			name: "mixed state",
			state: [][]int{
				{1, 0, 2},
				{0, 1, 0},
				{2, 0, 1},
			},
			want: [3][3]int{{1, 0, 2}, {0, 1, 0}, {2, 0, 1}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBoard()
			b.SetState(tt.state)
			
			got := b.GetState()
			if got != tt.want {
				t.Errorf("After SetState(), GetState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMakeMove(t *testing.T) {
	tests := []struct {
		name      string
		position  int
		player    int
		initial   [][]int
		want      bool
		wantState [][]int
	}{
		{
			name:     "valid move empty cell",
			position: 4,
			player:   1,
			initial: [][]int{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			want: true,
			wantState: [][]int{
				{0, 0, 0},
				{0, 1, 0},
				{0, 0, 0},
			},
		},
		{
			name:     "invalid move occupied cell",
			position: 4,
			player:   2,
			initial: [][]int{
				{0, 0, 0},
				{0, 1, 0},
				{0, 0, 0},
			},
			want: false,
			wantState: [][]int{
				{0, 0, 0},
				{0, 1, 0},
				{0, 0, 0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBoard()
			b.SetState(tt.initial)
			
			got := b.MakeMove(tt.position, tt.player)
			
			if got != tt.want {
				t.Errorf("MakeMove() = %v, want %v", got, tt.want)
			}
			
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					if b.state[i][j] != tt.wantState[i][j] {
						t.Errorf("board state at [%d][%d] = %d, want %d", 
							i, j, b.state[i][j], tt.wantState[i][j])
					}
				}
			}
		})
	}
}

func TestIsGameOver(t *testing.T) {
	tests := []struct {
		name      string
		state     [][]int
		wantOver  bool
		wantWinner int
	}{
		{
			name: "game in progress",
			state: [][]int{
				{1, 0, 2},
				{0, 1, 0},
				{0, 0, 0},
			},
			wantOver:  false,
			wantWinner: 0,
		},
		{
			name: "row win",
			state: [][]int{
				{1, 1, 1},
				{2, 2, 0},
				{0, 0, 0},
			},
			wantOver:  true,
			wantWinner: 1,
		},
		{
			name: "column win",
			state: [][]int{
				{2, 1, 0},
				{2, 1, 0},
				{2, 0, 0},
			},
			wantOver:  true,
			wantWinner: 2,
		},
		{
			name: "diagonal win",
			state: [][]int{
				{1, 2, 0},
				{2, 1, 0},
				{0, 0, 1},
			},
			wantOver:  true,
			wantWinner: 1,
		},
		{
			name: "draw game",
			state: [][]int{
				{1, 2, 1},
				{1, 2, 2},
				{2, 1, 1},
			},
			wantOver:  true,
			wantWinner: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBoard()
			b.SetState(tt.state)
			
			gotOver, gotWinner := b.IsGameOver()
			
			if gotOver != tt.wantOver || gotWinner != tt.wantWinner {
				t.Errorf("IsGameOver() = (%v, %v), want (%v, %v)", 
					gotOver, gotWinner, tt.wantOver, tt.wantWinner)
			}
		})
	}
}

func TestGetAvailableMoves(t *testing.T) {
	tests := []struct {
		name  string
		state [][]int
		want  []int
	}{
		{
			name: "empty board",
			state: [][]int{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
			want: []int{0, 1, 2, 3, 4, 5, 6, 7, 8},
		},
		{
			name: "partially filled board",
			state: [][]int{
				{1, 0, 2},
				{0, 1, 0},
				{2, 0, 0},
			},
			want: []int{1, 3, 5, 7, 8},
		},
		{
			name: "full board",
			state: [][]int{
				{1, 2, 1},
				{2, 1, 2},
				{1, 2, 1},
			},
			want: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBoard()
			b.SetState(tt.state)
			
			got := b.GetAvailableMoves()
			
			if len(got) != len(tt.want) {
				t.Errorf("GetAvailableMoves() returned %v moves, want %v moves", 
					len(got), len(tt.want))
			}
			
			// Check if all expected moves are present
			for i, move := range tt.want {
				if got[i] != move {
					t.Errorf("GetAvailableMoves()[%d] = %v, want %v", 
						i, got[i], move)
				}
			}
		})
	}
} 
