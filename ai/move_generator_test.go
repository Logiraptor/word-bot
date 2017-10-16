package ai_test

import (
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"

	"github.com/Logiraptor/word-bot/core"
	"github.com/Logiraptor/word-bot/wordlist"

	"github.com/Logiraptor/word-bot/ai"
)

type MoveGenConstructor func(words []string) ai.MoveGenerator

func tiles(s string) []core.Tile {
	output := []core.Tile{}
	for _, l := range s {
		output = append(output, core.Rune2Letter(unicode.ToLower(l)).ToTile(unicode.IsUpper(l)))
	}
	return output
}

func move(i, j int, dir core.Direction, word string, score core.Score) core.ScoredMove {
	return core.ScoredMove{
		PlacedTiles: core.PlacedTiles{
			Row:       i,
			Col:       j,
			Direction: dir,
			Word:      tiles(word),
		},
		Score: score,
	}
}

func TestSpeedyMoveGen(t *testing.T) {
	MoveGeneratorContract(t, func(words []string) ai.MoveGenerator {
		wordDB := wordlist.NewTrie()
		wordGaddag := wordlist.NewGaddag()

		for _, word := range words {
			wordDB.AddWord(word)
			wordGaddag.AddWord(word)
		}

		return ai.NewSpeedyAI(wordDB, wordGaddag)
	})
}

func TestSmartyMoveGen(t *testing.T) {
	MoveGeneratorContract(t, func(words []string) ai.MoveGenerator {
		wordDB := wordlist.NewTrie()

		for _, word := range words {
			wordDB.AddWord(word)
		}

		return ai.NewSmartyAI(wordDB, wordDB)
	})
}

func collectMoves(board *core.Board, rack core.Rack, moveGen ai.MoveGenerator) []core.ScoredMove {
	output := []core.ScoredMove{}
	moveGen.GenerateMoves(board, rack, func(t core.Turn) bool {
		if m, ok := t.(core.ScoredMove); ok {
			output = append(output, m)
		}
		return true
	})
	return output
}

func MoveGeneratorContract(t *testing.T, makeMoveGenerator MoveGenConstructor) {
	words := []string{"cab"}
	expectedMoves := []core.ScoredMove{
		move(7, 7, core.Horizontal, "cab", 6),
	}
	rackTiles := core.MakeTiles(core.MakeWord("abc"), "xxx")

	ai := makeMoveGenerator(words)
	board := core.NewBoard()
	rack := core.NewConsumableRack(rackTiles)
	moves := collectMoves(board, rack, ai)
	assert.Subset(t, moves, expectedMoves)
}
