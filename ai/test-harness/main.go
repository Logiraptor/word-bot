package main

import (
	"fmt"
	"reflect"

	"github.com/Logiraptor/word-bot/ai"
	"github.com/Logiraptor/word-bot/core"
	"github.com/Logiraptor/word-bot/definitions"
	"github.com/Logiraptor/word-bot/wordlist"
)

var wordDB *wordlist.Trie

func init() {
	wordDB = wordlist.NewTrie()
	err := definitions.LoadWords("../../words.txt", wordDB)
	if err != nil {
		panic(err)
	}
}

func diff(a, b []core.Turn) []core.Turn {
	output := []core.Turn{}
outer:
	for i := range a {
		for j := range b {
			if reflect.DeepEqual(a[i], b[j]) {
				continue outer
			}
		}
		output = append(output, a[i])
	}
	return output
}

func contains(l [][]core.Tile, e []core.Tile) bool {
	for i := range l {
		if reflect.DeepEqual(l[i], e) {
			return true
		}
	}
	return false
}

func main() {
	tiles := core.NewConsumableRack(core.MakeTiles(core.MakeWord("asdjdha"), "xxxxxx "))
	board := core.NewBoard()

	board.PlaceTiles(core.PlacedTiles{core.MakeTiles(core.MakeWord("doggo"), "xxxxx"), 7, 7, core.Horizontal})
	board.PlaceTiles(core.PlacedTiles{core.MakeTiles(core.MakeWord("ar"), "xx"), 7, 8, core.Vertical})

	smarty := ai.NewSmartyAI(wordDB, wordDB)

	smartyMoves := []core.Turn{}
	smarty.GenerateMoves(board, tiles, func(t core.Turn) bool {
		smartyMoves = append(smartyMoves, t)
		return true
	})

	bruteMoves := []core.Turn{}
	ai.BruteForce(board, tiles, wordDB, func(t core.Turn) {
		bruteMoves = append(bruteMoves, t)
	})

	d := diff(bruteMoves, smartyMoves)
	for _, x := range d {
		m, ok := x.(core.ScoredMove)
		if !ok {
			continue
		}
		b := board.Clone()
		fmt.Println("BEFORE")
		b.Print()
		b.PlaceTiles(m.PlacedTiles)
		fmt.Println("AFTER", m)
		b.Print()
		smarty.Search(board, m.Row, m.Col, m.Direction, core.NewConsumableRack(m.Word), wordDB, nil, func([]core.Tile) {})
		fmt.Scanln()
	}
}
