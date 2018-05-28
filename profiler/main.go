package main

import (
	"fmt"
	"time"

	"github.com/Logiraptor/word-bot/ai"
	"github.com/Logiraptor/word-bot/core"
	"github.com/Logiraptor/word-bot/wordlist"
	"github.com/pkg/profile"
)

var wordDB = wordlist.MakeDefaultWordList()

func main() {
	defer profile.Start(profile.CPUProfile).Stop()
	start := time.Now()
	for time.Since(start) < time.Second*10 {
		g := ai.PlayGame(wordDB, speedy, speedy)
		for _, move := range g.Moves {
			fmt.Println(move.Score, move.Player, move.Tiles)
		}
	}
}

func smarty(board *core.Board) *ai.Player {
	return ai.NewPlayer(ai.NewSmartyAI(wordDB, wordDB))
}

func speedy(board *core.Board) *ai.Player {
	return ai.NewPlayer(ai.NewSmartyAI(wordDB, wordDB))
}

func brute(board *core.Board) *ai.Player {
	return ai.NewPlayer(ai.NewBrute(wordDB))
}
