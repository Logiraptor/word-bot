package main

import (
	"net/http"
	"os"

	"github.com/Logiraptor/word-bot/web"
	"github.com/Logiraptor/word-bot/wordlist"
)

var wordDB *wordlist.Trie

func init() {
	wordDB = wordlist.MakeDefaultWordList()
}

func main() {
	s := web.Server{
		SearchSpace: wordDB,
		WordTree:    wordDB,
	}
	http.HandleFunc("/play", s.GetMove)
	http.HandleFunc("/validate", s.ValidateEndpoint)
	http.HandleFunc("/render", s.RenderBoard)
	http.HandleFunc("/save", s.SaveGame)
	http.Handle("/", http.FileServer(http.Dir("frontend/public")))

	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
