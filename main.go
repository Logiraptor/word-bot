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

func main() {
	words, err := loadWords()
	if err != nil {
		fmt.Println(err)
		return
	}

	wordDB := NewTrie()
	for _, word := range words {
		wordDB.AddWord(word)
	}

	board := NewBoard()
	board.PlaceTiles("hello", 7, 7, Horizontal)
	board.PlaceTiles("hello", 7, 7, Vertical)
	board.Print()
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
