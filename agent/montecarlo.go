package agent

import (
	"math"
	"math/rand"

	"github.com/jpotts18/tictactoe/game"
)

type Episode struct {
	state  string
	action int
	reward float64
}

type MonteCarloAgent struct {
	BaseAgent
	qTable     map[string][]float64
	returns    map[string]map[int][]float64
	epsilon    float64
	gamma      float64
	episode    []Episode
	Evaluating bool
}


func NewMonteCarloAgent(player int) *MonteCarloAgent {
	return &MonteCarloAgent{
		BaseAgent:  BaseAgent{Player: player},
		qTable:     make(map[string][]float64),
		returns:    make(map[string]map[int][]float64),
		epsilon:    0.9,
		gamma:      0.99,
		episode:    make([]Episode, 0),
		Evaluating: false,
	}
}

func (m *MonteCarloAgent) GetQValues(state string) []float64 {
	if _, exists := m.qTable[state]; !exists {
		m.qTable[state] = make([]float64, 9)
	}
	return m.qTable[state]
}

func (m *MonteCarloAgent) GetMove(board *game.Board, player int) int {
	state := m.GetStateKey(board)
	moves := board.GetEmptyCells()
	
	if len(moves) == 0 {
		return -1
	}

	if rand.Float64() < m.epsilon {
		return moves[rand.Intn(len(moves))]
	}

	return m.getBestAction(state, moves)
}

func (m *MonteCarloAgent) getBestAction(state string, moves []int) int {
	qValues := m.GetQValues(state)
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

func (m *MonteCarloAgent) Learn(state string, action int, reward float64, newState string) {
	m.episode = append(m.episode, Episode{state, action, reward})

	// Only update at the end of the episode
	if newState == "" {
		m.updateEpisode()
		m.episode = make([]Episode, 0)
		m.epsilon = math.Max(0.1, m.epsilon*0.99995)
	}
}

func (m *MonteCarloAgent) updateEpisode() {
	// Track first-visit for each state-action pair in this episode
	visited := make(map[string]map[int]bool)
	
	// Calculate returns for each state-action pair
	G := 0.0
	for i := len(m.episode) - 1; i >= 0; i-- {
		exp := m.episode[i]
		G = m.gamma*G + exp.reward
		
		// Only update on first visit to each state-action pair
		if _, exists := visited[exp.state]; !exists {
			visited[exp.state] = make(map[int]bool)
		}
		if !visited[exp.state][exp.action] {
			visited[exp.state][exp.action] = true
			
			if _, exists := m.returns[exp.state]; !exists {
				m.returns[exp.state] = make(map[int][]float64)
			}
			if _, exists := m.returns[exp.state][exp.action]; !exists {
				m.returns[exp.state][exp.action] = make([]float64, 0)
			}
			
			// Add this return to our running history
			m.returns[exp.state][exp.action] = append(m.returns[exp.state][exp.action], G)
			
			// Update Q-value with new average
			sum := 0.0
			for _, r := range m.returns[exp.state][exp.action] {
				sum += r
			}
			average := sum / float64(len(m.returns[exp.state][exp.action]))
			
			qValues := m.GetQValues(exp.state)
			qValues[exp.action] = average
			m.qTable[exp.state] = qValues
		}
	}
} 
