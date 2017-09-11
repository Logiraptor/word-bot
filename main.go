package main

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var wordDB *Trie

func init() {
	words, err := loadWords()
	if err != nil {
		panic(err)
	}

	wordDB = NewTrie()
	for _, word := range words {
		wordDB.AddWord(word)
	}
}

func main() {
	board := NewBoard()
	board.PlaceTiles(MakeTiles(MakeWord("hello"), "xxxxx"), 7, 7, Horizontal)
	board.PlaceTiles(MakeTiles(MakeWord("hello"), "xxxxx"), 7, 7, Vertical)
	board.Print()

	bob := NewBruteForceAI(board)
	moves := bob.FindMoves(MakeTiles(MakeWord("alfresc"), "xxxxxxx"))
	for _, move := range moves[:5] {
		fmt.Println("I would play:", move)
	}

}

func loadWords() ([]string, error) {
	f, err := os.Open("./words.txt")
	if err != nil {
		return nil, err
	}

	words := make([]string, 0, 80000)
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	return words, nil
}

func defineWord(word string) (string, error) {
	form := url.Values{
		"dictWord": {word},
	}
	resp, err := http.PostForm("https://scrabble.hasbro.com/en-us/tools", form)
	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return "", err
	}

	def := doc.Find(".word-definition")
	definition := def.Text()
	matcher := regexp.MustCompile("(?is)" + word + "(.*)")
	core := matcher.FindStringSubmatch(definition)
	if core == nil {
		return "", errors.New("No definition")
	}
	return strings.Join(strings.Fields(core[1]), " "), nil
}
