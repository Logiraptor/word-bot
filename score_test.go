package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func toTiles(word string) []Tile {
	return MakeTiles(MakeWord(word), strings.Repeat("x", len(word)))
}

func assertScore(t *testing.T, b *Board, score Score, word []Tile, row, col int, dir Direction) {
	t.Helper()
	result := b.Score(word, row, col, dir)
	if !assert.Equal(t, score, result) {
		b.PlaceTiles(word, row, col, dir)
		b.Print()
	}
}

func TestFirstWord(t *testing.T) {
	b := NewBoard()
	assertScore(t, b, 10, toTiles("dog"), 7, 7, Horizontal)

	b = NewBoard()
	assertScore(t, b, 14, toTiles("goats"), 7, 7, Horizontal)
}

func TestSecondWord(t *testing.T) {
	b := NewBoard()
	b.PlaceTiles(toTiles("dog"), 7, 7, Horizontal)
	assertScore(t, b, 8, toTiles("goats"), 7, 9, Vertical)
}

func TestMultiWord(t *testing.T) {
	b := NewBoard()
	b.PlaceTiles(toTiles("barn"), 7, 7, Horizontal)
	words := b.FindNewWords(toTiles("bob"), 6, 7, Horizontal)
	assert.Len(t, words, 4)
	// 8 for bob
	// 6 for bb
	// 4 for br
	// 3 for oa
	assertScore(t, b, 21, toTiles("bob"), 6, 7, Horizontal)
}

func TestPlacingJustAnS(t *testing.T) {
	b := NewBoard()
	b.PlaceTiles(toTiles("dog"), 7, 7, Horizontal)
	assertScore(t, b, 6, toTiles("s"), 7, 10, Horizontal)
}

func TestBlankTiles(t *testing.T) {
	b := NewBoard()
	assertScore(t, b, 8, MakeTiles(MakeWord("dog"), "x x"), 7, 7, Horizontal)
}

func TestConversions(t *testing.T) {
	tile := rune2Tile('o', true)
	r := letter2Rune(tile.ToLetter())
	assert.Equal(t, r, 'o')
}
