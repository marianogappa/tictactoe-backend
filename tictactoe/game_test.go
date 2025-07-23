package tictactoe

import (
	"testing"
)

func TestNewGame(t *testing.T) {
	game := New()

	if game == nil {
		t.Fatal("New() returned nil")
	}

	if game.IsGameEnded {
		t.Error("New game should not be ended")
	}

	if game.CurrentPlayerID != 0 {
		t.Error("New game should start with player 0")
	}

	// Check that board is empty
	for i, cell := range game.Board {
		if cell != 0 {
			t.Errorf("Cell %d should be empty, got %d", i, cell)
		}
	}
}

func TestMoveAction(t *testing.T) {
	game := New()

	// Test valid move
	action := MoveAction{PlayerID: 0, Position: 0}
	if !action.IsPossible(*game) {
		t.Error("Valid move should be possible")
	}

	err := game.RunAction(action)
	if err != nil {
		t.Errorf("Valid move should not return error: %v", err)
	}

	// Check that symbol was placed
	if game.Board[0] != 1 { // Player 0 should be symbol 1 (X)
		t.Errorf("Symbol should be placed at position 0, got %d", game.Board[0])
	}

	// Check that turn switched
	if game.CurrentPlayerID != 1 {
		t.Error("Turn should switch to player 1")
	}

	// Test invalid move (position already taken)
	action2 := MoveAction{PlayerID: 1, Position: 0}
	if action2.IsPossible(*game) {
		t.Error("Move to occupied position should not be possible")
	}
}

func TestWinCondition(t *testing.T) {
	game := New()

	// Create a winning scenario for player 0 (horizontal line)
	moves := []MoveAction{
		{PlayerID: 0, Position: 0}, // X at 0
		{PlayerID: 1, Position: 3}, // O at 3
		{PlayerID: 0, Position: 1}, // X at 1
		{PlayerID: 1, Position: 4}, // O at 4
		{PlayerID: 0, Position: 2}, // X at 2 (winning move)
	}

	for i, move := range moves {
		err := game.RunAction(move)
		if err != nil {
			t.Errorf("Move %d should not return error: %v", i, err)
		}

		// Check if game ends on the last move
		if i == len(moves)-1 {
			if !game.IsGameEnded {
				t.Error("Game should be ended after winning move")
			}
			if game.WinnerPlayerID != 0 {
				t.Errorf("Player 0 should be winner, got %d", game.WinnerPlayerID)
			}
		}
	}
}

func TestDrawCondition(t *testing.T) {
	game := New()

	// Create a true draw scenario
	// Final board:
	// X | O | X
	// X | O | O
	// O | X | X
	moves := []MoveAction{
		{PlayerID: 0, Position: 0}, // X at 0
		{PlayerID: 1, Position: 1}, // O at 1
		{PlayerID: 0, Position: 2}, // X at 2
		{PlayerID: 1, Position: 4}, // O at 4
		{PlayerID: 0, Position: 3}, // X at 3
		{PlayerID: 1, Position: 5}, // O at 5
		{PlayerID: 0, Position: 7}, // X at 7
		{PlayerID: 1, Position: 6}, // O at 6
		{PlayerID: 0, Position: 8}, // X at 8
	}

	for _, move := range moves {
		err := game.RunAction(move)
		if err != nil {
			t.Errorf("Move should not return error: %v", err)
		}
	}

	if !game.IsGameEnded {
		t.Error("Game should be ended (draw)")
	}
	if game.WinnerPlayerID != -1 {
		t.Errorf("Should be a draw (winner -1), got %d", game.WinnerPlayerID)
	}
}

func TestToClientGameState(t *testing.T) {
	game := New()

	clientState := game.ToClientGameState(0)

	if clientState.YouPlayerID != 0 {
		t.Errorf("YouPlayerID should be 0, got %d", clientState.YouPlayerID)
	}
	if clientState.ThemPlayerID != 1 {
		t.Errorf("ThemPlayerID should be 1, got %d", clientState.ThemPlayerID)
	}
	if clientState.YourSymbol != 1 {
		t.Errorf("YourSymbol should be 1 (X), got %d", clientState.YourSymbol)
	}
	if clientState.TheirSymbol != 2 {
		t.Errorf("TheirSymbol should be 2 (O), got %d", clientState.TheirSymbol)
	}
}
