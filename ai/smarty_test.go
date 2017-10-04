package ai_test

import (
	"testing"

	"github.com/Logiraptor/word-bot/ai"
	"github.com/Logiraptor/word-bot/core"
	"github.com/Logiraptor/word-bot/definitions"
	"github.com/Logiraptor/word-bot/wordlist"
)

var wordDB *wordlist.Trie

func init() {
	wordDB = wordlist.NewTrie()
	err := definitions.LoadWords("../shortwords.txt", wordDB)
	if err != nil {
		panic(err)
	}
}

func BenchmarkSmarty(b *testing.B) {
	tiles := core.MakeTiles(core.MakeWord("bdhrigs"), "xxxxxx ")
	board := core.NewBoard()
	smarty := ai.NewSmartyAI(board, wordDB, wordDB)

	board.PlaceTiles(core.PlacedTiles{core.MakeTiles(core.MakeWord("aaaaaaaaaaaaaaa"), "xxxxxxxxxxxxxxx"), 0, 7, core.Vertical})
	board.PlaceTiles(core.PlacedTiles{core.MakeTiles(core.MakeWord("aaaaaaaaaaaaaa"), "xxxxxxxxxxxxxxx"), 7, 0, core.Horizontal})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		smarty.FindMoves(tiles)
	}
}

func BenchmarkSearch(b *testing.B) {
	rack := core.NewConsumableRack(core.MakeTiles(core.MakeWord("bdhrigs"), "xxxxxx "))
	board := core.NewBoard()
	smarty := ai.NewSmartyAI(board, wordDB, wordDB)

	board.PlaceTiles(core.PlacedTiles{core.MakeTiles(core.MakeWord("aaaaaaaaaaaaaaa"), "xxxxxxxxxxxxxxx"), 0, 7, core.Vertical})
	board.PlaceTiles(core.PlacedTiles{core.MakeTiles(core.MakeWord("aaaaaaaaaaaaaa"), "xxxxxxxxxxxxxxx"), 7, 0, core.Horizontal})
	prev := []core.Tile{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		smarty.Search(8, 8, core.Horizontal, rack, wordDB, prev, func([]core.Tile) {})
	}
}
