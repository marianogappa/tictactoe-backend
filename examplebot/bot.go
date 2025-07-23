package examplebot

import (
	"math/rand"

	"github.com/marianogappa/tictactoe/tictactoe"
)

// SimpleBot implements a basic tictactoe strategy
type SimpleBot struct{}

// New creates a new SimpleBot
func New() *SimpleBot {
	return &SimpleBot{}
}

// ChooseAction implements the Bot interface
func (b *SimpleBot) ChooseAction(state tictactoe.ClientGameState) tictactoe.Action {
	if state.IsGameEnded {
		return nil
	}

	// Deserialize possible actions
	var actions []tictactoe.Action
	for _, rawAction := range state.PossibleActions {
		if action, err := tictactoe.DeserializeAction(rawAction); err == nil {
			actions = append(actions, action)
		}
	}

	if len(actions) == 0 {
		return nil
	}

	// Strategy:
	// 1. Win if possible
	// 2. Block opponent from winning
	// 3. Take center if available
	// 4. Take corners
	// 5. Take random available spot

	mySymbol := state.YourSymbol
	theirSymbol := state.TheirSymbol

	// Check for winning move
	for _, action := range actions {
		if moveAction, ok := action.(tictactoe.MoveAction); ok {
			if b.isWinningMove(state.Board, moveAction.Position, mySymbol) {
				return moveAction
			}
		}
	}

	// Check for blocking move
	for _, action := range actions {
		if moveAction, ok := action.(tictactoe.MoveAction); ok {
			if b.isWinningMove(state.Board, moveAction.Position, theirSymbol) {
				return moveAction
			}
		}
	}

	// Take center if available (position 4)
	for _, action := range actions {
		if moveAction, ok := action.(tictactoe.MoveAction); ok {
			if moveAction.Position == 4 {
				return moveAction
			}
		}
	}

	// Take corners (positions 0, 2, 6, 8)
	corners := []int{0, 2, 6, 8}
	for _, corner := range corners {
		for _, action := range actions {
			if moveAction, ok := action.(tictactoe.MoveAction); ok {
				if moveAction.Position == corner {
					return moveAction
				}
			}
		}
	}

	// Take any available position
	if len(actions) > 0 {
		return actions[rand.Intn(len(actions))]
	}

	return nil
}

// isWinningMove checks if placing a symbol at the given position would result in a win
func (b *SimpleBot) isWinningMove(board [9]int, position int, symbol int) bool {
	// Create a copy of the board and simulate the move
	testBoard := board
	testBoard[position] = symbol

	// Check all win patterns
	winPatterns := [][3]int{
		// Rows
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8},
		// Columns
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8},
		// Diagonals
		{0, 4, 8}, {2, 4, 6},
	}

	for _, pattern := range winPatterns {
		if testBoard[pattern[0]] == symbol &&
			testBoard[pattern[1]] == symbol &&
			testBoard[pattern[2]] == symbol {
			return true
		}
	}

	return false
}
