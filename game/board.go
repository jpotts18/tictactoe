package game

import (
	"strconv"
)

// Board represents the game grid
type Board struct {
	state [3][3]int
}

// NewBoard creates a new empty board
func NewBoard() *Board {
	return &Board{}
}

// GetState returns the current board state
func (b *Board) GetState() [3][3]int {
	return b.state
}

// SetState sets the board state to the provided state
func (b *Board) SetState(state [][]int) {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			b.state[i][j] = state[i][j]
		}
	}
}

// SetCell sets a value at the specified position
func (b *Board) SetCell(pos int, value int) bool {
	row := pos / 3
	col := pos % 3
	if row < 0 || row >= 3 || col < 0 || col >= 3 {
		return false
	}
	if b.state[row][col] != 0 {
		return false
	}
	b.state[row][col] = value
	return true
}

// GetCell returns the value at the specified position
func (b *Board) GetCell(pos int) int {
	row := pos / 3
	col := pos % 3
	return b.state[row][col]
}

// IsEmpty checks if a position is empty
func (b *Board) IsEmpty(pos int) bool {
	return b.GetCell(pos) == 0
}

// GetEmptyCells returns positions of all empty cells
func (b *Board) GetEmptyCells() []int {
	var cells []int
	for i := 0; i < 9; i++ {
		if b.IsEmpty(i) {
			cells = append(cells, i)
		}
	}
	return cells
}

// String returns a string representation of the board
func (b *Board) String() string {
	var s string
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if b.state[i][j] == 0 {
				s += "_"
			} else {
				s += strconv.Itoa(b.state[i][j])
			}
			if j < 2 {
				s += "|"
			}
		}
		if i < 2 {
			s += "\n-----\n"
		}
	}
	return s
}

// MakeMove places a player's mark (1 or 2) at the specified position
func (b *Board) MakeMove(pos int, player int) bool {
	row := pos / 3
	col := pos % 3
	if b.state[row][col] == 0 {
		b.state[row][col] = player
		return true
	}
	return false
}

// GetAvailableMoves returns a slice of valid move positions
func (b *Board) GetAvailableMoves() []int {
	var moves []int
	for i := 0; i < 9; i++ {
		row := i / 3
		col := i % 3
		if b.state[row][col] == 0 {
			moves = append(moves, i)
		}
	}
	return moves
}

// IsGameOver checks if the game is over and returns the winner (0 for draw)
func (b *Board) IsGameOver() (bool, int) {
	// Check rows
	for i := 0; i < 3; i++ {
		if b.state[i][0] != 0 && b.state[i][0] == b.state[i][1] && b.state[i][1] == b.state[i][2] {
			return true, b.state[i][0]
		}
	}

	// Check columns
	for i := 0; i < 3; i++ {
		if b.state[0][i] != 0 && b.state[0][i] == b.state[1][i] && b.state[1][i] == b.state[2][i] {
			return true, b.state[0][i]
		}
	}

	// Check diagonals
	if b.state[0][0] != 0 && b.state[0][0] == b.state[1][1] && b.state[1][1] == b.state[2][2] {
		return true, b.state[0][0]
	}
	if b.state[0][2] != 0 && b.state[0][2] == b.state[1][1] && b.state[1][1] == b.state[2][0] {
		return true, b.state[0][2]
	}

	// Check for draw
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if b.state[i][j] == 0 {
				return false, 0
			}
		}
	}
	return true, 0
}

// FromStateString creates a new board from a state string representation
func FromStateString(state string) *Board {
	board := NewBoard()
	for i := 0; i < 9; i++ {
		row := i / 3
		col := i % 3
		board.state[row][col] = int(state[i] - '0')
	}
	return board
} 
