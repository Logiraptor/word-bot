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

	wordDB = wordlist.NewTrie()
	for _, word := range words {
		wordDB.AddWord(word)
	}
}

func BenchmarkSmarty(b *testing.B) {
	tiles := core.MakeTiles(core.MakeWord("bdhrigs"), "xxxxxxx")
	for i := 0; i < b.N; i++ {
		b := core.NewBoard()
		ai := ai.NewSmartyAI(b, wordDB, wordDB)
		ai.FindMoves(tiles)
	}
}
