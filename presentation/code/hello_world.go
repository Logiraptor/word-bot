package main

import (
	"github.com/Logiraptor/word-bot/core"
)

func main() {
	b := core.NewBoard()
	b.PlaceTiles(core.PlacedTiles{
		Col:  7, Row: 7, Direction: core.Horizontal,
		Word: core.MakeTiles(core.MakeWord("hello"), "xxxxx"),
	})
	b.PlaceTiles(core.PlacedTiles{
		Col:  11, Row: 6, Direction: core.Vertical,
		Word: core.MakeTiles(core.MakeWord("wrld"), "xxxx"),
	})
	b.Print()
}
