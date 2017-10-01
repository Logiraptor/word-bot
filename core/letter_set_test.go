package core

import "testing"
import "github.com/stretchr/testify/assert"

func TestLetterSetEmpty(t *testing.T) {
	s := NewEmptyTileSet()
	a := Rune2Letter('a').ToTile(false)
	z := Rune2Letter('z').ToTile(false)
	blank := Rune2Letter('a').ToTile(true)
	assert.False(t, s.CanConsume(a))
	assert.False(t, s.CanConsume(z))
	assert.False(t, s.CanConsume(blank))
}

func TestLetterSetAddAndConsume(t *testing.T) {
	s := NewEmptyTileSet()
	a := Rune2Letter('a').ToTile(false)
	z := Rune2Letter('z').ToTile(false)
	blank := Rune2Letter('a').ToTile(true)
	s.Add(a)
	assert.True(t, s.CanConsume(a))
	assert.False(t, s.CanConsume(z))
	assert.False(t, s.CanConsume(blank))

	s, consumed := s.Consume(a)
	assert.Equal(t, consumed, a)
	assert.False(t, s.CanConsume(a))
	assert.False(t, s.CanConsume(z))
	assert.False(t, s.CanConsume(blank))
}

func TestLetterSetBlanksCanConsumeAnything(t *testing.T) {
	s := NewEmptyTileSet()
	a := Rune2Letter('a').ToTile(false)
	z := Rune2Letter('z').ToTile(false)
	blank := Rune2Letter('a').ToTile(true)
	s.Add(blank)
	assert.True(t, s.CanConsume(a))
	assert.True(t, s.CanConsume(z))
	assert.True(t, s.CanConsume(blank))
}

func TestLetterSetBlanksAreConsumedLast(t *testing.T) {
	s := NewEmptyTileSet()
	a := Rune2Letter('a').ToTile(false)
	blank := Rune2Letter('a').ToTile(true)
	s.Add(blank)
	s.Add(a)
	s, consumed1 := s.Consume(a)
	s, consumed2 := s.Consume(blank)
	assert.Equal(t, consumed1, a)
	assert.Equal(t, consumed2, blank)
}
