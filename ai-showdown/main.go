package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os/exec"

	"runtime"
	"sync"
	"time"

	"github.com/Logiraptor/word-bot/ai"
	"github.com/Logiraptor/word-bot/core"
	"github.com/Logiraptor/word-bot/persist"
	"github.com/Logiraptor/word-bot/wordlist"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var wordDB *wordlist.Trie
var commitHash []byte

func init() {
	wordDB = wordlist.MakeDefaultWordList()

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

	db, err := persist.NewDB("smart-results.db")
	if err != nil {
		panic(err)
	}

	numWorkers := runtime.NumCPU() / 2
	jobs := make(chan Job, numWorkers*2)

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(db, &wg, jobs)
	}

	smarty := ai.NewSmartyAI(wordDB, wordDB)
	weighted := ai.NewMoveChooser("Weighted - From Data"+time.Now().Format("02-15:04"), smarty, ai.NewLeaveWeighter(db))
	numIters := 1000
	for i := 0; i < numIters; i++ {
		jobs <- Job{
			p1: func(b *core.Board) *Player {
				return NewPlayer(smarty, 1)
			},
			p2: func(b *core.Board) *Player {
				return NewPlayer(weighted, 2)
			},
		}
		if i%100 == 0 {
			fmt.Println("Enqueued", i, "/", numIters, "jobs")
		}
	}

	close(jobs)
	fmt.Println("Done enqueuing, waiting for final jobs to terminate")

	wg.Wait()
}

type Job struct {
	p1, p2 func(b *core.Board) *Player
}

func worker(db *persist.DB, wg *sync.WaitGroup, jobs <-chan Job) {
	for j := range jobs {
		g := playGame(j.p1, j.p2)
		err := db.SaveGame(g)
		if err != nil {
			fmt.Println("ERROR Saving Game", err)
		}
	}
	wg.Done()
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

func (p *Player) takeTurn(board *core.Board, bag core.Bag) (core.Bag, []core.Tile, core.ScoredMove, bool) {
	var turn core.Turn
	p.ai.FindMove(board, bag, p.rack, func(t core.Turn) bool {
		turn = t
		return true
	})
	if turn == nil {
		return bag, nil, core.ScoredMove{}, false
	}

	switch move := turn.(type) {
	case core.ScoredMove:
		if !board.ValidateMove(move.PlacedTiles, wordDB) {
			fmt.Printf("%s played an invalid move: %v!\n", p.name, move)
			return bag, nil, core.ScoredMove{}, false
		}

		newRack, ok := p.rack.Play(move.Word)
		if !ok {
			return bag, nil, core.ScoredMove{}, false
		}

		p.rack = newRack

		score := board.Score(move.PlacedTiles)
		board.PlaceTiles(move.PlacedTiles)

		leave := newRack.Rack
		bag, p.rack.Rack = bag.FillRack(p.rack.Rack, 7-len(p.rack.Rack))

		p.score += score

		return bag, leave, move, true
	case core.Pass:
		return bag, nil, core.ScoredMove{}, false
	case core.Exchange:
		newRack := core.NewConsumableRack(nil)
		bag, newRack.Rack = bag.FillRack(newRack.Rack, 7-len(newRack.Rack))
		bag = bag.Replace(p.rack.Rack)
		p.rack = newRack
		bag, p.rack.Rack = bag.FillRack(p.rack.Rack, 7-len(p.rack.Rack))
		return bag, nil, core.ScoredMove{}, false
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
		leave      []core.Tile
	)

	for bag.Count() > 0 && (p1Ok || p2Ok) {
		if bag, leave, move, p1Ok = p1.takeTurn(board, bag); p1Ok {
			game.AddMove(p1.name, leave, move)
		}
		if bag, leave, move, p2Ok = p2.takeTurn(board, bag); p2Ok {
			game.AddMove(p2.name, leave, move)
		}

		//fmt.Println(bag.Count(), "Tiles left:", p1.name, p1.score, "to", p2.name, p2.score)
	}

	if swapped {
		p1, p2 = p2, p1
	}

	//fmt.Println("GAME OVER")
	//fmt.Printf("Final Score %s = %d\n", p1.name, p1.score)
	//fmt.Printf("Final Score %s = %d\n", p2.name, p2.score)

	//board.Print()

	return game
}
