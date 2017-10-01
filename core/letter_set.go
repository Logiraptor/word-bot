package core

// type TileSet [27]byte

// func NewEmptyTileSet() TileSet {
// 	return TileSet{}
// }

// func (l TileSet) indexOf(tile Tile) int {
// 	if tile.IsBlank() {
// 		return 26
// 	}
// 	return int(tile.ToLetter())
// }

// func (l TileSet) Consume(tile Tile) (TileSet, Tile) {
// 	if tile.IsBlank() {
// 		letter := tile.ToLetter().ToTile(false)
// 		if l[letter] > 0 {
// 			return l.Consume(letter)
// 		}
// 	}
// 	l[l.indexOf(tile)]--
// 	return l, tile
// }

// func (l TileSet) CanConsume(tile Tile) bool {
// 	return l[26] > 0 || l[l.indexOf(tile)] > 0
// }

// func (l TileSet) Add(tile Tile) TileSet {
// 	l[l.indexOf(tile)]++
// 	return l
// }

type TileSet [2]int64

func NewEmptyTileSet() TileSet {
	return TileSet{}
}

func indexOf(tile Tile) int {
	if tile.IsBlank() {
		return 26
	}
	return int(tile.ToLetter())
}

func (l TileSet) Consume(tile Tile) (TileSet, Tile) {
	if tile.IsBlank() {
		letter := tile.ToLetter().ToTile(false)
		if l[letter] > 0 {
			return l.Consume(letter)
		}
	}
	i := indexOf(tile)
	sub := i / 16
	rem := i % 16
	l[sub] -= 0x1 << (uint(rem) << 2)
	return l, tile
}

func (l TileSet) CanConsume(tile Tile) bool {
	if l[1]&(0xf<<(10<<2)) > 0 {
		return true
	}

	i := indexOf(tile)
	sub := i / 16
	rem := i % 16
	var bitMask int64 = 0xf << (uint(rem) << 2)
	return l[sub]&bitMask != 0
}

func (l *TileSet) Add(tile Tile) {
	i := indexOf(tile)
	sub := i / 16
	rem := i % 16
	l[sub] += 0x1 << (uint(rem) << 2)
}
