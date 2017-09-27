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

	wordDB = wordlist.NewTrie()
	for _, word := range words {
		wordDB.AddWord(word)
	}
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
