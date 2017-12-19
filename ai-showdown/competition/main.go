package main

import (
	"errors"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/Logiraptor/word-bot/ai"
	"github.com/Logiraptor/word-bot/core"
	"github.com/Logiraptor/word-bot/persist"
	"github.com/Logiraptor/word-bot/wordlist"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var wordDB = wordlist.MakeDefaultWordList()
var wordGaddag = wordlist.MakeDefaultWordListGaddag()

func main() {
	rand.Seed(time.Now().Unix())
	connectionString, err := getConnectionString()
	if err != nil {
		log.Printf("Failed to get connection string: %s", err.Error())
		return
	}
	gormDB, err := gorm.Open("postgres", connectionString)
	if err != nil {
		log.Printf("Failed to connect to db: %s", err.Error())
		return
	}
	db, err := persist.NewDBConn(gormDB)
	if err != nil {
		log.Printf("Failed to automigrate db: %s", err.Error())
		return
	}

	var (
		repo     CompetitorRepo = DBCompetitorRepo{DB: db}
		database Database       = db
	)

	for {
		time.Sleep(time.Second * 1)
		err := runIteration(repo, wordDB, database)
		if err != nil {
			log.Printf("Error running iteration: %s", err.Error())
		}
	}
}

type CompetitorPair struct {
	Competitor1, Competitor2 ai.AI
	NumPlays                 int
}

type CompetitorRepo interface {
	CompetitorPairs() ([]CompetitorPair, error)
}

type DBCompetitorRepo struct {
	DB *persist.DB
}

func (d DBCompetitorRepo) CompetitorPairs() ([]CompetitorPair, error) {
	matchups := []persist.Matchup{
		{Player1: "Smarty", Player2: "Smarty", NumGames: 0},
		{Player1: "Smarty", Player2: "Speedy", NumGames: 0},
		{Player1: "Speedy", Player2: "Speedy", NumGames: 0},
	}
	err := d.DB.GetMatchups(matchups)
	if err != nil {
		return nil, err
	}

	log.Println(matchups)

	var output []CompetitorPair
	for _, matchup := range matchups {
		output = append(output, CompetitorPair{
			Competitor1: nameToAi(matchup.Player1),
			Competitor2: nameToAi(matchup.Player2),
			NumPlays:    matchup.NumGames,
		})
	}
	return output, nil
}

func nameToAi(name string) ai.AI {
	switch  {
	case strings.HasPrefix(name, "Smarty"):
		return ai.NewSmartyAI(wordDB, wordDB)
	case strings.HasPrefix(name, "Speedy"):
		return ai.NewSpeedyAI(wordDB, wordGaddag)
	}
	log.Printf("Failed to decode ai name %q, defaulting to smarty", name)
	return ai.NewSmartyAI(wordDB, wordDB)
}

func getConnectionString() (string, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return "", errors.New("missing DATABASE_URL env var")
	}
	return dbURL, nil
}

type Database interface {
	SaveGame(g persist.Game) error
}

type ByNumPlays []CompetitorPair

func (b ByNumPlays) Len() int {
	return len(b)
}

func (b ByNumPlays) Less(i, j int) bool {
	return b[i].NumPlays < b[j].NumPlays
}

func (b ByNumPlays) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func runIteration(repo CompetitorRepo, wordList core.WordList, db Database) error {
	pairs, err := repo.CompetitorPairs()
	if err != nil {
		log.Printf("Error fetching competitor record: %s", err.Error())
	}
	if len(pairs) == 0 {
		log.Println("No competition pairs available")
		return nil
	}
	log.Printf("%d pairs available", len(pairs))
	sort.Sort(ByNumPlays(pairs))
	nextPair := pairs[0]
	log.Printf("Playing Game between %q and %q", nextPair.Competitor1.Name(), nextPair.Competitor2.Name())
	game := ai.PlayGame(wordList, func(board *core.Board) *ai.Player {
		return ai.NewPlayer(nextPair.Competitor1)
	}, func(board *core.Board) *ai.Player {
		return ai.NewPlayer(nextPair.Competitor2)
	})
	log.Printf("Finished, saving game results with %d moves", len(game.Moves))

	return db.SaveGame(game)
}
