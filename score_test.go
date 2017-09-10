package main

import "testing"
import "github.com/stretchr/testify/assert"

func assertScore(t *testing.T, b *Board, score int, word string, row, col int, dir Direction) {
	t.Helper()
	result := b.Score(word, row, col, dir)
	if !assert.Equal(t, score, result) {
		b.PlaceTiles(word, row, col, dir)
		b.Print()
	}
}

func TestFirstWord(t *testing.T) {
	b := NewBoard()
	assertScore(t, b, 10, "dog", 7, 7, Horizontal)

	b = NewBoard()
	assertScore(t, b, 14, "goats", 7, 7, Horizontal)
}

func TestSecondWord(t *testing.T) {
	b := NewBoard()
	b.PlaceTiles("dog", 7, 7, Horizontal)
	assertScore(t, b, 8, "goats", 7, 9, Vertical)
}

func TestMultiWord(t *testing.T) {
	b := NewBoard()
	b.PlaceTiles("barn", 7, 7, Horizontal)
	words := b.FindNewWords("bob", 6, 7, Horizontal)
	assert.Len(t, words, 4)
	// 8 for bob
	// 6 for bb
	// 4 for br
	// 3 for oa
	assertScore(t, b, 21, "bob", 6, 7, Horizontal)
}

func TestPlacingJustAnS(t *testing.T) {
	b := NewBoard()
	b.PlaceTiles("dog", 7, 7, Horizontal)
	assertScore(t, b, 6, "s", 7, 10, Horizontal)
}
