package core

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
)

// ConsumableRack manages a rack of tiles and allows efficient consumption of tiles
type ConsumableRack struct {
	Rack     []Tile
	consumed int
}

func NewConsumableRack(tiles []Tile) ConsumableRack {
	return ConsumableRack{
		Rack:     tiles,
		consumed: 0,
	}
}

// Consume uses up a tile in the rack
func (c ConsumableRack) Consume(i int) ConsumableRack {
	return ConsumableRack{
		Rack:     c.Rack,
		consumed: c.consumed | (1 << uint(i)),
	}
}

// CanConsume returns true if the ith tile is available to use
func (c ConsumableRack) CanConsume(i int) bool {
	return c.consumed&(1<<uint(i)) == 0
}

// WordList is used to validate words
type WordList interface {
	Contains(Word) bool
}

// Tile represents an actual physical tile on the board
type Tile int

// Letter represents an abstract letter (but more efficient to use than a rune)
type Letter int

// Direction is either Horizontal or Verical
type Direction bool

// Score is a point value
type Score int

// Bonus is a point modifier
type Bonus = Score

// Word is a collection of letters
type Word []Letter

func bonusToString(b Bonus) string {
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

func (t Tile) String() string {
	if t.IsBlank() {
		return "_"
	}
	return string(tile2Rune(t))
}

// IsNoTile returns true if a tile is literally non existent
func (t Tile) IsNoTile() bool {
	return t == -1
}

const blankTileBit = 1 << 10
const letterMask = 1<<7 - 1

// Direction constants
const (
	Vertical   Direction = true
	Horizontal           = !Vertical
)

// Shorthands for bonuses
const (
	xx Bonus = iota
	DW
	TW
	DL
	TL
)

// Longer names for bonus spaces
const (
	None         = xx
	DoubleWord   = DW
	TripleWord   = TW
	DoubleLetter = DL
	TripleLetter = TL
)

// PlacedWord represents a set of tiles placed on a board
type PlacedWord struct {
	Word      []Tile
	Row, Col  int
	Direction Direction
}

func (p PlacedWord) String() string {
	word := tiles2String(p.Word)
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

// PointValue returns the Score associated with a Tile
func (t Tile) PointValue() Score {
	if t.IsBlank() {
		return 0
	}
	return letterValues[t]
}

// IsBlank returns true for blank tiles
func (t Tile) IsBlank() bool {
	return t&blankTileBit != 0
}

// ToLetter converts a tile to the letter it represents
func (t Tile) ToLetter() Letter {
	return Letter(t & letterMask)
}

var normalBonus = [...][15]Bonus{
	{TW, xx, xx, DL, xx, xx, xx, TW, xx, xx, xx, DL, xx, xx, TW},
	{xx, DW, xx, xx, xx, TL, xx, xx, xx, TL, xx, xx, xx, DW, xx},
	{xx, xx, DW, xx, xx, xx, DL, xx, DL, xx, xx, xx, DW, xx, xx},
	{DL, xx, xx, DW, xx, xx, xx, DL, xx, xx, xx, DW, xx, xx, DL},
	{xx, xx, xx, xx, DW, xx, xx, xx, xx, xx, DW, xx, xx, xx, xx},
	{xx, TL, xx, xx, xx, TL, xx, xx, xx, TL, xx, xx, xx, TL, xx},
	{xx, xx, DL, xx, xx, xx, DL, xx, DL, xx, xx, xx, DL, xx, xx},
	{TW, xx, xx, DL, xx, xx, xx, DW, xx, xx, xx, DL, xx, xx, TW},
	{xx, xx, DL, xx, xx, xx, DL, xx, DL, xx, xx, xx, DL, xx, xx},
	{xx, TL, xx, xx, xx, TL, xx, xx, xx, TL, xx, xx, xx, TL, xx},
	{xx, xx, xx, xx, DW, xx, xx, xx, xx, xx, DW, xx, xx, xx, xx},
	{DL, xx, xx, DW, xx, xx, xx, DL, xx, xx, xx, DW, xx, xx, DL},
	{xx, xx, DW, xx, xx, xx, DL, xx, DL, xx, xx, xx, DW, xx, xx},
	{xx, DW, xx, xx, xx, TL, xx, xx, xx, TL, xx, xx, xx, DW, xx},
	{TW, xx, xx, DL, xx, xx, xx, TW, xx, xx, xx, DL, xx, xx, TW},
}

// Cell represents a spot on the board
type Cell struct {
	Bonus Bonus
	Tile  Tile
}

// Board is a regular scrabble board
type Board struct {
	Cells [15][15]Cell
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
	return b
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
func (b *Board) ValidateMove(word []Tile, row, col int, direction Direction, wordList WordList) bool {

	// Check that it connects to other words
	connectsToOtherWords := false
	dRow, dCol := direction.Offsets()
	wordPos := 0
	for progress := 0; wordPos < len(word); progress++ {
		tileRow := row + dRow*progress
		tileCol := col + dCol*progress

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

	words := b.FindNewWords(word, row, col, direction)
	for _, word := range words {
		if !wordList.Contains(tiles2Word(word.Word)) {
			return false
		}
	}
	return true
}

// Score computes the score for a given move. Score does not validate the move.
func (b *Board) Score(word []Tile, row, col int, direction Direction) Score {
	words := b.FindNewWords(word, row, col, direction)

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic caught while scoring: %s which yielded %s\n", PlacedWord{word, row, col, direction}, words)
			panic(r)
		}
	}()

	total := Score(0)
	for _, word := range words {
		wordScore := b.scoreWord(word.Word, word.Row, word.Col, word.Direction)
		total += wordScore
	}
	return total
}

func (b *Board) OutOfBounds(row, col int) bool {
	return row < 0 || row >= 15 || col < 0 || col >= 15
}

func (b *Board) scoreWord(word []Tile, row, col int, direction Direction) Score {
	dRow, dCol := direction.Offsets()
	sum := Bonus(0)
	wordBonus := Bonus(1)
	additionalBonus := Bonus(0)
	lettersUsed := 0
	for i, letter := range word {
		tileRow := row + (i * dRow)
		tileCol := col + (i * dCol)
		if b.OutOfBounds(tileRow, tileCol) {
			panic(fmt.Sprintf("attempted to score word %s - OUT OF BOUNDS (moving [%d,%d])", PlacedWord{word, row, col, direction}, dRow, dCol))
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
func (b *Board) FindNewWords(word []Tile, row, col int, direction Direction) []PlacedWord {
	dRow, dCol := direction.Offsets()
	words := []PlacedWord{}
	wordLetters := make([]Tile, 0, len(word))
	progress := 0
	wordPos := 0

	tileRow := row + dRow*progress
	tileCol := col + dCol*progress
	for !b.OutOfBounds(tileRow, tileCol) && wordPos < len(word) {

		if b.HasTile(tileRow, tileCol) {
			wordLetters = append(wordLetters, b.Cells[tileRow][tileCol].Tile)
		} else {
			subWord, ok := b.GrowWord(word[wordPos], tileRow, tileCol, !direction)
			if ok {
				words = append(words, subWord)
			}
			wordLetters = append(wordLetters, word[wordPos])
			wordPos++
		}

		progress++
		tileRow = row + dRow*progress
		tileCol = col + dCol*progress
	}

	// Grow placed word
	lhs := b.scan(row-dRow, col-dCol, -dRow, -dCol)
	reverse(lhs)
	rhs := b.scan(
		row+dRow*len(wordLetters),
		col+dCol*len(wordLetters),
		dRow, dCol)
	wordLetters = append(append(lhs, wordLetters...), rhs...)

	words = append(words, PlacedWord{
		Col: col - dCol*len(lhs), Row: row - dRow*len(lhs),
		Direction: direction, Word: wordLetters,
	})

	return words
}

// GrowWord finds the full contiguous word at a given position
func (b *Board) GrowWord(l Tile, row, col int, dir Direction) (PlacedWord, bool) {
	dRow, dCol := dir.Offsets()

	lhs := b.scan(row-dRow, col-dCol, -dRow, -dCol)
	reverse(lhs)
	rhs := b.scan(row+dRow, col+dCol, dRow, dCol)
	word := append(append(lhs, l), rhs...)

	return PlacedWord{
		Col:       col - len(lhs)*dCol,
		Row:       row - len(lhs)*dRow,
		Direction: dir,
		Word:      word,
	}, len(word) > 1
}

// PlaceTiles places tiles. It does not validate the move.
func (b *Board) PlaceTiles(tiles []Tile, row, col int, direction Direction) []Tile {
	dRow, dCol := direction.Offsets()
	progress := 0
	wordPos := 0
	for wordPos < len(tiles) {
		tileRow := row + progress*dRow
		tileCol := col + progress*dCol

		if b.OutOfBounds(tileRow, tileCol) {
			panic(fmt.Sprintf("attempted to place word %s - OUT OF BOUNDS (moving [%d,%d]) (at [%d,%d])",
				PlacedWord{tiles, row, col, direction}, dRow, dCol, tileRow, tileCol))
		}

		if !b.HasTile(tileRow, tileCol) {
			b.Cells[tileRow][tileCol].Tile = tiles[wordPos]
			wordPos++
		}
		progress++
	}
	return tiles
}

// Print prints the board to the console
func (b *Board) Print() {
	for i, row := range b.Cells {
		for j, cell := range row {
			_, _ = i, j

			letter := ' '
			cellColor := color.New(color.FgBlack)
			if b.HasTile(i, j) {
				letter = tile2Rune(cell.Tile)
				cellColor = cellColor.Add(color.BgMagenta)
			} else {
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
			}

			cellColor.Printf(" %c ", letter)
		}
		fmt.Println()
	}
}

func (b *Board) scan(row, col, dRow, dCol int) []Tile {
	letters := []Tile{}
	for col >= 0 && col < 15 &&
		row >= 0 && row < 15 &&
		b.HasTile(row, col) {
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

// Rune2Letter returns the Letter corresponding to a rune in ['a'..'z']
func Rune2Letter(r rune) Letter {
	return Letter(r - 'a')
}

// Rune2Tile constructs a tile
func Rune2Tile(r rune, blank bool) Tile {
	return letter2Tile(Rune2Letter(r), blank)
}

func letter2Tile(l Letter, blank bool) Tile {
	if blank {
		return Tile(l | blankTileBit)
	}
	return Tile(l)
}

func letter2Rune(t Letter) rune {
	return rune(t + 'a')
}

func tile2Rune(t Tile) rune {
	return letter2Rune(t.ToLetter())
}

func tiles2String(tiles []Tile) string {
	word := ""
	for _, l := range tiles {
		word += string(tile2Rune(l))
	}
	return word
}

func tiles2Word(tiles []Tile) Word {
	word := make(Word, len(tiles))
	for i, l := range tiles {
		word[i] = l.ToLetter()
	}
	return word
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
		output[i] = letter2Tile(letter, mask[i] == ' ')
	}
	return output
}
