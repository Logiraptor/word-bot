package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	assertScore(t, b, 8, toTiles("oats"), 7, 9, Vertical)
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

func TestUsing7Letters(t *testing.T) {
	b := NewBoard()
	assertScore(t, b, 134, toTiles("alfresco"), 7, 7, Horizontal)
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

func TestValidation(t *testing.T) {
	b := NewBoard()
	result := b.ValidateMove(toTiles("dugz"), 7, 7, Horizontal)
	assert.Equal(t, false, result)
}

func TestMultiValidationSuccess(t *testing.T) {
	b := NewBoard()
	b.PlaceTiles(toTiles("handy"), 7, 7, Horizontal)
	result := b.ValidateMove(toTiles("stone"), 6, 5, Horizontal)
	assert.Equal(t, true, result)
}

func TestMultiValidation(t *testing.T) {
	b := NewBoard()
	b.PlaceTiles(toTiles("handy"), 7, 7, Horizontal)
	result := b.ValidateMove(toTiles("stones"), 6, 5, Horizontal)
	assert.Equal(t, false, result)
}

func TestValidationDanglingWords(t *testing.T) {
	b := NewBoard()
	b.PlaceTiles(toTiles("handy"), 7, 7, Horizontal)
	result := b.ValidateMove(toTiles("stones"), 5, 5, Horizontal)
	assert.Equal(t, false, result)
}

func TestValidationOverflowingWords(t *testing.T) {
	b := NewBoard()
	x := b.ValidateMove(toTiles("alfresco"), 7, 7, Horizontal)
	assert.Equal(t, true, x)
	b.PlaceTiles(toTiles("alfresco"), 7, 7, Horizontal)
	result := b.ValidateMove(toTiles("tabarded"), 6, 14, Horizontal)
	if !assert.Equal(t, false, result) {
		b.Print()
	}
}

// Regression tests

func TestJint(t *testing.T) {
	b := NewBoard()
	b.PlaceTiles(toTiles("tusseh"), 14, 0, Horizontal)
	result := b.ValidateMove(toTiles("jin"), 11, 0, Vertical)
	assert.False(t, result)
}

func TestBridgingWords(t *testing.T) {
	b := NewBoard()

	b.PlaceTiles(toTiles("bat"), 7, 7, Horizontal)
	b.PlaceTiles(toTiles("oard"), 7, 7, Vertical)
	b.PlaceTiles(toTiles("ravel"), 7, 9, Vertical)

	result := b.ValidateMove(toTiles("ae"), 10, 7, Horizontal)
	assert.True(t, result)

	words := b.FindNewWords(toTiles("ae"), 10, 7, Horizontal)
	assert.Len(t, words, 1)
	if !assert.Equal(t, "rave", tiles2String(words[0].word)) {
		b.Print()
	}
}

func TestConsumableRack(t *testing.T) {
	c := ConsumableRack{
		rack:     toTiles("abcd"),
		consumed: 0,
	}

	assert.Equal(t, true, c.CanConsume(1))
	next := c.Consume(1)

	assert.Equal(t, true, c.CanConsume(1))
	assert.Equal(t, false, next.CanConsume(1))
}
