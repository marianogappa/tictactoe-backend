//go:build tinygo
// +build tinygo

package main

import (
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/marianogappa/tictactoe/examplebot"
	"github.com/marianogappa/tictactoe/tictactoe"
)

func main() {
	js.Global().Set("tictactoeNew", js.FuncOf(tictactoeNew))
	js.Global().Set("tictactoeRunAction", js.FuncOf(tictactoeRunAction))
	js.Global().Set("tictactoeBotRunAction", js.FuncOf(tictactoeBotRunAction))
	select {}
}

var (
	state *tictactoe.GameState
	bot   tictactoe.Bot
)

func tictactoeNew(this js.Value, p []js.Value) interface{} {
	state = tictactoe.New()
	bot = examplebot.New()

	nbs, err := json.Marshal(state.ToClientGameState(0))
	if err != nil {
		panic(err)
	}

	buffer := js.Global().Get("Uint8Array").New(len(nbs))
	js.CopyBytesToJS(buffer, nbs)
	return buffer
}

func tictactoeRunAction(this js.Value, p []js.Value) interface{} {
	jsonBytes := make([]byte, p[0].Length())
	js.CopyBytesToGo(jsonBytes, p[0])

	newBytes := _runAction(jsonBytes)

	buffer := js.Global().Get("Uint8Array").New(len(newBytes))
	js.CopyBytesToJS(buffer, newBytes)
	return buffer
}

func tictactoeBotRunAction(this js.Value, p []js.Value) interface{} {
	if !state.IsGameEnded {
		action := bot.ChooseAction(state.ToClientGameState(1))

		err := state.RunAction(action)
		if err != nil {
			panic(fmt.Errorf("running action: %w", err))
		}
	}

	nbs, err := json.Marshal(state.ToClientGameState(0))
	if err != nil {
		panic(fmt.Errorf("marshalling game state: %w", err))
	}

	buffer := js.Global().Get("Uint8Array").New(len(nbs))
	js.CopyBytesToJS(buffer, nbs)
	return buffer
}

func _runAction(bs []byte) []byte {
	action, err := tictactoe.DeserializeAction(bs)
	if err != nil {
		panic(err)
	}
	err = state.RunAction(action)
	if err != nil {
		panic(err)
	}
	nbs, err := json.Marshal(state.ToClientGameState(0))
	if err != nil {
		panic(err)
	}
	return nbs
}
