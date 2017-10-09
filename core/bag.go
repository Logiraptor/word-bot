package core

import (
	"math/bits"
	"math/rand"
	"strings"
)

// A-9, B-2, C-2, D-4, E-12, F-2, G-3, H-2, I-9, J-1, K-1, L-4, M-2, N-6, O-8, P-2, Q-1, R-6, S-4, T-6, U-4, V-2, W-2, X-1, Y-2, Z-1 and Blanks-2.

const allLetters = "aaaaaaaaabbccddddeeeeeeeeeeeeffggghhiiiiiiiiijkllllmmnnnnnnooooooooppqrrrrrrssssttttttuuuuvvwwxyyz"

var allTiles = MakeTiles(MakeWord(allLetters+"aa"), strings.Repeat("x", len(allLetters))+"  ")

type Bag struct {
	tiles    []Tile
	consumed [2]uint64
}

func NewConsumableBag() Bag {
	allTilesCopy := make([]Tile, len(allTiles))
	copy(allTilesCopy, allTiles)
	b := Bag{
		tiles: allTilesCopy,
	}
	return b
}

func (c Bag) validate() {
	if (c.consumed[1] & (1 << (uint(len(allTiles)) - 64))) != 0 {
		panic("Invalid bag")
	}
}

// Shuffle randomizes the order of tiles inside the bag
func (c Bag) Shuffle() Bag {
	c.validate()
	result := c
	result.tiles = make([]Tile, len(allTiles))
	copy(result.tiles, allTiles)
	for i := len(result.tiles) - 1; i > 0; i-- {
		j := rand.Intn(i)
		result.tiles[i], result.tiles[j] = result.tiles[j], result.tiles[i]

		subFieldI := i / 64
		remI := i % 64
		bitI := (result.consumed[subFieldI] & (1 << uint(remI))) >> uint(remI)

		subFieldJ := j / 64
		remJ := j % 64
		bitJ := (result.consumed[subFieldJ] & (1 << uint(remJ))) >> uint(remJ)

		result.consumed[subFieldI] ^= (-bitJ ^ result.consumed[subFieldI]) & (1 << uint(remJ))
		result.consumed[subFieldJ] ^= (-bitI ^ result.consumed[subFieldJ]) & (1 << uint(remI))
	}
	result.validate()
	return result
}

// Consume uses up a tile in the rack
func (c Bag) Consume(i int) Bag {
	c.validate()
	result := c
	subField := i / 64
	rem := i % 64
	result.consumed[subField] |= (1 << uint(rem))
	result.validate()
	return result
}

// CanConsume returns true if the ith tile is available to use
func (c Bag) CanConsume(i int) bool {
	c.validate()
	subField := i / 64
	rem := i % 64
	return c.consumed[subField]&(1<<uint(rem)) == 0
}

func (c Bag) FillRack(tiles []Tile, n int) (Bag, []Tile) {
	c.validate()
	i := 0
	for x := 0; x < len(c.tiles) && i < n; x++ {
		if c.CanConsume(x) {
			tiles = append(tiles, c.tiles[x])
			c = c.Consume(x)
			i++
		}
	}
	c.validate()
	return c, tiles
}

func (c Bag) ConsumeTiles(tiles []Tile) Bag {
	c.validate()
	for _, t := range tiles {
		for i, x := range c.tiles {
			if tilesEqual(x, t) && c.CanConsume(i) {
				c = c.Consume(i)
				break
			}
		}
	}
	c.validate()
	return c
}

func (c Bag) Count() int {
	c.validate()
	consumed := bits.OnesCount64(c.consumed[0]) + bits.OnesCount64(c.consumed[1])
	return len(allTiles) - consumed
}
