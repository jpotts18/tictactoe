package agent

import (
	"math"
	"math/rand"

	"github.com/jpotts18/tictactoe/game"
)

type QAgent struct {
	BaseAgent
	qTable     map[string][]float64
	epsilon    float64
	alpha      float64
	gamma      float64
	Evaluating bool
}

func NewQAgent(player int) *QAgent {
	return &QAgent{
		BaseAgent:  BaseAgent{Player: player},
		qTable:     make(map[string][]float64),
		epsilon:    0.9,
		alpha:      0.1,
		gamma:      0.99,
		Evaluating: false,
	}
}

func (q *QAgent) GetQValues(state string) []float64 {
	if _, exists := q.qTable[state]; !exists {
		q.qTable[state] = make([]float64, 9)
	}
	return q.qTable[state]
}

func (q *QAgent) GetMove(board *game.Board, player int) int {
	moves := board.GetAvailableMoves()
	if len(moves) == 0 {
		return -1
	}

	stateKey := q.GetStateKey(board)
	qValues := q.GetQValues(stateKey)

	if q.Evaluating {
		return q.getBestAction(moves, qValues)
	}

	if rand.Float64() < q.epsilon {
		return moves[rand.Intn(len(moves))]
	}

	return q.getBestAction(moves, qValues)
}

func (q *QAgent) getBestAction(moves []int, qValues []float64) int {
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

func (q *QAgent) Learn(state string, action int, reward float64, newState string) {
	oldQValues := q.GetQValues(state)
	oldValue := oldQValues[action]

	var maxNextQ float64
	if newState != "" {
		nextQValues := q.GetQValues(newState)
		nextBoard := game.FromStateString(newState)
		nextMoves := nextBoard.GetAvailableMoves()
		if len(nextMoves) > 0 {
			maxNextQ = nextQValues[q.getBestAction(nextMoves, nextQValues)]
		}
	}

	newValue := oldValue + q.alpha*(reward + q.gamma*maxNextQ - oldValue)
	oldQValues[action] = newValue
	q.qTable[state] = oldQValues
	
	q.epsilon = math.Max(0.1, q.epsilon*0.99995)
}

func (q *QAgent) Save(filename string) error {
	return SaveQTable(filename + ".qlearning", q.qTable)
}

func (q *QAgent) Load(filename string) error {
	qtable, err := LoadQTable(filename + ".qlearning")
	if err != nil {
		return err
	}
	q.qTable = qtable
	return nil
}
