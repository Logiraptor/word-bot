package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

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

type AI interface {
	FindMoves(rack []Tile) []ScoredMove
}

type Player struct {
	name  string
	rack  Rack
	score Score
}

func (p *Player) Play(ai AI, board *Board, bag *Bag) {
	moves := ai.FindMoves(p.rack)
	if len(moves) > 0 {
		bestMove := moves[0]
		fmt.Println(p.name, "would play:", bestMove)
		p.score += bestMove.Score

		used := board.PlaceTiles(bestMove.word, bestMove.row, bestMove.col, bestMove.direction)
		p.rack.Remove(used)
		p.rack = append(p.rack, Rack(bag.Draw(7-len(p.rack)))...)
		fmt.Println(p.name, "rack is now", p.rack)

		board.Print()
	} else {
		fmt.Println(p.name, "passes")
	}
}

type GameState struct {
	board   *Board
	bag     *Bag
	players []*Player
}

func main() {
	rand.Seed(time.Now().Unix())

	gs := new(GameState)
	gs.board = NewBoard()
	gs.board.Print()

	gs.bag = NewBag()
	gs.bag.Shuffle()

	bob := new(Player)
	bob.name = "bob"
	bob.rack = Rack(gs.bag.Draw(7))

	alice := new(Player)
	alice.name = "alice"
	alice.rack = Rack(gs.bag.Draw(7))

	// bf := NewBruteForceAI(gs.board)
	ai := NewSmartyAI(gs.board)

	gs.players = []*Player{bob, alice}

	for len(bob.rack) > 0 || len(alice.rack) > 0 {
		for _, player := range gs.players {
			player.Play(ai, gs.board, gs.bag)
		}

		fmt.Println("Bob:", bob.score, "Alice:", alice.score)

		gs.board.Print()
	}

	fmt.Println("Bob:", bob.score, "Alice:", alice.score)
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
