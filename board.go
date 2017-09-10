package main

import (
	"fmt"

	"github.com/fatih/color"
)

type Letter = int
type Bonus = byte
type Direction bool
type Score = int

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
	word      []Letter
	row, col  int
	direction Direction
}

func (p PlacedWord) String() string {
	word := letters2Word(p.word)
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
	bonus    Bonus
	value    Letter
	adjacent bool
}

type Board struct {
	cells [15][15]Cell
}

func NewBoard() *Board {
	b := new(Board)
	for i, row := range b.cells {
		for j := range row {
			b.cells[i][j].value = -1
			b.cells[i][j].bonus = normalBonus[i][j]
		}
	}
	return b
}

func (b *Board) Score(word string, row, col int, direction Direction) Score {
	words := b.FindNewWords(word, row, col, direction)
	total := 0
	for _, word := range words {
		wordScore := b.scoreWord(word.word, word.row, word.col, word.direction)
		total += wordScore
	}
	return total
}

func (b *Board) scoreWord(word []Letter, row, col int, direction Direction) Score {
	dRow, dCol := direction.Offsets()
	var sum = 0
	wordBonus := 1
	for i, letter := range word {
		letterBonus := 1
		if b.cells[row+(i*dRow)][col+(i*dCol)].value == -1 {
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
		sum += letterValues[letter] * letterBonus
	}
	return sum * wordBonus
}

func (b *Board) FindNewWords(word string, row, col int, direction Direction) []PlacedWord {
	dRow, dCol := direction.Offsets()
	words := []PlacedWord{}
	wordLetters := make([]Letter, len(word))

	for i, letter := range word {
		if b.cells[row+(dRow*i)][col+(dCol*i)].value == -1 {
			subWord, ok := b.GrowWord(letter, row+(dRow*i), col+(dCol*i), !direction)
			if ok {
				words = append(words, subWord)
			}
		}
		wordLetters[i] = letter2Token(letter)
	}

	words = append(words, PlacedWord{
		col: col, row: row,
		direction: direction, word: wordLetters,
	})

	return words
}

func (b *Board) GrowWord(r rune, row, col int, dir Direction) (PlacedWord, bool) {
	dRow, dCol := dir.Offsets()

	lhs := b.scan(row-dRow, col-dCol, -dRow, -dCol)
	rhs := b.scan(row+dRow, col+dCol, dRow, dCol)
	word := append(append(lhs, letter2Token(r)), rhs...)

	return PlacedWord{
		col:       col - len(lhs),
		row:       row,
		direction: dir,
		word:      word,
	}, len(word) > 1
}

func (b *Board) PlaceTiles(word string, row, col int, direction Direction) {
	dRow, dCol := direction.Offsets()

	for i, letter := range word {
		b.cells[row+i*dRow][col+i*dCol].value = letter2Token(letter)
	}
}

func (b *Board) Print() {
	for i, row := range b.cells {
		for j, cell := range row {
			_, _ = i, j

			letter := ' '
			if cell.value > -1 {
				letter = token2Letter(cell.value)
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

func (b *Board) scan(row, col, dRow, dCol int) []Letter {
	letters := []Letter{}
	for col > 0 && col < 15 &&
		row > 0 && row < 15 &&
		b.cells[row][col].value != -1 {
		letters = append(letters, b.cells[row][col].value)
		row += dRow
		col += dCol
	}
	return letters
}

func letter2Token(r rune) Letter {
	return int(r - 'a')
}

func token2Letter(t Letter) rune {
	return rune(t + 'a')
}

func letters2Word(letters []Letter) string {
	word := ""
	for _, l := range letters {
		word += string(token2Letter(l))
	}
	return word
}
