package agent

import (
	"math"
	"math/rand"

	"github.com/jpotts18/tictactoe/game"
)

type SarsaAgent struct {
	BaseAgent
	qTable     map[string][]float64
	epsilon    float64
	alpha      float64
	gamma      float64
	lastState  string
	lastAction int
	Evaluating bool
}

func NewSarsaAgent(player int) *SarsaAgent {
	return &SarsaAgent{
		BaseAgent:  BaseAgent{Player: player},
		qTable:     make(map[string][]float64),
		epsilon:    0.9,
		alpha:      0.1,
		gamma:      0.99,
		lastState:  "",
		lastAction: -1,
		Evaluating: false,
	}
}

func (s *SarsaAgent) GetQValues(state string) []float64 {
	if _, exists := s.qTable[state]; !exists {
		s.qTable[state] = make([]float64, 9)
	}
	return s.qTable[state]
}

func (s *SarsaAgent) GetMove(board *game.Board, player int) int {
	state := s.GetStateKey(board)
	moves := board.GetEmptyCells()
	
	if len(moves) == 0 {
		return -1
	}

	if rand.Float64() < s.epsilon {
		return moves[rand.Intn(len(moves))]
	}

	return s.getBestAction(state, moves)
}

func (s *SarsaAgent) getBestAction(state string, moves []int) int {
	qValues := s.GetQValues(state)
	bestMove := moves[0]
	bestValue := qValues[bestMove]

	for _, move := range moves[1:] {
		if qValues[move] > bestValue {
			bestMove = move
			bestValue = qValues[move]
		}
	}
	return bestMove
}

func (s *SarsaAgent) Learn(state string, action int, reward float64, newState string) {
	if s.lastState != "" {
		currentQValues := s.GetQValues(newState)
		nextAction := -1
		
		if newState != "" {
			boardState := game.FromStateString(newState)
			moves := boardState.GetAvailableMoves()
			if len(moves) > 0 {
				if rand.Float64() < s.epsilon {
					nextAction = moves[rand.Intn(len(moves))]
				} else {
					nextAction = s.getBestAction(state, moves)
				}
			}
		}

		oldQValues := s.GetQValues(s.lastState)
		oldValue := oldQValues[s.lastAction]
		
		var nextQValue float64
		if nextAction != -1 {
			nextQValue = currentQValues[nextAction]
		}

		newValue := oldValue + s.alpha*(reward + s.gamma*nextQValue - oldValue)
		oldQValues[s.lastAction] = newValue
		s.qTable[s.lastState] = oldQValues
	}

	s.lastState = state
	s.lastAction = action
	s.epsilon = math.Max(0.1, s.epsilon*0.99995)
}
