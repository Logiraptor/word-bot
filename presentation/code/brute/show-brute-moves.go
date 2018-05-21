package brute

import (
	"github.com/Logiraptor/word-bot/ai"
	"github.com/Logiraptor/word-bot/presentation/code/suggestions"
	"github.com/Logiraptor/word-bot/wordlist"
)

func main() {
	wordList := wordlist.MakeDefaultWordList()
	suggestions.PrintSuggestions(ai.NewBrute(wordList))
}
