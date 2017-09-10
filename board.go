package main

import (
	"fmt"

	"github.com/fatih/color"
)

type Tile int
type Letter int
type Direction bool
type Score int
type Bonus = Score
type Word []Letter

const NoTile = -1
const BlankTileBit = 1 << 10
const LetterMask = 1<<7 - 1

const (
	Vertical   Direction = true
	Horizontal           = !Vertical
)

const (
	__ Bonus = iota
	DW
	TW
	DL
	TL
)

const (
	None         = __
	DoubleWord   = DW
	TripleWord   = TW
	DoubleLetter = DL
	TripleLetter = TL
)

type PlacedWord struct {
	word      []Tile
	row, col  int
	direction Direction
}

func (p PlacedWord) String() string {
	word := tiles2String(p.word)
	return fmt.Sprintf("(%d,%d,%v: %s)", p.row, p.col, p.direction, word)
}

func (d Direction) String() string {
	if d == Horizontal {
		return "Horizontal"
	}
	return "Vertical"
}

func (d Direction) Offsets() (dRow, dCol int) {
	if d == Horizontal {
		return 0, 1
	}
	return 1, 0
}

func (t Tile) PointValue() Score {
	if t.IsBlank() {
		return 0
	}
	return letterValues[t]
}

func (t Tile) IsBlank() bool {
	return t&BlankTileBit != 0
}

func (t Tile) ToLetter() Letter {
	if t.IsBlank() {
		return Letter(t & LetterMask)
	}
	return Letter(t & LetterMask)
}

var normalBonus [15][15]Bonus = [...][15]Bonus{
	{TW, __, __, DL, __, __, __, TW, __, __, __, DL, __, __, TW},
	{__, DW, __, __, __, TL, __, __, __, TL, __, __, __, DW, __},
	{__, __, DW, __, __, __, DL, __, DL, __, __, __, DW, __, __},
	{DL, __, __, DW, __, __, __, DL, __, __, __, DW, __, __, DL},
	{__, __, __, __, DW, __, __, __, __, __, DW, __, __, __, __},
	{__, TL, __, __, __, TL, __, __, __, TL, __, __, __, TL, __},
	{__, __, DL, __, __, __, DL, __, DL, __, __, __, DL, __, __},
	{TW, __, __, DL, __, __, __, DW, __, __, __, DL, __, __, TW},
	{__, __, DL, __, __, __, DL, __, DL, __, __, __, DL, __, __},
	{__, TL, __, __, __, TL, __, __, __, TL, __, __, __, TL, __},
	{__, __, __, __, DW, __, __, __, __, __, DW, __, __, __, __},
	{DL, __, __, DW, __, __, __, DL, __, __, __, DW, __, __, DL},
	{__, __, DW, __, __, __, DL, __, DL, __, __, __, DW, __, __},
	{__, DW, __, __, __, TL, __, __, __, TL, __, __, __, DW, __},
	{TW, __, __, DL, __, __, __, TW, __, __, __, DL, __, __, TW},
}

type Cell struct {
	bonus Bonus
	tile  Tile
}

type Board struct {
	cells [15][15]Cell
}

func NewBoard() *Board {
	b := new(Board)
	for i, row := range b.cells {
		for j := range row {
			b.cells[i][j].tile = NoTile
			b.cells[i][j].bonus = normalBonus[i][j]
		}
	}
	return b
}

func (b *Board) Score(word []Tile, row, col int, direction Direction) Score {
	words := b.FindNewWords(word, row, col, direction)
	total := Score(0)
	for _, word := range words {
		wordScore := b.scoreWord(word.word, word.row, word.col, word.direction)
		total += wordScore
	}
	return total
}

func (b *Board) scoreWord(word []Tile, row, col int, direction Direction) Score {
	dRow, dCol := direction.Offsets()
	sum := Bonus(0)
	wordBonus := Bonus(1)
	for i, letter := range word {
		letterBonus := Bonus(1)
		if b.cells[row+(i*dRow)][col+(i*dCol)].tile == NoTile {
			bonus := b.cells[row+(i*dRow)][col+(i*dCol)].bonus
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
		}
		sum += letter.PointValue() * letterBonus
	}
	return sum * wordBonus
}

