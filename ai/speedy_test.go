package ai_test

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Logiraptor/word-bot/ai"
	"github.com/Logiraptor/word-bot/core"
	"github.com/Logiraptor/word-bot/definitions"
	"github.com/Logiraptor/word-bot/wordlist"
)

var wordGaddag *wordlist.Gaddag

func init() {
	wordGaddag = wordlist.NewGaddag()
	err := definitions.LoadWords("../words.txt", wordGaddag)
	if err != nil {
		panic(err)
	}
}

func filter(scoredMoves []core.Turn, pred func(core.ScoredMove) bool) []core.ScoredMove {
	output := []core.ScoredMove{}
	for i := range scoredMoves {
		if m, ok := scoredMoves[i].(core.ScoredMove); ok {
			if pred(m) {
				output = append(output, m)
			}
		}
	}
	return output
}

func unique(turns []core.Turn) []core.Turn {
	output := []core.Turn{}
outer:
	for _, m := range turns {
		if contains(output, m) {
			continue outer
		}
		output = append(output, m)
	}
	return output
}

func contains(turns []core.Turn, t core.Turn) bool {
	for _, x := range turns {
		if reflect.DeepEqual(x, t) {
			return true
		}
	}
	return false
}

func intersection(a, b []core.Turn) []core.Turn {
	output := []core.Turn{}
	for i := range a {
		if contains(a, a[i]) && contains(b, a[i]) && (!contains(output, a[i])) {
			output = append(output, a[i])
		}
	}
	return output
}

func compareSets(board *core.Board, aName, bName string, a, b []core.Turn) {
	hasBlank := func(t []core.Tile) bool {
		for _, x := range t {
			if x.IsBlank() {
				return true
			}
		}
		return false
	}
	vertical := func(c core.ScoredMove) bool {
		return c.Direction == core.Vertical
	}
	horizontal := func(c core.ScoredMove) bool {
		return c.Direction == core.Horizontal
	}
	withBlank := func(c core.ScoredMove) bool {
		return hasBlank(c.Word)
	}
	valid := func(c core.ScoredMove) bool {
		return board.ValidateMove(c.PlacedTiles, wordDB)
	}
	fmt.Printf("%s has %d vertical\n", aName, len(filter(a, vertical)))
	fmt.Printf("%s has %d vertical\n", bName, len(filter(b, vertical)))

	fmt.Printf("%s has %d horizontal\n", aName, len(filter(a, horizontal)))
	fmt.Printf("%s has %d horizontal\n", bName, len(filter(b, horizontal)))

	fmt.Printf("%s has %d withBlank\n", aName, len(filter(a, withBlank)))
	fmt.Printf("%s has %d withBlank\n", bName, len(filter(b, withBlank)))

	fmt.Printf("%s has %d valid\n", aName, len(filter(a, valid)))
	fmt.Printf("%s has %d valid\n", bName, len(filter(b, valid)))

	fmt.Printf("%s has %d unique\n", aName, len(unique(a)))
	fmt.Printf("%s has %d unique\n", bName, len(unique(b)))

	fmt.Printf("Intersection is %d\n", len(intersection(a, b)))

	fmt.Printf("%s unique has %d vertical\n", aName, len(filter(unique(a), vertical)))
	fmt.Printf("%s unique has %d vertical\n", bName, len(filter(unique(b), vertical)))

	fmt.Printf("%s unique has %d horizontal\n", aName, len(filter(unique(a), horizontal)))
	fmt.Printf("%s unique has %d horizontal\n", bName, len(filter(unique(b), horizontal)))

	fmt.Printf("%s unique has %d withBlank\n", aName, len(filter(unique(a), withBlank)))
	fmt.Printf("%s unique has %d withBlank\n", bName, len(filter(unique(b), withBlank)))

	fmt.Printf("%s unique has %d valid\n", aName, len(filter(unique(a), valid)))
	fmt.Printf("%s unique has %d valid\n", bName, len(filter(unique(b), valid)))
}

func dumpTurns(filename string, moves []core.Turn) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	wr := csv.NewWriter(f)
	defer wr.Flush()
	for _, t := range moves {
		if m, ok := t.(core.ScoredMove); ok {
			wr.Write([]string{
				strconv.Itoa(m.Row),
				strconv.Itoa(m.Col),
				core.Tiles2String(m.Word),
				m.Direction.String(),
			})
		}
	}
}

