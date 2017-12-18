package ai_test

import (
	"testing"

	"github.com/Logiraptor/word-bot/ai"
	"github.com/Logiraptor/word-bot/core"
	"github.com/Logiraptor/word-bot/wordlist"
	"github.com/stretchr/testify/assert"
)

var wordDB *wordlist.Trie

func init() {
	wordDB = wordlist.MakeDefaultWordList()
}

func TestSmartyMatchesBrute(t *testing.T) {
	if testing.Short() {
		t.Skip("Brute is too slow")
	}
	tiles := core.NewConsumableRack(core.MakeTiles(core.MakeWord("asdjdha"), "xxxxxx "))
	board := core.NewBoard()

	board.PlaceTiles(core.PlacedTiles{core.MakeTiles(core.MakeWord("doggo"), "xxxxx"), 7, 7, core.Horizontal})
	board.PlaceTiles(core.PlacedTiles{core.MakeTiles(core.MakeWord("ar"), "xx"), 7, 8, core.Vertical})

	smarty := ai.NewSmartyAI(wordDB, wordDB)
	defer smarty.Kill()

	smartyMoves := []core.Turn{}
	smarty.GenerateMoves(board, tiles, func(t core.Turn) bool {
		smartyMoves = append(smartyMoves, t)
		return true
	})

	bruteMoves := []core.Turn{}
	ai.BruteForce(board, tiles, wordDB, func(t core.Turn) {
		bruteMoves = append(bruteMoves, t)
	})

	smartyMoves = unique(smartyMoves)
	bruteMoves = unique(bruteMoves)

	assert.Subset(t, bruteMoves, smartyMoves)
	if !assert.Equal(t, len(bruteMoves), len(smartyMoves)) {
		compareSets(board, "brute", "smarty", bruteMoves, smartyMoves)
		return
	}
}

func BenchmarkSmarty(b *testing.B) {
	tiles := core.NewConsumableRack(core.MakeTiles(core.MakeWord("bdhrigs"), "xxxxxx "))
	board := core.NewBoard()
	smarty := ai.NewSmartyAI(wordDB, wordDB)
	bag := core.NewConsumableBag()
	defer smarty.Kill()

	board.PlaceTiles(core.PlacedTiles{core.MakeTiles(core.MakeWord("aaaaaaaaaaaaaaa"), "xxxxxxxxxxxxxxx"), 0, 7, core.Vertical})
	board.PlaceTiles(core.PlacedTiles{core.MakeTiles(core.MakeWord("aaaaaaaaaaaaaa"), "xxxxxxxxxxxxxxx"), 7, 0, core.Horizontal})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		smarty.FindMove(board, bag, tiles, func(core.Turn) bool { return true })
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