func (b *Board) FindNewWords(word []Tile, row, col int, direction Direction) []PlacedWord {
	dRow, dCol := direction.Offsets()
	words := []PlacedWord{}
	wordLetters := make([]Tile, len(word))

	for i, letter := range word {
		if b.cells[row+(dRow*i)][col+(dCol*i)].tile == NoTile {
			subWord, ok := b.GrowWord(letter, row+(dRow*i), col+(dCol*i), !direction)
			if ok {
				words = append(words, subWord)
			}
		}
		wordLetters[i] = letter
	}

	// Grow placed word
	lhs := b.scan(row-dRow, col-dCol, -dRow, -dCol)
	rhs := b.scan(
		row+dRow*len(wordLetters),
		col+dCol*len(wordLetters),
		dRow, dCol)
	wordLetters = append(append(lhs, wordLetters...), rhs...)

	words = append(words, PlacedWord{
		col: col - dCol*len(lhs), row: row - dRow*len(lhs),
		direction: direction, word: wordLetters,
	})

	return words
}

func (b *Board) GrowWord(l Tile, row, col int, dir Direction) (PlacedWord, bool) {
	dRow, dCol := dir.Offsets()

	lhs := b.scan(row-dRow, col-dCol, -dRow, -dCol)
	rhs := b.scan(row+dRow, col+dCol, dRow, dCol)
	word := append(append(lhs, l), rhs...)

	return PlacedWord{
		col:       col - len(lhs),
		row:       row,
		direction: dir,
		word:      word,
	}, len(word) > 1
}

func (b *Board) PlaceTiles(tiles []Tile, row, col int, direction Direction) {
	dRow, dCol := direction.Offsets()

	for i, tile := range tiles {
		b.cells[row+i*dRow][col+i*dCol].tile = tile
	}
}

func (b *Board) Print() {
	for i, row := range b.cells {
		for j, cell := range row {
			_, _ = i, j

			letter := ' '
			if cell.tile != NoTile {
				letter = tile2Rune(cell.tile)
			}
			cellColor := color.New(color.FgBlack)

			switch cell.bonus {
			case DoubleWord:
				cellColor = cellColor.Add(color.BgRed)
			case TripleWord:
				cellColor = cellColor.Add(color.BgRed)
			case DoubleLetter:
				cellColor = cellColor.Add(color.BgBlue)
			case TripleLetter:
				cellColor = cellColor.Add(color.BgBlue)
			case None:
				cellColor = cellColor.Add(color.BgWhite)
			}

			cellColor.Printf(" %c ", letter)
		}
		fmt.Println()
	}
}

func (b *Board) scan(row, col, dRow, dCol int) []Tile {
	letters := []Tile{}
	for col > 0 && col < 15 &&
		row > 0 && row < 15 &&
		b.cells[row][col].tile != NoTile {
		letters = append(letters, b.cells[row][col].tile)
		row += dRow
		col += dCol
	}
	return letters
}

func rune2Letter(r rune) Letter {
	return Letter(r - 'a')
}

func rune2Tile(r rune, blank bool) Tile {
	return letter2Tile(rune2Letter(r), blank)
}

func letter2Tile(l Letter, blank bool) Tile {
	if blank {
		return Tile(l | BlankTileBit)
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

func MakeWord(word string) Word {
	output := make(Word, len(word))
	for i, r := range word {
		output[i] = rune2Letter(r)
	}
	return output
}

// MakeTiles should be used like: MakeTiles(word, "xx x")
func MakeTiles(word Word, mask string) []Tile {
	output := make([]Tile, len(word))
	for i, letter := range word {
		output[i] = letter2Tile(letter, mask[i] == ' ')
	}
	return output
}
