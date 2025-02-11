package agent

import (
	"math/rand"

	"github.com/jpotts18/tictactoe/game"
)

type RandomAgent struct {
	BaseAgent
	rng *rand.Rand
}

func NewRandomAgent(player int) *RandomAgent {
	return &RandomAgent{
		BaseAgent: BaseAgent{Player: player},
		rng:       rand.New(rand.NewSource(42)),
	}
}

func (r *RandomAgent) GetMove(board *game.Board, player int) int {
	moves := board.GetAvailableMoves()
	if len(moves) == 0 {
		return -1
	}
	return moves[r.rng.Intn(len(moves))]
}

func (r *RandomAgent) Learn(oldState string, action int, reward float64, newState string) {
	// Random agent doesn't learn
} 
