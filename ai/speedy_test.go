package ai_test

import (
	"testing"

	"github.com/Logiraptor/word-bot/ai"
	"github.com/Logiraptor/word-bot/core"
	"github.com/Logiraptor/word-bot/definitions"
	"github.com/Logiraptor/word-bot/wordlist"
)

var wordGaddag *wordlist.Gaddag

func init() {
	wordGaddag = wordlist.NewGaddag()
	err := definitions.LoadWords("../words.txt", wordGaddag)
	if err != nil {
		panic(err)
	}
}

func BenchmarkSpeedy(b *testing.B) {
	tiles := core.NewConsumableRack(core.MakeTiles(core.MakeWord("bdhrigs"), "xxxxxx "))
	board := core.NewBoard()
	speedy := ai.NewSpeedyAI(wordDB, wordDB)
	bag := core.NewConsumableBag()
	defer speedy.Kill()

	board.PlaceTiles(core.PlacedTiles{core.MakeTiles(core.MakeWord("aaaaaaaaaaaaaaa"), "xxxxxxxxxxxxxxx"), 0, 7, core.Vertical})
	board.PlaceTiles(core.PlacedTiles{core.MakeTiles(core.MakeWord("aaaaaaaaaaaaaa"), "xxxxxxxxxxxxxxx"), 7, 0, core.Horizontal})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		speedy.FindMove(board, bag, tiles, func(core.Turn) bool { return true })
	}
}

func BenchmarkSpeedySearch(b *testing.B) {
	rack := core.NewConsumableRack(core.MakeTiles(core.MakeWord("bdhrigs"), "xxxxxx "))
	board := core.NewBoard()
	speedy := ai.NewSpeedyAI(wordDB, wordDB)
	defer speedy.Kill()

	board.PlaceTiles(core.PlacedTiles{core.MakeTiles(core.MakeWord("aaaaaaaaaaaaaaa"), "xxxxxxxxxxxxxxx"), 0, 7, core.Vertical})
	board.PlaceTiles(core.PlacedTiles{core.MakeTiles(core.MakeWord("aaaaaaaaaaaaaa"), "xxxxxxxxxxxxxxx"), 7, 0, core.Horizontal})
	prev := []core.Tile{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		speedy.Search(board, 8, 8, core.Horizontal, rack, wordDB, prev, func([]core.Tile) {})
	}
}
