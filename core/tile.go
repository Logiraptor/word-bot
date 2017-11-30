package core

import (
	"fmt"
	"unicode"
)

// Tile represents an actual physical tile on the board
type Tile int

func (t Tile) ToRune() rune {
	return t.ToLetter().ToRune()
}

func (t Tile) String() string {
	if t.IsBlank() {
		return string(unicode.ToUpper(t.ToRune()))
	}
	return string(t.ToRune())
}

func (t Tile) GoString() string {
	return fmt.Sprintf("core.Rune2Letter(%q).ToTile(%#v)", t.ToRune(), t.IsBlank())
}

// PointValue returns the Score associated with a Tile
func (t Tile) PointValue() Score {
	if t.IsBlank() {
		return 0
	}
	return letterValues[t.ToLetter()]
}

// IsBlank returns true for blank tiles
func (t Tile) IsBlank() bool {
	return t&blankTileBit != 0
}

// ToLetter converts a tile to the letter it represents
func (t Tile) ToLetter() Letter {
	return Letter(t & letterMask)
}

// IsNoTile returns true if a tile is literally non existent
func (t Tile) IsNoTile() bool {
	return t == -1
}

// Flag returns the value of the ith flag.
func (t Tile) Flag(i uint) bool {
	return ((t&flagMask)>>flagOffset)&(1<<i) != 0
}

func (t Tile) SetFlag(i uint, value bool) Tile {
	if value {
		return t | (1 << (i + flagOffset))
	}
	return t & (^(1 << (i + flagOffset)))
}

const flagMask = 0xff00
const flagOffset = 8
const blankTileBit = 1 << 7
const letterMask = 1<<7 - 1

// Letter represents an abstract letter (but more efficient to use than a rune)
type Letter int

func (l Letter) ToRune() rune {
	return rune(l + 'a')
}

func (l Letter) ToTile(blank bool) Tile {
	if blank {
		return Tile(l | blankTileBit)
	}
	return Tile(l)
}

// Rune2Letter returns the Letter corresponding to a rune in ['a'..'z']
func Rune2Letter(r rune) Letter {
	return Letter(r - 'a')
}

func Tiles2String(tiles []Tile) string {
	word := ""
	for _, l := range tiles {
		word += l.String()
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
