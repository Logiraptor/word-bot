package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Logiraptor/word-bot/ai"
	"github.com/Logiraptor/word-bot/core"
)

type DB interface {
	Save([]core.PlacedTiles) error
}

type Server struct {
	WordTree    ai.WordTree
	SearchSpace core.WordList
	DB          DB
}

type AI interface {
	FindMoves(rack []core.Tile) []core.ScoredMove
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

func (t TileJS) ToTile() core.Tile {
	return core.Rune2Letter([]rune(t.Letter)[0]).ToTile(t.Blank)
}

type Move struct {
	Tiles []TileJS `json:"tiles"`
	Row   int      `json:"row"`
	Col   int      `json:"col"`
	Dir   string   `json:"direction"` // vertical / horizontal
}

func (m Move) ToPlacedTiles() core.PlacedTiles {
	dir := core.Horizontal
	if m.Dir == "vertical" {
		dir = core.Vertical
	}
	return core.PlacedTiles{
		Word:      jsTilesToTiles(m.Tiles),
		Row:       m.Row,
		Col:       m.Col,
		Direction: dir,
	}
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
		tiles = append(tiles, core.Rune2Letter(letter).ToTile(t.Blank))
	}
	return tiles
}

func tiles2JsTiles(tiles []core.Tile) []TileJS {
	jsTiles := []TileJS{}
	for _, t := range tiles {
		jsTiles = append(jsTiles, TileJS{
			Blank:  t.IsBlank(),
			Letter: string(t.ToRune()),
			Value:  t.PointValue(),
		})
	}
	return jsTiles
}

func (s Server) GetMove(rw http.ResponseWriter, req *http.Request) {
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
		b.PlaceTiles(move.ToPlacedTiles())
	}

	ai := ai.NewSmartyAI(b, s.SearchSpace, s.WordTree)
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

func (s Server) RenderBoard(rw http.ResponseWriter, req *http.Request) {
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
		pt := move.ToPlacedTiles()
		output.Scores[i] = b.Score(pt)
		b.PlaceTiles(pt)
	}

	for i, row := range b.Cells {
		for j, cell := range row {
			if !cell.Tile.IsNoTile() {
				output.Board[i][j] = TileJS{
					Blank:  cell.Tile.IsBlank(),
					Letter: string(cell.Tile.ToRune()),
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

func (s Server) SaveGame(rw http.ResponseWriter, req *http.Request) {
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

	placedWords := make([]core.PlacedTiles, 0, len(moves.Moves))
	for _, m := range moves.Moves {
		placedWords = append(placedWords, m.ToPlacedTiles())
	}

	err = s.DB.Save(placedWords)
	if err != nil {
		http.Error(rw, "Saving failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
