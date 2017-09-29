package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"

	"github.com/Logiraptor/word-bot/ai"
	"github.com/Logiraptor/word-bot/core"
	"github.com/Logiraptor/word-bot/definitions"
	"github.com/Logiraptor/word-bot/wordlist"
)

var wordDB *wordlist.Trie

func init() {
	builder := wordlist.NewTrieBuilder(151434)
	err := definitions.LoadWords("../words.txt", builder)
	if err != nil {
		panic(err)
	}

	wordDB = builder.Build()
}

type AI interface {
	FindMoves(rack []core.Tile) []ai.ScoredMove
	Kill()
}

func main() {

	resultFile, err := os.Create("results.csv")
	if err != nil {
		panic(err)
	}
	defer resultFile.Close()

	wr := csv.NewWriter(resultFile)

	wr.Write([]string{
		"Smarty 1",
		"Smarty 2",
	})

	for i := 0; i < 1000; i++ {

		p1, p2 := playGame(
			func(b *core.Board) *Player {
				return NewPlayer(ai.NewSmartyAI(b, wordDB, wordDB), "Smarty 1")
			}, func(b *core.Board) *Player {
				return NewPlayer(ai.NewSmartyAI(b, wordDB, wordDB), "Smarty 2")
			},
		)

		fmt.Println("GAME OVER")
		fmt.Printf("Final Score %s = %d\n", p1.name, p1.score)
		fmt.Printf("Final Score %s = %d\n", p2.name, p2.score)

		wr.Write([]string{
			strconv.Itoa(int(p1.score)),
			strconv.Itoa(int(p2.score)),
		})
	}
	wr.Flush()
}

type Player struct {
	ai    AI
	name  string
	rack  core.ConsumableRack
	score core.Score
}

func NewPlayer(ai AI, name string) *Player {
	return &Player{
		ai:   ai,
		name: name,
		rack: core.NewConsumableRack(nil),
	}
}

func (p *Player) takeTurn(board *core.Board, bag core.ConsumableBag) (core.ConsumableBag, bool) {
	moves := p.ai.FindMoves(p.rack.Rack)
	if len(moves) == 0 {
		return bag, false
	}

	move := moves[0]
	if !board.ValidateMove(move.Word, move.Row, move.Col, move.Direction, wordDB) {
		fmt.Printf("%s played an invalid move: %v!\n", p.name, move)
		return bag, false
	}

	newRack, ok := p.rack.Play(move.Word)
	if !ok {
		return bag, false
	}

	p.rack = newRack

	score := board.Score(move.Word, move.Row, move.Col, move.Direction)
	board.PlaceTiles(move.Word, move.Row, move.Col, move.Direction)

	bag, p.rack.Rack = bag.FillRack(p.rack.Rack, 7-len(p.rack.Rack))

	p.score += score

	return bag, true
}

func playGame(a, b func(board *core.Board) *Player) (p1, p2 *Player) {
	swapped := false
	if rand.Intn(2) == 0 {
		swapped = true
		a, b = b, a
	}

	board := core.NewBoard()
	p1 = a(board)
	p2 = b(board)
	bag := core.NewConsumableBag().Shuffle()

	bag, p1.rack.Rack = bag.FillRack(p1.rack.Rack, 7)
	bag, p2.rack.Rack = bag.FillRack(p2.rack.Rack, 7)

	var (
		p1Ok, p2Ok = true, true
	)

	for bag.Count() > 0 && (p1Ok || p2Ok) {
		bag, p1Ok = p1.takeTurn(board, bag)
		bag, p2Ok = p2.takeTurn(board, bag)
	}

	p1.ai.Kill()
	p2.ai.Kill()

	if swapped {
		p1, p2 = p2, p1
	}

	return p1, p2
}
