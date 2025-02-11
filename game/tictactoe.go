package game

import (
	"errors"
	"fmt"
)

type TicTacToe struct {
	board         *Board
	currentPlayer int
	isGameOver    bool
	winner        int
}

func NewTicTacToe() *TicTacToe {
	return &TicTacToe{
		board:         NewBoard(),
		currentPlayer: 1,
		isGameOver:    false,
		winner:        0,
	}
}

func (g *TicTacToe) MakeMove(position int) error {
	if g.isGameOver {
		return errors.New("game is already over")
	}

	if !g.board.SetCell(position, g.currentPlayer) {
		return errors.New("invalid move")
	}

	if g.checkWinner() {
		g.isGameOver = true
		g.winner = g.currentPlayer
	} else if len(g.board.GetEmptyCells()) == 0 {
		g.isGameOver = true
		g.winner = 0 // Draw
	} else {
		g.switchPlayer()
	}

	return nil
}

func (g *TicTacToe) checkWinner() bool {
	// Check rows
	for i := 0; i < 3; i++ {
		if g.checkLine(i*3, i*3+1, i*3+2) {
			return true
		}
	}

	// Check columns
	for i := 0; i < 3; i++ {
		if g.checkLine(i, i+3, i+6) {
			return true
		}
	}

	// Check diagonals
	if g.checkLine(0, 4, 8) {
		return true
	}
	if g.checkLine(2, 4, 6) {
		return true
	}

	return false
}

func (g *TicTacToe) checkLine(pos1, pos2, pos3 int) bool {
	p1 := g.board.GetCell(pos1)
	return p1 != 0 && p1 == g.board.GetCell(pos2) && p1 == g.board.GetCell(pos3)
}

func (g *TicTacToe) switchPlayer() {
	if g.currentPlayer == 1 {
		g.currentPlayer = 2
	} else {
		g.currentPlayer = 1
	}
}

func (g *TicTacToe) GetCurrentPlayer() int {
	return g.currentPlayer
}

func (g *TicTacToe) IsGameOver() bool {
	return g.isGameOver
}

func (g *TicTacToe) GetWinner() int {
	return g.winner
}

func (g *TicTacToe) GetBoard() *Board {
	return g.board
}

func (g *TicTacToe) GetAvailableMoves() []int {
	return g.board.GetEmptyCells()
}

func (g *TicTacToe) String() string {
	status := fmt.Sprintf("Current player: %d\n", g.currentPlayer)
	if g.isGameOver {
		if g.winner == 0 {
			status += "Game Over: Draw!\n"
		} else {
			status += fmt.Sprintf("Game Over: Player %d wins!\n", g.winner)
		}
	}
	status += g.board.String()
	return status
} 
