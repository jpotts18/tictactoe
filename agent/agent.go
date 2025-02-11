package agent

import (
	"strconv"

	"github.com/jpotts18/tictactoe/game"
)

// Agent defines the interface that all game-playing agents must implement
type Agent interface {
	// GetMove selects the next move given the current board state and player
	GetMove(board *game.Board, player int) int
	
	// Learn updates the agent's knowledge based on the game experience
	Learn(oldState string, action int, reward float64, newState string)
	
	// GetStateKey converts a board state into a string representation
	GetStateKey(board *game.Board) string
}

// BaseAgent provides common functionality for all agents
type BaseAgent struct {
	Player int
}

func (b *BaseAgent) GetStateKey(board *game.Board) string {
	state := ""
	for i := 0; i < 9; i++ {
		state += strconv.Itoa(board.GetCell(i))
	}
	return state
}


