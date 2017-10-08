package core

import "fmt"

// Rack manages a rack of tiles and allows efficient consumption of tiles
type Rack struct {
	Rack     []Tile
	consumed int
}

func NewConsumableRack(tiles []Tile) Rack {
	return Rack{
		Rack:     tiles,
		consumed: 0,
	}
}

// Consume uses up a tile in the rack
func (c Rack) Consume(i int) Rack {
	return Rack{
		Rack:     c.Rack,
		consumed: c.consumed | (1 << uint(i)),
	}
}

// CanConsume returns true if the ith tile is available to use
func (c Rack) CanConsume(i int) bool {
	return c.consumed&(1<<uint(i)) == 0
}

func (c Rack) Play(tiles []Tile) (Rack, bool) {
	newTiles := make([]Tile, len(c.Rack))
	copy(newTiles, c.Rack)

outer:
	for _, t := range tiles {
		for i := range newTiles {
			if tilesEqual(newTiles[i], t) {
				newTiles[i] = newTiles[len(newTiles)-1]
				newTiles = newTiles[:len(newTiles)-1]
				continue outer
			}
		}
		fmt.Printf("Cannot play tile '%s' with rack: %s, full rack is: %s\n", t.String(), tiles2String(newTiles), tiles2String(c.Rack))
		return Rack{}, false
	}

	return NewConsumableRack(newTiles), true
}

func tilesEqual(a, b Tile) bool {
	if a.IsBlank() && b.IsBlank() {
		return true
	}
	return a.ToLetter() == b.ToLetter()
}
