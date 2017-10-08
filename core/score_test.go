package core

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func toTiles(word string) []Tile {
	return MakeTiles(MakeWord(word), strings.Repeat("x", len(word)))
}

func wordEqual(aStr string, bWord Word) bool {
	aWord := MakeWord(aStr)
	return reflect.DeepEqual(aWord, bWord)
}

type fakeWordList struct{}

func (fakeWordList) Contains(word Word) bool {
	words := []string{
		"alfresco",
		"rave",
		"oh",
		"na",
		"en",
		"stone",
		"stones",
	}

	for _, x := range words {
		if wordEqual(x, word) {
			return true
		}
	}
	return false
}

func assertScore(t *testing.T, b *Board, score Score, word []Tile, row, col int, dir Direction) {
	t.Helper()
	result := b.Score(PlacedTiles{word, row, col, dir})
	if !assert.Equal(t, score, result) {
		b.PlaceTiles(PlacedTiles{word, row, col, dir})
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
	b.PlaceTiles(PlacedTiles{toTiles("dog"), 7, 7, Horizontal})
	assertScore(t, b, 8, toTiles("oats"), 7, 9, Vertical)
}

func TestMultiWord(t *testing.T) {
	b := NewBoard()
	b.PlaceTiles(PlacedTiles{toTiles("barn"), 7, 7, Horizontal})
	words := b.FindNewWords(PlacedTiles{toTiles("bob"), 6, 7, Horizontal})
	assert.Len(t, words, 4)
	// 8 for bob
	// 6 for bb
	// 4 for br
	// 3 for oa
	assertScore(t, b, 21, toTiles("bob"), 6, 7, Horizontal)
}

func TestPlacingJustAnS(t *testing.T) {
	b := NewBoard()
	b.PlaceTiles(PlacedTiles{toTiles("dog"), 7, 7, Horizontal})
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
	tile := Rune2Letter('o').ToTile(true)
	r := tile.ToRune()
	assert.Equal(t, r, 'o')
}

func TestValidation(t *testing.T) {
	b := NewBoard()
	result := b.ValidateMove(PlacedTiles{toTiles("dugz"), 7, 7, Horizontal}, fakeWordList{})
	assert.Equal(t, false, result)
}

func TestMultiValidationSuccess(t *testing.T) {
	b := NewBoard()
	b.PlaceTiles(PlacedTiles{toTiles("handy"), 7, 7, Horizontal})
	result := b.ValidateMove(PlacedTiles{toTiles("stone"), 6, 5, Horizontal}, fakeWordList{})
	assert.Equal(t, true, result)
}

func TestMultiValidation(t *testing.T) {
	b := NewBoard()
	b.PlaceTiles(PlacedTiles{toTiles("handy"), 7, 7, Horizontal})
	result := b.ValidateMove(PlacedTiles{toTiles("stones"), 6, 5, Horizontal}, fakeWordList{})
	assert.Equal(t, false, result)
}

func TestValidationDanglingWords(t *testing.T) {
	b := NewBoard()
	b.PlaceTiles(PlacedTiles{toTiles("handy"), 7, 7, Horizontal})
	result := b.ValidateMove(PlacedTiles{toTiles("stones"), 5, 5, Horizontal}, fakeWordList{})
	assert.Equal(t, false, result)
}

func TestValidationOverflowingWords(t *testing.T) {
	b := NewBoard()
	x := b.ValidateMove(PlacedTiles{toTiles("alfresco"), 7, 7, Horizontal}, fakeWordList{})
	assert.Equal(t, true, x)
	b.PlaceTiles(PlacedTiles{toTiles("alfresco"), 7, 7, Horizontal})
	result := b.ValidateMove(PlacedTiles{toTiles("tabarded"), 6, 14, Horizontal}, fakeWordList{})
	if !assert.Equal(t, false, result) {
		b.Print()
	}
}

// Regression tests

func TestJint(t *testing.T) {
	b := NewBoard()
	b.PlaceTiles(PlacedTiles{toTiles("tusseh"), 14, 0, Horizontal})
	result := b.ValidateMove(PlacedTiles{toTiles("jin"), 11, 0, Vertical}, fakeWordList{})
	assert.False(t, result)
}

func TestBridgingWords(t *testing.T) {
	b := NewBoard()

	b.PlaceTiles(PlacedTiles{toTiles("bat"), 7, 7, Horizontal})
	b.PlaceTiles(PlacedTiles{toTiles("oard"), 7, 7, Vertical})
	b.PlaceTiles(PlacedTiles{toTiles("ravel"), 7, 9, Vertical})

	result := b.ValidateMove(PlacedTiles{toTiles("ae"), 10, 7, Horizontal}, fakeWordList{})
	assert.True(t, result)

	words := b.FindNewWords(PlacedTiles{toTiles("ae"), 10, 7, Horizontal})
	assert.Len(t, words, 1)
	if !assert.Equal(t, "rave", Tiles2String(words[0].Word)) {
		b.Print()
	}
}

func TestConsumableRack(t *testing.T) {
	c := Rack{
		Rack:     toTiles("abcd"),
		consumed: 0,
	}

	assert.Equal(t, true, c.CanConsume(1))
	next := c.Consume(1)

	assert.Equal(t, true, c.CanConsume(1))
	assert.Equal(t, false, next.CanConsume(1))
}
