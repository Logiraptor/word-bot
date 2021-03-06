package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Logiraptor/word-bot/ai"
	"github.com/Logiraptor/word-bot/core"
	"github.com/Logiraptor/word-bot/wordlist"
)

type DB interface {
	Save([]core.PlacedTiles) error
}

type Server struct {
	WordTree    *wordlist.Trie
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
	Flags  []uint
}

func (t TileJS) ToTile() core.Tile {
	tile := core.Rune2Letter([]rune(t.Letter)[0]).ToTile(t.Blank)
	for _, f := range t.Flags {
		tile = tile.SetFlag(f, true)
	}
	return tile
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
		tile := core.Rune2Letter(letter).ToTile(t.Blank)
		for _, f := range t.Flags {
			tile = tile.SetFlag(f, true)
		}
		tiles = append(tiles, tile)
	}
	return tiles
}

func tiles2JsTiles(tiles []core.Tile) []TileJS {
	jsTiles := []TileJS{}
	for _, t := range tiles {
		jsTiles = append(jsTiles, tile2JsTile(t))
	}
	return jsTiles
}

func tile2JsTile(t core.Tile) TileJS {
	tJS := TileJS{
		Blank:  t.IsBlank(),
		Letter: string(t.ToRune()),
		Value:  t.PointValue(),
		Flags:  []uint{},
	}
	for i := uint(0); i < 8; i++ {
		if t.Flag(i) {
			tJS.Flags = append(tJS.Flags, i)
		}
	}
	return tJS
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
	bag := core.NewConsumableBag()
	bag = bag.ConsumeTiles(jsTilesToTiles(moves.Rack))
	for _, move := range moves.Moves {
		pt := move.ToPlacedTiles()
		b.PlaceTiles(pt)
		bag = bag.ConsumeTiles(pt.Word)
	}

	ai := ai.NewSmartyAI(s.SearchSpace, s.WordTree)
	defer ai.Kill()
	var play core.ScoredMove
	ai.FindMove(b, bag, core.NewConsumableRack(jsTilesToTiles(moves.Rack)), func(turn core.Turn) bool {
		if sm, ok := turn.(core.ScoredMove); ok {
			play = sm
		}
		return true
	})

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

	var output = Render(moves)

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

func Render(moves MoveRequest) RenderedBoard {
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
				output.Board[i][j] = tile2JsTile(cell.Tile)
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
	return output
}

func (s Server) ValidateEndpoint(rw http.ResponseWriter, req *http.Request) {
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

	var output = s.Validate(moves)

	json.NewEncoder(rw).Encode(output)
}

func (s Server) Validate(moves MoveRequest) []bool {
	output := make([]bool, len(moves.Moves))
	b := core.NewBoard()
	for i, m := range moves.Moves {
		output[i] = b.ValidateMove(m.ToPlacedTiles(), s.SearchSpace)
		b.PlaceTiles(m.ToPlacedTiles())
	}
	return output
}

func RemainingTiles(moves MoveRequest) []TileJS {
	b := core.NewConsumableBag()
	for _, m := range moves.Moves {
		b = b.ConsumeTiles(m.ToPlacedTiles().Word)
	}
	b = b.ConsumeTiles(jsTilesToTiles(moves.Rack))
	var remaining []core.Tile = b.Remaining()
	return tiles2JsTiles(remaining)
}
