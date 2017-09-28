package main

import (
	"net/http"
	"os"

	"github.com/Logiraptor/word-bot/definitions"
	"github.com/Logiraptor/word-bot/web"
	"github.com/Logiraptor/word-bot/wordlist"
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
	s := web.Server{
		SearchSpace: wordDB,
		WordTree:    wordDB,
	}
	http.HandleFunc("/play", s.GetMove)
	http.HandleFunc("/render", s.RenderBoard)
	http.Handle("/", http.FileServer(http.Dir("frontend/public")))

	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
