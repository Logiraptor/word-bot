package core

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
)

// WordList is used to validate words
type WordList interface {
	Contains(Word) bool
}

// Direction is either Horizontal or Verical
type Direction bool

func (d Direction) GoString() string {
	return fmt.Sprintf("core.%s", d.String())
}

// Score is a point value
type Score int

// Bonus is a point modifier
type Bonus = Score

// Word is a collection of letters
type Word []Letter

func (b Bonus) ToString() string {
	switch b {
	case DW:
		return "DW"
	case TW:
		return "TW"
	case TL:
		return "TL"
	case DL:
		return "DL"
	}
	return ""
}

// PlacedTiles represents a set of tiles placed on a board
type PlacedTiles struct {
	Word      []Tile
	Row, Col  int
	Direction Direction
}

func (p PlacedTiles) String() string {
	word := Tiles2String(p.Word)
	return fmt.Sprintf("(%d,%d,%v: %s)", p.Row, p.Col, p.Direction, word)
}

func (d Direction) String() string {
	if d == Horizontal {
		return "Horizontal"
	}
	return "Vertical"
}

// Offsets returns a direction vector for moving in the direction
func (d Direction) Offsets() (dRow, dCol int) {
	if d == Horizontal {
		return 0, 1
	}
	return 1, 0
}

// Cell represents a spot on the board
type Cell struct {
	Bonus Bonus
	Tile  Tile
}

// Board is a regular scrabble board
type Board struct {
	Cells [15][15]Cell
	ValidatedMoves []PlacedTiles
}

// NewBoard initializes an empty board
func NewBoard() *Board {
	b := new(Board)
	for i, row := range b.Cells {
		for j := range row {
			b.Cells[i][j].Tile = -1
			b.Cells[i][j].Bonus = normalBonus[i][j]
		}
	}
	b.ValidatedMoves = []PlacedTiles{}
	return b
}

func (b *Board) Clone() *Board {
	output := new(Board)
	for i, row := range b.Cells {
		for j := range row {
			output.Cells[i][j] = b.Cells[i][j]
		}
	}
	return output
}

// Save encodes a representation of the board to the given file
func (b *Board) Save(filename string) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0660)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(b)
}

// HasTile returns true if the given spot is occupied
func (b *Board) HasTile(row, col int) bool {
	if b.OutOfBounds(row, col) {
		return false
	}
	return !b.Cells[row][col].Tile.IsNoTile()
}

// ValidateMove returns true if the given move is legal
func (b *Board) ValidateMove(move PlacedTiles, wordList WordList) bool {
	b.ValidatedMoves = append(b.ValidatedMoves, move)
	// Check that it connects to other words
	connectsToOtherWords := false
	dRow, dCol := move.Direction.Offsets()
	wordPos := 0
	for progress := 0; wordPos < len(move.Word); progress++ {
		tileRow := move.Row + dRow*progress
		tileCol := move.Col + dCol*progress

		if b.OutOfBounds(tileRow, tileCol) {
			return false
		}
		if !b.HasTile(tileRow, tileCol) {
			wordPos++
		}

		if !connectsToOtherWords {
			if tileRow == 7 && tileCol == 7 {
				connectsToOtherWords = true
			} else if b.HasTile(tileRow-1, tileCol) ||
				b.HasTile(tileRow+1, tileCol) ||
				b.HasTile(tileRow, tileCol-1) ||
				b.HasTile(tileRow, tileCol+1) {
				connectsToOtherWords = true
			}
		}
	}

	if !connectsToOtherWords {
		return false
	}

	words := b.FindNewWords(move)
	for _, word := range words {
		if !wordList.Contains(tiles2Word(word.Word)) {
			return false
		}
	}
	return true
}

// Score computes the score for a given move. Score does not validate the move.
func (b *Board) Score(move PlacedTiles) Score {
	words := b.FindNewWords(move)

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic caught while scoring: %s which yielded %s\n", move, words)
			panic(r)
		}
	}()

	total := Score(0)
	for _, word := range words {
		wordScore := b.scoreWord(word)
		total += wordScore
	}
	return total
}

func (b *Board) OutOfBounds(row, col int) bool {
	return row < 0 || row >= 15 || col < 0 || col >= 15
}

func (b *Board) scoreWord(move PlacedTiles) Score {
	dRow, dCol := move.Direction.Offsets()
	sum := Bonus(0)
	wordBonus := Bonus(1)
	additionalBonus := Bonus(0)
	lettersUsed := 0
	for i, letter := range move.Word {
		tileRow := move.Row + (i * dRow)
		tileCol := move.Col + (i * dCol)
		if b.OutOfBounds(tileRow, tileCol) {
			panic(fmt.Sprintf("attempted to score word %s - OUT OF BOUNDS (moving [%d,%d])", move, dRow, dCol))
		}
		letterBonus := Bonus(1)
		if !b.HasTile(tileRow, tileCol) {
			bonus := b.Cells[tileRow][tileCol].Bonus
			switch bonus {
			case DoubleLetter:
				letterBonus *= 2
			case TripleLetter:
				letterBonus *= 3
			case TripleWord:
				wordBonus *= 3
			case DoubleWord:
				wordBonus *= 2
			}
			lettersUsed++
		}
		sum += letter.PointValue() * letterBonus
	}
	if lettersUsed >= 7 {
		additionalBonus = 50
	}
	return sum*wordBonus + additionalBonus
}

