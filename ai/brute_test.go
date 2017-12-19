package ai_test

import (
	"reflect"
	"testing"

	"github.com/Logiraptor/word-bot/ai"
	"github.com/Logiraptor/word-bot/core"
	"github.com/stretchr/testify/assert"
)

func subset(t *testing.T, super, sub [][]core.Tile) bool {
	good := true
outer:
	for i := range sub {
		for j := range super {
			a := super[j]
			b := sub[i]
			if len(a) == 0 && 0 == len(b) {
				continue outer
			}
			if reflect.DeepEqual(a, b) {
				continue outer
			}
		}
		t.Logf("%v missing in set %v", sub[i], super)
		good = false
	}
	return good
}

func permContains(l [][]core.Tile, e []core.Tile) bool {
	for i := range l {
		if reflect.DeepEqual(l[i], e) {
			return true
		}
	}
	return false
}

func TestPermute(t *testing.T) {
	set := core.MakeTiles(core.MakeWord("abc"), "xxx")
	result := ai.Permute(set)
	expected := [][]core.Tile{
		core.MakeTiles(core.MakeWord("abc"), "xxx"),
		core.MakeTiles(core.MakeWord("acb"), "xxx"),
		core.MakeTiles(core.MakeWord("bac"), "xxx"),
		core.MakeTiles(core.MakeWord("bca"), "xxx"),
		core.MakeTiles(core.MakeWord("cab"), "xxx"),
		core.MakeTiles(core.MakeWord("cba"), "xxx"),

		core.MakeTiles(core.MakeWord("ab"), "xx"),
		core.MakeTiles(core.MakeWord("ac"), "xx"),

		core.MakeTiles(core.MakeWord("ba"), "xx"),
		core.MakeTiles(core.MakeWord("bc"), "xx"),

		core.MakeTiles(core.MakeWord("cb"), "xx"),
		core.MakeTiles(core.MakeWord("ca"), "xx"),

		core.MakeTiles(core.MakeWord("a"), "x"),
		core.MakeTiles(core.MakeWord("b"), "x"),
		core.MakeTiles(core.MakeWord("c"), "x"),
	}

	assert.True(t, subset(t, result, expected))
	assert.True(t, subset(t, expected, result))
	assert.Equal(t, len(expected), len(result))
}

func TestPermuteBlank(t *testing.T) {
	set := core.MakeTiles(core.MakeWord("a"), " ")
	result := ai.Permute(set)
	expected := [][]core.Tile{
		core.MakeTiles(core.MakeWord("a"), " "),
		core.MakeTiles(core.MakeWord("b"), " "),
		core.MakeTiles(core.MakeWord("c"), " "),
		core.MakeTiles(core.MakeWord("d"), " "),
		core.MakeTiles(core.MakeWord("e"), " "),
		core.MakeTiles(core.MakeWord("f"), " "),
		core.MakeTiles(core.MakeWord("g"), " "),
		core.MakeTiles(core.MakeWord("h"), " "),
		core.MakeTiles(core.MakeWord("i"), " "),
		core.MakeTiles(core.MakeWord("j"), " "),
		core.MakeTiles(core.MakeWord("k"), " "),
		core.MakeTiles(core.MakeWord("l"), " "),
		core.MakeTiles(core.MakeWord("m"), " "),
		core.MakeTiles(core.MakeWord("n"), " "),
		core.MakeTiles(core.MakeWord("o"), " "),
		core.MakeTiles(core.MakeWord("p"), " "),
		core.MakeTiles(core.MakeWord("q"), " "),
		core.MakeTiles(core.MakeWord("r"), " "),
		core.MakeTiles(core.MakeWord("s"), " "),
		core.MakeTiles(core.MakeWord("t"), " "),
		core.MakeTiles(core.MakeWord("u"), " "),
		core.MakeTiles(core.MakeWord("v"), " "),
		core.MakeTiles(core.MakeWord("w"), " "),
		core.MakeTiles(core.MakeWord("x"), " "),
		core.MakeTiles(core.MakeWord("y"), " "),
		core.MakeTiles(core.MakeWord("z"), " "),
	}

	assert.True(t, subset(t, result, expected))
	assert.True(t, subset(t, expected, result))
	assert.Equal(t, len(expected), len(result))
}

func TestPermuteSpotCheck(t *testing.T) {
	tiles := core.MakeTiles(core.MakeWord("asdjdha"), "xxxxxx ")
	perms := ai.Permute(tiles)

	assert.True(t, permContains(perms, core.MakeTiles(core.MakeWord("adds"), "xxxx")))
	assert.True(t, permContains(perms, core.MakeTiles(core.MakeWord("odds"), " xxx")))
}

func BenchmarkBrute(b *testing.B) {
	tiles := core.NewConsumableRack(core.MakeTiles(core.MakeWord("bdhrigs"), "xxxxxx "))
	board := core.NewBoard()
	brute := ai.NewBrute(wordDB)

	board.PlaceTiles(core.PlacedTiles{core.MakeTiles(core.MakeWord("aaaaaaaaaaaaaaa"), "xxxxxxxxxxxxxxx"), 0, 7, core.Vertical})
	board.PlaceTiles(core.PlacedTiles{core.MakeTiles(core.MakeWord("aaaaaaaaaaaaaa"), "xxxxxxxxxxxxxxx"), 7, 0, core.Horizontal})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		brute.GenerateMoves(board, tiles, func(core.Turn) bool { return true })
	}
}
