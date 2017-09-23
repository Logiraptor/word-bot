package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"word-bot/ai"
	"word-bot/core"
	"word-bot/wordlist"

	"github.com/PuerkitoBio/goquery"
)

var wordDB *wordlist.Trie

func init() {
	words, err := loadWords()
	if err != nil {
		panic(err)
	}

	wordDB = wordlist.NewTrie()
	for _, word := range words {
		wordDB.AddWord(word)
	}
}

type AI interface {
	FindMoves(rack []core.Tile) []ai.ScoredMove
}

func toTiles(word string) []core.Tile {
	return core.MakeTiles(core.MakeWord(word), strings.Repeat("x", len(word)))
}

type MoveRequest struct {
	Moves []Move   `json:"moves"`
	Rack  []TileJS `json:"rack"`
}

type TileJS struct {
	Letter string
	Blank  bool
	Value  core.Score
	Bonus  string
}

type Move struct {
	Tiles []TileJS `json:"tiles"`
	Row   int      `json:"row"`
	Col   int      `json:"col"`
	Dir   string   `json:"direction"` // vertical / horizontal
}

type ScoredMoveJS struct {
	Tiles []TileJS   `json:"tiles"`
	Row   int        `json:"row"`
	Col   int        `json:"col"`
	Dir   string     `json:"direction"` // vertical / horizontal
	Score core.Score `json:"score"`
}

type RenderedBoard struct {
	Board  [15][15]TileJS
	Scores []core.Score
}

func jsTilesToTiles(jsTiles []TileJS) []core.Tile {
	tiles := []core.Tile{}
	for _, t := range jsTiles {
		letters := []rune(t.Letter)
		letter := 'a'
		if len(letters) > 0 {
			letter = letters[0]
		}
		tiles = append(tiles, core.Rune2Tile(letter, t.Blank))
	}
	return tiles
}

func tiles2JsTiles(tiles []core.Tile) []TileJS {
	jsTiles := []TileJS{}
	for _, t := range tiles {
		jsTiles = append(jsTiles, TileJS{
			Blank:  t.IsBlank(),
			Letter: string(core.Tile2Rune(t)),
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

	b := core.NewBoard()
	for _, move := range moves.Moves {
		dir := core.Vertical
		if move.Dir == "horizontal" {
			dir = core.Horizontal
		}

		b.PlaceTiles(jsTilesToTiles(move.Tiles), move.Row, move.Col, dir)
	}

	ai := ai.NewSmartyAI(b, wordDB, wordDB)
	play := ai.FindMoves(jsTilesToTiles(moves.Rack))[0]

	dirString := "horizontal"
	if play.Direction == core.Vertical {
		dirString = "vertical"
	}

	json.NewEncoder(rw).Encode(ScoredMoveJS{
		Tiles: tiles2JsTiles(play.Word),
		Row:   play.Row,
		Col:   play.Col,
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
	output.Scores = make([]core.Score, len(moves.Moves))
	b := core.NewBoard()
	for i, move := range moves.Moves {
		dir := core.Vertical
		if move.Dir == "horizontal" {
			dir = core.Horizontal
		}
		output.Scores[i] = b.Score(jsTilesToTiles(move.Tiles), move.Row, move.Col, dir)
		b.PlaceTiles(jsTilesToTiles(move.Tiles), move.Row, move.Col, dir)
	}

	for i, row := range b.Cells {
		for j, cell := range row {
			if !cell.Tile.IsNoTile() {
				output.Board[i][j] = TileJS{
					Blank:  cell.Tile.IsBlank(),
					Letter: string(core.Tile2Rune(cell.Tile)),
					Value:  cell.Tile.PointValue(),
					Bonus:  cell.Bonus.ToString(),
				}
			} else {
				output.Board[i][j] = TileJS{
					Blank:  true,
					Letter: "",
					Value:  -1,
					Bonus:  cell.Bonus.ToString(),
				}
			}
		}
	}

	json.NewEncoder(rw).Encode(output)
}

func main() {
	http.HandleFunc("/play", getMove)
	http.HandleFunc("/render", renderBoard)
	http.Handle("/", http.FileServer(http.Dir("frontend/public")))

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
