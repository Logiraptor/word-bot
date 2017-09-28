package main

import (
	"net/http"
	"os"
	"word-bot/definitions"
	"word-bot/web"
	"word-bot/wordlist"
)

var wordDB *wordlist.Trie

func init() {
	words, err := definitions.LoadWords("./words.txt")
	if err != nil {
		panic(err)
	}

	builder := wordlist.NewTrieBuilder()
	for _, word := range words {
		builder.AddWord(word)
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
