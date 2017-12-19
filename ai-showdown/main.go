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
	numIterations := 1000
	for i := 0; i < numIterations; i++ {
		jobs <- Job{
			p1: func(b *core.Board) *ai.Player {
				return ai.NewPlayer(smarty)
			},
			p2: func(b *core.Board) *ai.Player {
				return ai.NewPlayer(weighted)
			},
		}
		if i%100 == 0 {
			fmt.Println("Enqueued", i, "/", numIterations, "jobs")
		}
	}

	close(jobs)
	fmt.Println("Done enqueuing, waiting for final jobs to terminate")

	wg.Wait()
}

type Job struct {
	p1, p2 func(b *core.Board) *ai.Player
}

func worker(db *persist.DB, wg *sync.WaitGroup, jobs <-chan Job) {
	for j := range jobs {
		g := ai.PlayGame(wordDB, j.p1, j.p2)
		err := db.SaveGame(g)
		if err != nil {
			fmt.Println("ERROR Saving Game", err)
		}
	}
	wg.Done()
}
