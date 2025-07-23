package tictactoe

import (
	"encoding/json"
	"errors"
	"fmt"
)

// GameState represents the state of a TicTacToe game
type GameState struct {
	// Board is a 3x3 grid represented as a 9-element array
	// 0-2: top row, 3-5: middle row, 6-8: bottom row
	// 0 = empty, 1 = X, 2 = O
	Board [9]int `json:"board"`

	// CurrentPlayerID is the player whose turn it is (0 or 1)
	CurrentPlayerID int `json:"currentPlayerID"`

	// PlayerSymbols maps player IDs to their symbols (1=X, 2=O)
	PlayerSymbols map[int]int `json:"playerSymbols"`

	// IsGameEnded is true if the game is finished
	IsGameEnded bool `json:"isGameEnded"`

	// WinnerPlayerID is the ID of the winning player, or -1 for draw/no winner
	WinnerPlayerID int `json:"winnerPlayerID"`

	// PossibleActions is a list of possible moves the current player can make
	PossibleActions []json.RawMessage `json:"possibleActions"`
}

// Player represents a player in the game
type Player struct {
	ID     int `json:"id"`
	Symbol int `json:"symbol"` // 1 = X, 2 = O
}

// Action represents a move in the game
type Action interface {
	IsPossible(g GameState) bool
	Run(g *GameState) error
	GetName() string
	GetPlayerID() int
}

// MoveAction represents placing a symbol on the board
type MoveAction struct {
	PlayerID int `json:"playerID"`
	Position int `json:"position"` // 0-8, position on the board
}

func (a MoveAction) GetName() string  { return "PLACE_SYMBOL" }
func (a MoveAction) GetPlayerID() int { return a.PlayerID }
func (a MoveAction) String() string {
	return fmt.Sprintf("Place symbol at position %d", a.Position)
}

func (a MoveAction) IsPossible(g GameState) bool {
	// Check if it's the player's turn
	if a.PlayerID != g.CurrentPlayerID {
		return false
	}
	// Check if position is valid and empty
	if a.Position < 0 || a.Position > 8 {
		return false
	}
	return g.Board[a.Position] == 0
}

func (a MoveAction) Run(g *GameState) error {
	if !a.IsPossible(*g) {
		return errors.New("move not possible")
	}

	// Place the symbol
	g.Board[a.Position] = g.PlayerSymbols[a.PlayerID]

	// Check for win or draw
	g.checkGameEnd()

	return nil
}

// NewTwoPlayerGame creates a new tictactoe game for two players
func NewTwoPlayerGame() *GameState {
	return &GameState{
		Board:           [9]int{},
		CurrentPlayerID: 0,
		PlayerSymbols: map[int]int{
			0: 1, // Player 0 is X
			1: 2, // Player 1 is O
		},
		IsGameEnded:     false,
		WinnerPlayerID:  -1,
		PossibleActions: []json.RawMessage{},
	}
}

// CalculatePossibleActions returns all possible moves for the current player
func (g *GameState) CalculatePossibleActions() []Action {
	if g.IsGameEnded {
		return []Action{}
	}

	var actions []Action
	for i := 0; i < 9; i++ {
		if g.Board[i] == 0 {
			action := MoveAction{
				PlayerID: g.CurrentPlayerID,
				Position: i,
			}
			actions = append(actions, action)
		}
	}
	return actions
}

// RunAction executes an action and updates the game state
func (g *GameState) RunAction(action Action) error {
	if action == nil {
		return nil
	}

	if g.IsGameEnded {
		return errors.New("game is already ended")
	}

	if !action.IsPossible(*g) {
		return errors.New("action not possible")
	}

	err := action.Run(g)
	if err != nil {
		return err
	}

	// Switch turns if game is not ended (in TicTacToe, every move yields turn)
	if !g.IsGameEnded {
		g.CurrentPlayerID = 1 - g.CurrentPlayerID
	}

	// Update possible actions
	g.PossibleActions = serializeActions(g.CalculatePossibleActions())

	return nil
}

// checkGameEnd checks if the game has ended (win or draw)
func (g *GameState) checkGameEnd() {
	// Check for wins
	winPatterns := [][3]int{
		// Rows
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8},
		// Columns
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8},
		// Diagonals
		{0, 4, 8}, {2, 4, 6},
	}

	for _, pattern := range winPatterns {
		if g.Board[pattern[0]] != 0 &&
			g.Board[pattern[0]] == g.Board[pattern[1]] &&
			g.Board[pattern[1]] == g.Board[pattern[2]] {
			g.IsGameEnded = true
			// Find which player has this symbol
			for playerID, symbol := range g.PlayerSymbols {
				if symbol == g.Board[pattern[0]] {
					g.WinnerPlayerID = playerID
					return
				}
			}
		}
	}

	// Check for draw (board full)
	full := true
	for _, cell := range g.Board {
		if cell == 0 {
			full = false
			break
		}
	}

	if full {
		g.IsGameEnded = true
		g.WinnerPlayerID = -1 // Draw
	}
}

// ToClientGameState converts the game state to a client-friendly format
func (g *GameState) ToClientGameState(playerID int) ClientGameState {
	opponentID := 1 - playerID

	return ClientGameState{
		Board:           g.Board,
		CurrentPlayerID: g.CurrentPlayerID,
		IsGameEnded:     g.IsGameEnded,
		WinnerPlayerID:  g.WinnerPlayerID,
		PossibleActions: g.PossibleActions,
		YouPlayerID:     playerID,
		ThemPlayerID:    opponentID,
		YourSymbol:      g.PlayerSymbols[playerID],
		TheirSymbol:     g.PlayerSymbols[opponentID],
	}
}

// ClientGameState represents the game state as seen by a client
type ClientGameState struct {
	Board           [9]int            `json:"board"`
	CurrentPlayerID int               `json:"currentPlayerID"`
	IsGameEnded     bool              `json:"isGameEnded"`
	WinnerPlayerID  int               `json:"winnerPlayerID"`
	PossibleActions []json.RawMessage `json:"possibleActions"`
	YouPlayerID     int               `json:"youPlayerID"`
	ThemPlayerID    int               `json:"themPlayerID"`
	YourSymbol      int               `json:"yourSymbol"`
	TheirSymbol     int               `json:"theirSymbol"`
}

// Serialize converts the game state to JSON
func (g *GameState) Serialize() ([]byte, error) {
	return json.Marshal(g)
}

// serializeActions converts actions to JSON for transmission
func serializeActions(actions []Action) []json.RawMessage {
	var serialized []json.RawMessage
	for _, action := range actions {
		if bs, err := json.Marshal(action); err == nil {
			serialized = append(serialized, bs)
		}
	}
	return serialized
}

// DeserializeAction converts JSON back to an Action
func DeserializeAction(data []byte) (Action, error) {
	var moveAction MoveAction
	if err := json.Unmarshal(data, &moveAction); err == nil {
		return moveAction, nil
	}

	return nil, errors.New("unknown action type")
}

// Bot interface for implementing tictactoe bots
type Bot interface {
	ChooseAction(ClientGameState) Action
}

// Initialize the game state with possible actions
func (g *GameState) Initialize() {
	g.PossibleActions = serializeActions(g.CalculatePossibleActions())
}

// New creates a new game and initializes it
func New() *GameState {
	game := NewTwoPlayerGame()
	game.Initialize()
	return game
}
