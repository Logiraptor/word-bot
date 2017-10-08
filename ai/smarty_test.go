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
	tiles := core.NewConsumableRack(core.MakeTiles(core.MakeWord("bdhrigs"), "xxxxxx "))
	board := core.NewBoard()
	smarty := ai.NewSmartyAI(wordDB, wordDB)
	defer smarty.Kill()

	board.PlaceTiles(core.PlacedTiles{core.MakeTiles(core.MakeWord("aaaaaaaaaaaaaaa"), "xxxxxxxxxxxxxxx"), 0, 7, core.Vertical})
	board.PlaceTiles(core.PlacedTiles{core.MakeTiles(core.MakeWord("aaaaaaaaaaaaaa"), "xxxxxxxxxxxxxxx"), 7, 0, core.Horizontal})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		smarty.FindMove(board, tiles, func(core.Turn) bool { return true })
	}
}

func BenchmarkSearch(b *testing.B) {
	rack := core.NewConsumableRack(core.MakeTiles(core.MakeWord("bdhrigs"), "xxxxxx "))
	board := core.NewBoard()
	smarty := ai.NewSmartyAI(wordDB, wordDB)
	defer smarty.Kill()

	board.PlaceTiles(core.PlacedTiles{core.MakeTiles(core.MakeWord("aaaaaaaaaaaaaaa"), "xxxxxxxxxxxxxxx"), 0, 7, core.Vertical})
	board.PlaceTiles(core.PlacedTiles{core.MakeTiles(core.MakeWord("aaaaaaaaaaaaaa"), "xxxxxxxxxxxxxxx"), 7, 0, core.Horizontal})
	prev := []core.Tile{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		smarty.Search(board, 8, 8, core.Horizontal, rack, wordDB, prev, func([]core.Tile) {})
	}
}
