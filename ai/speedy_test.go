package ai_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

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

func filter(scoredMoves []core.Turn, pred func(core.ScoredMove) bool) []core.ScoredMove {
	output := []core.ScoredMove{}
	for i := range scoredMoves {
		if m, ok := scoredMoves[i].(core.ScoredMove); ok {
			if pred(m) {
				output = append(output, m)
			}
		}
	}
	return output
}

func compareSets(board *core.Board, a, b []core.Turn) {
	hasBlank := func(t []core.Tile) bool {
		for _, x := range t {
			if x.IsBlank() {
				return true
			}
		}
		return false
	}
	vertical := func(c core.ScoredMove) bool {
		return c.Direction == core.Vertical
	}
	horizontal := func(c core.ScoredMove) bool {
		return c.Direction == core.Horizontal
	}
	withBlank := func(c core.ScoredMove) bool {
		return hasBlank(c.Word)
	}
	valid := func(c core.ScoredMove) bool {
		return board.ValidateMove(c.PlacedTiles, wordDB)
	}
	fmt.Printf("A has %d vertical\n", len(filter(a, vertical)))
	fmt.Printf("B has %d vertical\n", len(filter(b, vertical)))

	fmt.Printf("A has %d horizontal\n", len(filter(a, horizontal)))
	fmt.Printf("B has %d horizontal\n", len(filter(b, horizontal)))

	fmt.Printf("A has %d withBlank\n", len(filter(a, withBlank)))
	fmt.Printf("B has %d withBlank\n", len(filter(b, withBlank)))

	fmt.Printf("A has %d valid\n", len(filter(a, valid)))
	fmt.Printf("B has %d valid\n", len(filter(b, valid)))
}

func TestSpeedyMatchesSmarty(t *testing.T) {
	tiles := core.NewConsumableRack(core.MakeTiles(core.MakeWord("asdjdha"), "xxxxxx "))
	board := core.NewBoard()

	board.PlaceTiles(core.PlacedTiles{core.MakeTiles(core.MakeWord("doggo"), "xxxxx"), 7, 7, core.Vertical})
	board.PlaceTiles(core.PlacedTiles{core.MakeTiles(core.MakeWord("ar"), "xx"), 7, 8, core.Vertical})

	speedy := ai.NewSpeedyAI(wordDB, wordGaddag)
	defer speedy.Kill()
	smarty := ai.NewSmartyAI(wordDB, wordDB)
	defer smarty.Kill()

	speedyMoves := []core.Turn{}
	speedy.GenerateMoves(board, tiles, func(t core.Turn) bool {
		speedyMoves = append(speedyMoves, t)
		return true
	})

	smartyMoves := []core.Turn{}
	smarty.GenerateMoves(board, tiles, func(t core.Turn) bool {
		smartyMoves = append(smartyMoves, t)
		return true
	})

	if !assert.Equal(t, len(smartyMoves), len(speedyMoves)) {
		compareSets(board, smartyMoves, speedyMoves)
	}
	assert.Subset(t, speedyMoves, smartyMoves)
	assert.Subset(t, smartyMoves, speedyMoves)
}

func BenchmarkSpeedy(b *testing.B) {
	tiles := core.NewConsumableRack(core.MakeTiles(core.MakeWord("bdhrigs"), "xxxxxx "))
	board := core.NewBoard()
	speedy := ai.NewSpeedyAI(wordDB, wordGaddag)
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
	speedy := ai.NewSpeedyAI(wordDB, wordGaddag)
	defer speedy.Kill()

	board.PlaceTiles(core.PlacedTiles{core.MakeTiles(core.MakeWord("aaaaaaaaaaaaaaa"), "xxxxxxxxxxxxxxx"), 0, 7, core.Vertical})
	board.PlaceTiles(core.PlacedTiles{core.MakeTiles(core.MakeWord("aaaaaaaaaaaaaa"), "xxxxxxxxxxxxxxx"), 7, 0, core.Horizontal})
	prev := []core.Tile{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		speedy.Search(board, 8, 8, core.Horizontal, rack, wordGaddag, prev, func([]core.Tile) {})
	}
}
