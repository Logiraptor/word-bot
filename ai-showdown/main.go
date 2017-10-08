package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"time"

	"github.com/Logiraptor/word-bot/smarter"

	"github.com/Logiraptor/word-bot/ai"
	"github.com/Logiraptor/word-bot/core"
	"github.com/Logiraptor/word-bot/definitions"
	"github.com/Logiraptor/word-bot/persist"
	"github.com/Logiraptor/word-bot/wordlist"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var wordDB *wordlist.Trie
var commitHash []byte

func init() {
	builder := wordlist.NewTrieBuilder(151434)
	err := definitions.LoadWords("../words.txt", builder)
	if err != nil {
		panic(err)
	}

	wordDB = builder.Build()

	status, err := exec.Command("git", "status").Output()
	if err != nil {
		panic(err)
	}

	if bytes.Contains(status, []byte("Changes not staged for commit")) ||
		bytes.Contains(status, []byte("Changes to be committed")) {
		log.Print("There are uncommitted changes, please commit or discard to keep the logs accurate!")
		commitHash = []byte("???")
	} else {
		commitHash, err = exec.Command("git", "rev-parse", "--short", "HEAD").Output()
		if err != nil {
			panic(err)
		}
	}
}

func main() {

	rand.Seed(time.Now().Unix())

	db, err := persist.NewDB("results.db")
	if err != nil {
		panic(err)
	}

	smarty := ai.NewSmartyAI(wordDB, wordDB)
	mcts := smarter.NewMCTSAI(smarty)

	for i := 0; i < 10000; i++ {

		g := playGame(
			func(b *core.Board) *Player {
				return NewPlayer(smarty, 1)
			}, func(b *core.Board) *Player {
				return NewPlayer(mcts, 2)
			},
		)

		err := db.SaveGame(g)
		if err != nil {
			fmt.Println("ERROR SAVING GAME", err)
		}
	}
}

type Player struct {
	ai    ai.AI
	name  string
	rack  core.Rack
	score core.Score
}

func NewPlayer(ai ai.AI, n int) *Player {
	return &Player{
		ai:   ai,
		name: fmt.Sprintf("%s - %d - %s", ai.Name(), n, commitHash),
		rack: core.NewConsumableRack(nil),
	}
}

func (p *Player) takeTurn(board *core.Board, bag core.Bag) (core.Bag, core.ScoredMove, bool) {
	var turn core.Turn
	p.ai.FindMove(board, bag, p.rack, func(t core.Turn) bool {
		turn = t
		return true
	})
	if turn == nil {
		return bag, core.ScoredMove{}, false
	}

	switch move := turn.(type) {
	case core.ScoredMove:
		if !board.ValidateMove(move.PlacedTiles, wordDB) {
			fmt.Printf("%s played an invalid move: %v!\n", p.name, move)
			return bag, core.ScoredMove{}, false
		}

		newRack, ok := p.rack.Play(move.Word)
		if !ok {
			return bag, core.ScoredMove{}, false
		}

		p.rack = newRack

		score := board.Score(move.PlacedTiles)
		board.PlaceTiles(move.PlacedTiles)

		bag, p.rack.Rack = bag.FillRack(p.rack.Rack, 7-len(p.rack.Rack))

		p.score += score

		return bag, move, true
	case core.Pass:
		return bag, core.ScoredMove{}, false
	default:
		panic(fmt.Sprintf("%s played unknown turn type: %#v", p.ai.Name(), move))
	}
}

func playGame(a, b func(board *core.Board) *Player) persist.Game {
	game := persist.Game{}

	swapped := false
	if rand.Intn(2) == 0 {
		swapped = true
		a, b = b, a
	}

	board := core.NewBoard()
	p1 := a(board)
	p2 := b(board)

	bag := core.NewConsumableBag().Shuffle()

	bag, p1.rack.Rack = bag.FillRack(p1.rack.Rack, 7)
	bag, p2.rack.Rack = bag.FillRack(p2.rack.Rack, 7)

	var (
		p1Ok, p2Ok = true, true
		move       core.ScoredMove
	)

	for bag.Count() > 0 && (p1Ok || p2Ok) {
		if bag, move, p1Ok = p1.takeTurn(board, bag); p1Ok {
			game.AddMove(p1.name, move)
		}
		if bag, move, p2Ok = p2.takeTurn(board, bag); p2Ok {
			game.AddMove(p2.name, move)
		}
	}

	if swapped {
		p1, p2 = p2, p1
	}

	fmt.Println("GAME OVER")
	fmt.Printf("Final Score %s = %d\n", p1.name, p1.score)
	fmt.Printf("Final Score %s = %d\n", p2.name, p2.score)

	board.Print()

	return game
}
