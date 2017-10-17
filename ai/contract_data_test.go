package ai_test

import (
	"unicode"

	"github.com/Logiraptor/word-bot/core"
)

func move(i, j int, dir core.Direction, word string) core.PlacedTiles {
	return core.PlacedTiles{
		Row:       i,
		Col:       j,
		Direction: dir,
		Word:      tiles(word),
	}
}

func tiles(s string) []core.Tile {
	output := []core.Tile{}
	for _, l := range s {
		output = append(output, core.Rune2Letter(unicode.ToLower(l)).ToTile(unicode.IsUpper(l)))
	}
	return output
}

type MoveGenTestCase struct {
	dictionary    []string
	previousMoves []core.PlacedTiles
	rack          string
	expectedMoves []core.PlacedTiles
}

var moveGenTestData = []MoveGenTestCase{
	{
		dictionary: []string{"cab"},
		previousMoves: []core.PlacedTiles{
			move(7, 6, core.Horizontal, "cab"),
		},
		rack: "bc",
		expectedMoves: []core.PlacedTiles{
			move(6, 7, core.Vertical, "cb"),
		},
	},
}