func TestSpeedyMatchesSmarty(t *testing.T) {
	tiles := core.NewConsumableRack(core.MakeTiles(core.MakeWord("asdjdha"), "xxxxxx "))
	board := core.NewBoard()

	board.PlaceTiles(core.PlacedTiles{core.MakeTiles(core.MakeWord("doggo"), "xxxxx"), 7, 7, core.Horizontal})
	board.PlaceTiles(core.PlacedTiles{core.MakeTiles(core.MakeWord("ar"), "xx"), 7, 8, core.Vertical})

	speedy := ai.NewSpeedyAI(wordDB, wordGaddag)
	defer speedy.Kill()
	smarty := ai.NewSmartyAI(wordDB, wordDB)
	defer smarty.Kill()

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

	bruteMoves := []core.Turn{}
	ai.BruteForce(board, tiles, wordDB, func(t core.Turn) {
		bruteMoves = append(bruteMoves, t)
	})

	speedyMoves = unique(speedyMoves)
	smartyMoves = unique(smartyMoves)
	bruteMoves = unique(bruteMoves)

	dumpTurns("brute.csv", bruteMoves)
	dumpTurns("smarty.csv", smartyMoves)
	dumpTurns("speedy.csv", speedyMoves)
	// pt := core.PlacedTiles{
	// 	Word:      core.MakeTiles(core.MakeWord("add"), " xx"),
	// 	Row:       6,
	// 	Col:       10,
	// 	Direction: core.Horizontal,
	// }
	// fmt.Println("HERERERERERERERERERERER", board.ValidateMove(pt, wordDB))
	// smarty.Search(board, 6, 10, core.Horizontal, tiles, wordDB, nil, func(t []core.Tile) {
	// 	if reflect.DeepEqual(t, pt.Word) {
	// 		fmt.Println("GENERATED")
	// 	}
	// })

	assert.Subset(t, bruteMoves, smartyMoves)
	if !assert.Equal(t, len(bruteMoves), len(smartyMoves)) {
		compareSets(board, "brute", "smarty", bruteMoves, smartyMoves)
		return
	}
	// assert.Subset(t, speedyMoves, smartyMoves)
	assert.Subset(t, smartyMoves, speedyMoves)
	if !assert.Equal(t, len(smartyMoves), len(speedyMoves)) {
		compareSets(board, "smarty", "speedy", smartyMoves, speedyMoves)
	}
	board.Print()
}

func BenchmarkSpeedy(b *testing.B) {
	tiles := core.NewConsumableRack(core.MakeTiles(core.MakeWord("bdhrigs"), "xxxxxx "))
	board := core.NewBoard()
	speedy := ai.NewSpeedyAI(wordDB, wordGaddag)
	bag := core.NewConsumableBag()
	defer speedy.Kill()

	board.PlaceTiles(core.PlacedTiles{core.MakeTiles(core.MakeWord("aaaaaaaaaaaaaaa"), "xxxxxxxxxxxxxxx"), 0, 7, core.Vertical})
	board.PlaceTiles(core.PlacedTiles{core.MakeTiles(core.MakeWord("aaaaaaaaaaaaaa"), "xxxxxxxxxxxxxxx"), 7, 0, core.Horizontal})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		speedy.FindMove(board, bag, tiles, func(core.Turn) bool { return true })
	}
}

func BenchmarkSpeedySearch(b *testing.B) {
	rack := core.NewConsumableRack(core.MakeTiles(core.MakeWord("bdhrigs"), "xxxxxx "))
	board := core.NewBoard()
	speedy := ai.NewSpeedyAI(wordDB, wordGaddag)
	defer speedy.Kill()

	board.PlaceTiles(core.PlacedTiles{core.MakeTiles(core.MakeWord("aaaaaaaaaaaaaaa"), "xxxxxxxxxxxxxxx"), 0, 7, core.Vertical})
	board.PlaceTiles(core.PlacedTiles{core.MakeTiles(core.MakeWord("aaaaaaaaaaaaaa"), "xxxxxxxxxxxxxxx"), 7, 0, core.Horizontal})
	prev := []core.Tile{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		speedy.Search(board, 8, 8, core.Horizontal, rack, wordGaddag, prev, func(int, int, []core.Tile) {})
	}
}
