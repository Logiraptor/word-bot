package ai_test

import (
	"testing"
	"word-bot/ai"
	"word-bot/core"
	"word-bot/definitions"
	"word-bot/wordlist"
)

var wordDB *wordlist.Trie

func init() {
	words, err := definitions.LoadWords("../words.txt")
	if err != nil {
		panic(err)
	}

	// wordDB = wordlist.NewTrie()
	// for _, word := range words {
	// 	wordDB.AddWord(word)
	// }

	builder := wordlist.NewTrieBuilder()
	for _, word := range words {
		builder.AddWord(word)
	}
	wordDB = builder.Build()
}

func BenchmarkSmarty(b *testing.B) {
	tiles := core.MakeTiles(core.MakeWord("bdhrigs"), "xxxxxx ")
	board := core.NewBoard()
	smarty := ai.NewSmartyAI(board, wordDB, wordDB)

	board.PlaceTiles(core.MakeTiles(core.MakeWord("aaaaaaaaaaaaaaa"), "xxxxxxxxxxxxxxx"), 0, 7, core.Vertical)
	board.PlaceTiles(core.MakeTiles(core.MakeWord("aaaaaaaaaaaaaa"), "xxxxxxxxxxxxxxx"), 7, 0, core.Horizontal)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		smarty.FindMoves(tiles)
	}
}

func BenchmarkSearch(b *testing.B) {
	rack := core.NewConsumableRack(core.MakeTiles(core.MakeWord("bdhrigs"), "xxxxxx "))
	board := core.NewBoard()
	smarty := ai.NewSmartyAI(board, wordDB, wordDB)

	board.PlaceTiles(core.MakeTiles(core.MakeWord("aaaaaaaaaaaaaaa"), "xxxxxxxxxxxxxxx"), 0, 7, core.Vertical)
	board.PlaceTiles(core.MakeTiles(core.MakeWord("aaaaaaaaaaaaaa"), "xxxxxxxxxxxxxxx"), 7, 0, core.Horizontal)
	prev := []core.Tile{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		smarty.Search(8, 8, core.Horizontal, rack, wordDB, prev, func([]core.Tile) {})
	}
}
