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
	name          string
	dictionary    []string
	previousMoves []core.PlacedTiles
	rack          string
	expectedMoves []core.PlacedTiles
}

var moveGenTestData = []MoveGenTestCase{
	{
		name:          "first move with one word",
		dictionary:    []string{"cab"},
		previousMoves: []core.PlacedTiles{},
		rack:          "cab",
		expectedMoves: []core.PlacedTiles{
			move(6, 7, core.Vertical, "cab"),
		},
	},
	{
		name:       "vertical crossover",
		dictionary: []string{"cab"},
		previousMoves: []core.PlacedTiles{
			move(7, 6, core.Horizontal, "cab"),
		},
		rack: "bc",
		expectedMoves: []core.PlacedTiles{
			move(6, 7, core.Vertical, "cb"),
		},
	},
	{
		name:       "horizontal crossover",
		dictionary: []string{"cab"},
		previousMoves: []core.PlacedTiles{
			move(6, 7, core.Vertical, "cab"),
		},
		rack: "bc",
		expectedMoves: []core.PlacedTiles{
			move(7, 6, core.Horizontal, "cb"),
		},
	},
	{
		name:       "line up",
		dictionary: []string{"cab", "dc", "ea", "fb", "def"},
		previousMoves: []core.PlacedTiles{
			move(7, 7, core.Horizontal, "cab"),
		},
		rack: "def",
		expectedMoves: []core.PlacedTiles{
			move(6, 7, core.Horizontal, "def"),
		},
	},
	{
		name:       "corner play",
		dictionary: []string{"abc", "fb", "fd", "edc"},
		previousMoves: []core.PlacedTiles{
			move(7, 7, core.Horizontal, "abc"),
			move(5, 9, core.Vertical, "ed"),
		},
		rack: "f",
		expectedMoves: []core.PlacedTiles{
			move(6, 8, core.Horizontal, "f"),
		},
	},
}
