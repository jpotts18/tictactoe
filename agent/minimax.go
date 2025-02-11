package agent

import "github.com/jpotts18/tictactoe/game"

// MinimaxAgent implements a perfect player using the minimax algorithm
type MinimaxAgent struct {
	BaseAgent
}

func NewMinimaxAgent(player int) *MinimaxAgent {
	return &MinimaxAgent{
		BaseAgent: BaseAgent{Player: player},
	}
}

// ChooseAction evaluates all possible moves and selects the optimal one
func (m *MinimaxAgent) ChooseAction(board *game.Board) int {
	moves := board.GetAvailableMoves()
	if len(moves) == 0 {
		return -1
	}
	
	bestScore := -1000
	bestMove := moves[0]
	
	// Try each possible move and evaluate the resulting position
	for _, move := range moves {
		boardCopy := *board
		boardCopy.MakeMove(move, m.Player)
		// Evaluate position assuming opponent plays optimally
		score := m.minimax(&boardCopy, false, m.Player%2+1, 5) // depth of 5 moves ahead
		if score > bestScore {
			bestScore = score
			bestMove = move
		}
	}
	return bestMove
}

// minimax implements the recursive minimax algorithm
func (m *MinimaxAgent) minimax(board *game.Board, isMax bool, currentPlayer int, depth int) int {
	gameOver, winner := board.IsGameOver()
	
	// Base cases: game is over or reached depth limit
	if gameOver || depth == 0 {
		if winner == m.Player {
			return 1  // Win for this agent
		} else if winner == 0 {
			return 0  // Draw
		}
		return -1    // Loss for this agent
	}
	
	moves := board.GetAvailableMoves()
	if isMax {
		// Maximizing player: find the maximum score among all possible moves
		bestScore := -1000
		for _, move := range moves {
			boardCopy := *board
			boardCopy.MakeMove(move, m.Player)
			score := m.minimax(&boardCopy, false, currentPlayer%2+1, depth-1)
			if score > bestScore {
				bestScore = score
			}
		}
		return bestScore
	} else {
		// Minimizing player: find the minimum score among all possible moves
		bestScore := 1000
		for _, move := range moves {
			boardCopy := *board
			boardCopy.MakeMove(move, currentPlayer)
			score := m.minimax(&boardCopy, true, currentPlayer%2+1, depth-1)
			if score < bestScore {
				bestScore = score
			}
		}
		return bestScore
	}
}

// Learn is a no-op for MinimaxAgent as it doesn't learn
func (m *MinimaxAgent) Learn(oldState string, action int, reward float64, newState string) {
	// Minimax agent doesn't learn
}

func (m *MinimaxAgent) GetMove(board *game.Board, player int) int {
	bestScore := -1000
	bestMove := -1
	
	for _, move := range board.GetEmptyCells() {
		boardCopy := *board
		boardCopy.MakeMove(move, player)
		score := m.minimax(&boardCopy, false, player%2+1, 5)
		
		if score > bestScore {
			bestScore = score
			bestMove = move
		}
	}
	
	return bestMove
} 
