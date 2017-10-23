package ai_test

import (
	"testing"

	"github.com/Logiraptor/word-bot/ai"
	"github.com/Logiraptor/word-bot/core"
	"github.com/Logiraptor/word-bot/wordlist"
	"github.com/stretchr/testify/assert"
)

func TestBoardConstraints(t *testing.T) {
	board := core.NewBoard()
	board.PlaceTiles(move(7, 7, core.Horizontal, "abc"))
	wordDB := wordlist.NewGaddag()
	wordDB.AddWord("abc")
	wordDB.AddWord("ba")
	wordDB.AddWord("bc")
	wordDB.AddWord("ab")
	wordDB.AddWord("ac")
	a := core.Rune2Letter('a').ToTile(false)
	b := core.Rune2Letter('b').ToTile(false)
	c := core.Rune2Letter('c').ToTile(false)

	assert.False(t, ai.PermittedTiles(board, wordDB, 7, 6).AllowsTile(a))
	assert.False(t, ai.PermittedTiles(board, wordDB, 7, 6).AllowsTile(b))
	assert.False(t, ai.PermittedTiles(board, wordDB, 7, 6).AllowsTile(c))

	assert.False(t, ai.PermittedTiles(board, wordDB, 6, 7).AllowsTile(a))
	assert.True(t, ai.PermittedTiles(board, wordDB, 6, 7).AllowsTile(b))
	assert.False(t, ai.PermittedTiles(board, wordDB, 6, 7).AllowsTile(c))

	assert.True(t, ai.PermittedTiles(board, wordDB, 6, 8).AllowsTile(a))
	assert.False(t, ai.PermittedTiles(board, wordDB, 6, 8).AllowsTile(b))
	assert.False(t, ai.PermittedTiles(board, wordDB, 6, 8).AllowsTile(c))

	assert.True(t, ai.PermittedTiles(board, wordDB, 6, 9).AllowsTile(a))
	assert.True(t, ai.PermittedTiles(board, wordDB, 6, 9).AllowsTile(b))
	assert.False(t, ai.PermittedTiles(board, wordDB, 6, 9).AllowsTile(c))

	assert.False(t, ai.PermittedTiles(board, wordDB, 7, 10).AllowsTile(a))
	assert.False(t, ai.PermittedTiles(board, wordDB, 7, 10).AllowsTile(b))
	assert.False(t, ai.PermittedTiles(board, wordDB, 7, 10).AllowsTile(c))

	assert.False(t, ai.PermittedTiles(board, wordDB, 8, 7).AllowsTile(a))
	assert.True(t, ai.PermittedTiles(board, wordDB, 8, 7).AllowsTile(b))
	assert.True(t, ai.PermittedTiles(board, wordDB, 8, 7).AllowsTile(c))

	assert.True(t, ai.PermittedTiles(board, wordDB, 8, 8).AllowsTile(a))
	assert.False(t, ai.PermittedTiles(board, wordDB, 8, 8).AllowsTile(b))
	assert.True(t, ai.PermittedTiles(board, wordDB, 8, 8).AllowsTile(c))

	assert.False(t, ai.PermittedTiles(board, wordDB, 8, 9).AllowsTile(a))
	assert.False(t, ai.PermittedTiles(board, wordDB, 8, 9).AllowsTile(b))
	assert.False(t, ai.PermittedTiles(board, wordDB, 8, 9).AllowsTile(c))
}
