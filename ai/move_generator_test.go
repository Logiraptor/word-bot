package ai_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Logiraptor/word-bot/core"
	"github.com/Logiraptor/word-bot/wordlist"

	"github.com/Logiraptor/word-bot/ai"
)

type MoveGenConstructor func(words []string) ai.MoveGenerator

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

func TestBruteMoveGen(t *testing.T) {
	MoveGeneratorContract(t, func(words []string) ai.MoveGenerator {
		wordDB := wordlist.NewTrie()

		for _, word := range words {
			wordDB.AddWord(word)
		}

		return ai.NewBrute(wordDB)
	})
}

func collectMoves(board *core.Board, rack core.Rack, moveGen ai.MoveGenerator) []core.PlacedTiles {
	output := []core.PlacedTiles{}
	moveGen.GenerateMoves(board, rack, func(t core.Turn) bool {
		if m, ok := t.(core.ScoredMove); ok {
			output = append(output, board.NormalizeMove(m.PlacedTiles))
		}
		return true
	})
	return output
}

func MoveGeneratorContract(t *testing.T, makeMoveGenerator MoveGenConstructor) {
	for _, tc := range moveGenTestData {
		t.Run(tc.name, func(t *testing.T) {
			ai := makeMoveGenerator(tc.dictionary)
			board := core.NewBoard()
			for _, m := range tc.previousMoves {
				board.PlaceTiles(m)
			}
			rack := core.NewConsumableRack(tiles(tc.rack))
			moves := collectMoves(board, rack, ai)
			assert.Subset(t, moves, tc.expectedMoves, tc.name)
		})
	}
}