// FindNewWords returns all net-new words produced by the given move.
func (b *Board) FindNewWords(move PlacedTiles) []PlacedTiles {
	dRow, dCol := move.Direction.Offsets()
	words := []PlacedTiles{}
	progress := 0
	wordPos := 0

	var letters []Tile
	letters = b.scan(letters, move.Row-dRow, move.Col-dCol, -dRow, -dCol)
	reverse(letters)
	leftLen := len(letters)

	tileRow := move.Row + dRow*progress
	tileCol := move.Col + dCol*progress
	for !b.OutOfBounds(tileRow, tileCol) && wordPos < len(move.Word) {

		if b.HasTile(tileRow, tileCol) {
			letters = append(letters, b.Cells[tileRow][tileCol].Tile)
		} else {
			subWord, ok := b.GrowWord(move.Word[wordPos], tileRow, tileCol, !move.Direction)
			if ok {
				words = append(words, subWord)
			}
			letters = append(letters, move.Word[wordPos])
			wordPos++
		}

		progress++
		tileRow = move.Row + dRow*progress
		tileCol = move.Col + dCol*progress
	}

	// Grow placed word

	letters = b.scan(letters,
		move.Row+dRow*(len(letters)-leftLen),
		move.Col+dCol*(len(letters)-leftLen),
		dRow, dCol)

	words = append(words, PlacedTiles{
		Col:       move.Col - dCol*leftLen,
		Row:       move.Row - dRow*leftLen,
		Direction: move.Direction,
		Word:      letters,
	})

	return words
}

// GrowWord finds the full contiguous word at a given position
func (b *Board) GrowWord(l Tile, row, col int, dir Direction) (PlacedTiles, bool) {
	dRow, dCol := dir.Offsets()

	var letters []Tile
	letters = b.scan(letters, row-dRow, col-dCol, -dRow, -dCol)
	lenLeft := len(letters)
	reverse(letters)
	letters = append(letters, l)
	letters = b.scan(letters, row+dRow, col+dCol, dRow, dCol)

	return PlacedTiles{
		Col:       col - lenLeft*dCol,
		Row:       row - lenLeft*dRow,
		Direction: dir,
		Word:      letters,
	}, len(letters) > 1
}

// PlaceTiles places tiles. It does not validate the move.
func (b *Board) PlaceTiles(move PlacedTiles) []Tile {
	dRow, dCol := move.Direction.Offsets()
	progress := 0
	wordPos := 0
	for wordPos < len(move.Word) {
		tileRow := move.Row + progress*dRow
		tileCol := move.Col + progress*dCol

		if b.OutOfBounds(tileRow, tileCol) {
			panic(fmt.Sprintf("attempted to place word %s - OUT OF BOUNDS (moving [%d,%d]) (at [%d,%d])", move, dRow, dCol, tileRow, tileCol))
		}

		if !b.HasTile(tileRow, tileCol) {
			b.Cells[tileRow][tileCol].Tile = move.Word[wordPos]
			wordPos++
		}
		progress++
	}
	return move.Word
}

func (b *Board) NormalizeMove(move PlacedTiles) PlacedTiles {
	dRow, dCol := move.Direction.Offsets()
	i, j := move.Row, move.Col
	for b.HasTile(i, j) {
		i += dRow
		j += dCol
	}
	return PlacedTiles{
		Row:       i,
		Col:       j,
		Direction: move.Direction,
		Word:      move.Word,
	}
}

// Print prints the board to the console
func (b *Board) Print() {
	for i, row := range b.Cells {
		for j, cell := range row {
			letter := " "
			cellColor := color.New(color.FgBlack)
			if b.HasTile(i, j) {
				letter = cell.Tile.String()
				cellColor = cellColor.Add(color.BgMagenta)
			} else {
				cellColor = addBonusColor(cell, cellColor)
			}

			cellColor.Printf(" %s ", letter)
		}
		fmt.Println()
	}
}

func addBonusColor(cell Cell, cellColor *color.Color) *color.Color {
	switch cell.Bonus {
	case DoubleWord:
		cellColor = cellColor.Add(color.BgCyan)
	case TripleWord:
		cellColor = cellColor.Add(color.BgRed)
	case DoubleLetter:
		cellColor = cellColor.Add(color.BgBlue)
	case TripleLetter:
		cellColor = cellColor.Add(color.BgGreen)
	case None:
		cellColor = cellColor.Add(color.BgWhite)
	}
	return cellColor
}

func (b *Board) scan(letters []Tile, row, col, dRow, dCol int) []Tile {
	for b.HasTile(row, col) {
		letters = append(letters, b.Cells[row][col].Tile)
		row += dRow
		col += dCol
	}
	return letters
}

func reverse(tiles []Tile) {
	for i := 0; i < len(tiles)/2; i++ {
		j := len(tiles) - i - 1
		tiles[i], tiles[j] = tiles[j], tiles[i]
	}
}

// MakeWord converts a string to a Word
func MakeWord(word string) Word {
	output := make(Word, len(word))
	for i, r := range word {
		output[i] = Rune2Letter(r)
	}
	return output
}

// MakeTiles should be used like so: MakeTiles(word, "xx x") TODO: make this an example
func MakeTiles(word Word, mask string) []Tile {
	output := make([]Tile, len(word))
	for i, letter := range word {
		output[i] = letter.ToTile(mask[i] == ' ')
	}
	return output
}
