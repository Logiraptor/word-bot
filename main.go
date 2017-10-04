package main

import (
	"net/http"
	"os"

	"github.com/Logiraptor/word-bot/persist"

	"github.com/Logiraptor/word-bot/definitions"
	"github.com/Logiraptor/word-bot/web"
	"github.com/Logiraptor/word-bot/wordlist"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var wordDB *wordlist.Trie

func init() {
	builder := wordlist.NewTrieBuilder(151434)
	err := definitions.LoadWords("./words.txt", builder)
	if err != nil {
		panic(err)
	}

	wordDB = builder.Build()
}

func main() {
	var gdb *gorm.DB
	var err error

	// POSTGRES FORMAT
	// "host=myhost user=gorm dbname=gorm sslmode=disable password=mypassword"

	if os.Getenv("POSTGRES") == "" {
		gdb, err = gorm.Open("sqlite3", "results.db")
		if err != nil {
			panic(err)
		}
	} else {
		gdb, err = gorm.Open("postgres", os.Getenv("POSTGRES"))
		if err != nil {
			panic(err)
		}
	}

	db, err := persist.NewDB(gdb)
	if err != nil {
		panic(err)
	}

	s := web.Server{
		SearchSpace: wordDB,
		WordTree:    wordDB,
		DB:          db,
	}
	http.HandleFunc("/play", s.GetMove)
	http.HandleFunc("/render", s.RenderBoard)
	http.HandleFunc("/save", s.SaveGame)
	http.Handle("/", http.FileServer(http.Dir("frontend/public")))

	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
