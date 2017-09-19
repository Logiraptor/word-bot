package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
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

func toTiles(word string) []Tile {
	return MakeTiles(MakeWord(word), strings.Repeat("x", len(word)))
}

type MoveRequest struct {
	Moves []Move   `json:"moves"`
	Rack  []TileJS `json:"rack"`
}

type TileJS struct {
	Letter string
	Blank  bool
	Value  Score
	Bonus  string
}

type Move struct {
	Tiles []TileJS `json:"tiles"`
	Row   int      `json:"row"`
	Col   int      `json:"col"`
	Dir   string   `json:"direction"` // vertical / horizontal
}

type ScoredMoveJS struct {
	Tiles []TileJS `json:"tiles"`
	Row   int      `json:"row"`
	Col   int      `json:"col"`
	Dir   string   `json:"direction"` // vertical / horizontal
	Score Score    `json:"score"`
}

type RenderedBoard struct {
	Board  [15][15]TileJS
	Scores []Score
}

func jsTilesToTiles(jsTiles []TileJS) []Tile {
	tiles := []Tile{}
	for _, t := range jsTiles {
		letters := []rune(t.Letter)
		letter := 'a'
		if len(letters) > 0 {
			letter = letters[0]
		}
		tiles = append(tiles, rune2Tile(letter, t.Blank))
	}
	return tiles
}

func tiles2JsTiles(tiles []Tile) []TileJS {
	jsTiles := []TileJS{}
	for _, t := range tiles {
		jsTiles = append(jsTiles, TileJS{
			Blank:  t.IsBlank(),
			Letter: string(tile2Rune(t)),
			Value:  t.PointValue(),
		})
	}
	return jsTiles
}

func getMove(rw http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	var moves MoveRequest
	err := json.NewDecoder(req.Body).Decode(&moves)
	if err != nil {
		http.Error(rw, "JSON parsing failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	b := NewBoard()
	for _, move := range moves.Moves {
		dir := Vertical
		if move.Dir == "horizontal" {
			dir = Horizontal
		}

		b.PlaceTiles(jsTilesToTiles(move.Tiles), move.Row, move.Col, dir)
	}

	ai := NewSmartyAI(b)
	play := ai.FindMoves(jsTilesToTiles(moves.Rack))[0]

	dirString := "horizontal"
	if play.direction == Vertical {
		dirString = "vertical"
	}

	json.NewEncoder(rw).Encode(ScoredMoveJS{
		Tiles: tiles2JsTiles(play.word),
		Row:   play.row,
		Col:   play.col,
		Dir:   dirString,
		Score: play.Score,
	})
}

func renderBoard(rw http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	var moves MoveRequest
	err := json.NewDecoder(req.Body).Decode(&moves)
	if err != nil {
		http.Error(rw, "JSON parsing failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	var output RenderedBoard
	output.Scores = make([]Score, len(moves.Moves))
	b := NewBoard()
	for i, move := range moves.Moves {
		dir := Vertical
		if move.Dir == "horizontal" {
			dir = Horizontal
		}
		output.Scores[i] = b.Score(jsTilesToTiles(move.Tiles), move.Row, move.Col, dir)
		b.PlaceTiles(jsTilesToTiles(move.Tiles), move.Row, move.Col, dir)
	}

	for i, row := range b.Cells {
		for j, cell := range row {
			if cell.Tile != NoTile {
				output.Board[i][j] = TileJS{
					Blank:  cell.Tile.IsBlank(),
					Letter: string(tile2Rune(cell.Tile)),
					Value:  cell.Tile.PointValue(),
					Bonus:  bonusToString(cell.Bonus),
				}
			} else {
				output.Board[i][j] = TileJS{
					Blank:  true,
					Letter: "",
					Value:  -1,
					Bonus:  bonusToString(cell.Bonus),
				}
			}
		}
	}

	json.NewEncoder(rw).Encode(output)
}

func main() {
	http.HandleFunc("/play", getMove)
	http.HandleFunc("/render", renderBoard)
	http.Handle("/", http.FileServer(http.Dir("public")))

	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}

func checkWords() {

	words, err := loadWords()
	if err != nil {
		panic(err)
	}

	start := 0
	for i := range words {
		if words[i] == "zorilla" {
			start = i
		}
	}

	output, err := os.OpenFile("fast-words.csv", os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		panic(err)
	}
	defer output.Close()

	type Result struct {
		Word, Definition, Error string
	}

	jobs := make(chan string, 100)
	results := make(chan Result, 100)
	for i := 0; i < 100; i++ {
		go func() {
			for j := range jobs {
				definition, err := defineWord(j)
				errString := ""
				if err != nil {
					errString = err.Error()
				}
				results <- Result{
					Word:       j,
					Error:      errString,
					Definition: definition,
				}
			}
		}()
	}

	go func() {
		for _, word := range words[start:] {
			jobs <- word
		}
		close(jobs)
	}()

	wr := csv.NewWriter(output)
	for res := range results {
		err = wr.Write([]string{
			res.Word,
			res.Definition,
			res.Error,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(res.Word, err == nil)
		wr.Flush()
	}
	wr.Flush()

}

func aiVsHumanGame() {
	rand.Seed(time.Now().Unix())

	var err error
	b := NewBoard()
	// f, err := os.Open("./game.json")
	// if err != nil {
	// 	panic(err)
	// }
	// err = json.NewDecoder(f).Decode(b)
	// if err != nil {
	// 	panic(err)
	// }
	ai := NewSmartyAI(b)
	var (
		word      string
		row, col  int
		direction Direction
	)
	for {
		err = b.Save("./game.json")
		if err != nil {
			panic(err)
		}
		b.Print()

		fmt.Println("Enter my tiles as a string with spaces for blanks")
		in := bufio.NewReader(os.Stdin)
		line, err := in.ReadString('\n')
		if err != nil {
			panic(err)
		}

		moves := ai.FindMoves(toTiles(line[:len(line)-1]))
		if len(moves) > 0 {
			move := moves[0]
			fmt.Printf("I play %s (%s) at %d, %d for %d points\n", tiles2String(move.word), move.direction, move.row, move.col, move.Score)
			b.PlaceTiles(move.word, move.row, move.col, move.direction)
		}

		err = b.Save("./game.json")
		if err != nil {
			panic(err)
		}
		b.Print()

		fmt.Println("Enter Opponent's move as: tiles row col vertical?")
		fmt.Scanln(&word, &row, &col, &direction)

		tiles := toTiles(word)
		b.PlaceTiles(tiles, row, col, direction)
	}
}

func aiVsAiGame() {
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

	for len(bob.rack) > 0 && len(alice.rack) > 0 {
		for _, player := range gs.players {
			player.Play(ai, gs.board, gs.bag)
		}

		fmt.Println("Bob:", bob.score, "Alice:", alice.score)

		gs.board.Print()
	}

	fmt.Println("GAME OVER")
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
