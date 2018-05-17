package main

import (
	"fmt"

	"github.com/Logiraptor/word-bot/ai"
	"github.com/Logiraptor/word-bot/core"
	"github.com/Logiraptor/word-bot/wordlist"
	"github.com/pkg/profile"
)

var wordDB = wordlist.MakeDefaultWordList()

func main() {
	defer profile.Start(profile.TraceProfile).Stop()
	g := ai.PlayGame(wordDB, p1, p2)
	for _, move := range g.Moves {
		fmt.Println(move.Score, move.Player, move.Tiles)
	}
}

func p2(board *core.Board) *ai.Player {
	return ai.NewPlayer(ai.NewSmartyAI(wordDB, wordDB))
}

func p1(board *core.Board) *ai.Player {
	return ai.NewPlayer(ai.NewSmartyAI(wordDB, wordDB))
}
