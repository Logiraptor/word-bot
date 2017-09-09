package main

import (
	"fmt"

	"github.com/fatih/color"
)

type Letter = int
type Bonus = byte
type Direction bool

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

var normalBonus [15][15]Bonus = [...][15]Bonus{
	[...]Bonus{TW, __, __, DL, __, __, __, TW, __, __, __, DL, __, __, TW},
	[...]Bonus{__, DW, __, __, __, TL, __, __, __, TL, __, __, __, DW, __},
	[...]Bonus{__, __, DW, __, __, __, DL, __, DL, __, __, __, DW, __, __},
	[...]Bonus{DL, __, __, DW, __, __, __, DL, __, __, __, DW, __, __, DL},
	[...]Bonus{__, __, __, __, DW, __, __, __, __, __, DW, __, __, __, __},
	[...]Bonus{__, TL, __, __, __, TL, __, __, __, TL, __, __, __, TL, __},
	[...]Bonus{__, __, DL, __, __, __, DL, __, DL, __, __, __, DL, __, __},
	[...]Bonus{TW, __, __, DL, __, __, __, DW, __, __, __, DL, __, __, TW},
	[...]Bonus{__, __, DL, __, __, __, DL, __, DL, __, __, __, DL, __, __},
	[...]Bonus{__, TL, __, __, __, TL, __, __, __, TL, __, __, __, TL, __},
	[...]Bonus{__, __, __, __, DW, __, __, __, __, __, DW, __, __, __, __},
	[...]Bonus{DL, __, __, DW, __, __, __, DL, __, __, __, DW, __, __, DL},
	[...]Bonus{__, __, DW, __, __, __, DL, __, DL, __, __, __, DW, __, __},
	[...]Bonus{__, DW, __, __, __, TL, __, __, __, TL, __, __, __, DW, __},
	[...]Bonus{TW, __, __, DL, __, __, __, TW, __, __, __, DL, __, __, TW},
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

func (b *Board) PlaceTiles(word string, row, col int, direction Direction) {
	if direction == Horizontal {
		for i, letter := range word {
			b.cells[row][col+i].value = int(letter - 'a')
		}
	} else {
		for i, letter := range word {
			b.cells[row+i][col].value = int(letter - 'a')
		}
	}
}

func (b *Board) Print() {
	for i, row := range b.cells {
		for j, cell := range row {
			_, _ = i, j

			letter := ' '
			if cell.value > -1 {
				letter = rune(cell.value) + 'a'
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
