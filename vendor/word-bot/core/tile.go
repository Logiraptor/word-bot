package core

// Tile represents an actual physical tile on the board
type Tile int

func (t Tile) ToRune() rune {
	return t.ToLetter().ToRune()
}

func (t Tile) String() string {
	if t.IsBlank() {
		return "_"
	}
	return string(t.ToRune())
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

// IsNoTile returns true if a tile is literally non existent
func (t Tile) IsNoTile() bool {
	return t == -1
}

const blankTileBit = 1 << 10
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

func tiles2String(tiles []Tile) string {
	word := ""
	for _, l := range tiles {
		word += string(l.ToRune())
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
