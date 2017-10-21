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

var wordGaddag *wordlist.Gaddag

func init() {
	wordDB = wordlist.NewTrie()
	wordGaddag = wordlist.NewGaddag()
	err := definitions.LoadWords("../../words.txt", wordDB)
	if err != nil {
		panic(err)
	}
	err = definitions.LoadWords("../../words.txt", wordGaddag)
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
	defer smarty.Kill()
	speedy := ai.NewSpeedyAI(wordDB, wordGaddag)
	defer speedy.Kill()

	speedyMoves := []core.Turn{}
	speedy.GenerateMoves(board, tiles, func(t core.Turn) bool {
		speedyMoves = append(speedyMoves, t)
		return true
	})

	smartyMoves := []core.Turn{}
	smarty.GenerateMoves(board, tiles, func(t core.Turn) bool {
		smartyMoves = append(smartyMoves, t)
		return true
	})

	d := diff(smartyMoves, speedyMoves)
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
		speedy.Search(board, 6, 7, core.Vertical, core.NewConsumableRack(m.Word), wordGaddag, nil, func(int, int, []core.Tile, []core.Tile) {})

		fmt.Scanln()
	}
}

// package main

// import (
// 	"os"

// 	"github.com/Logiraptor/word-bot/wordlist"
// )

// func main() {

// 	f, err := os.Create("out.dot")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer f.Close()

// 	g := wordlist.NewGaddag()
// 	g.AddWord("shaded")
// 	g.AddWord("doggo")
// 	g.AddWord("oar")
// 	g.DumpToDot(f)

// }
