package main

import (
	"math/rand"
	"strings"
)

// A-9, B-2, C-2, D-4, E-12, F-2, G-3, H-2, I-9, J-1, K-1, L-4, M-2, N-6, O-8, P-2, Q-1, R-6, S-4, T-6, U-4, V-2, W-2, X-1, Y-2, Z-1 and Blanks-2.

type Bag []Tile

func NewBag() *Bag {
	const allLetters = "aaaaaaaaabbccddddeeeeeeeeeeeeffggghhiiiiiiiiijkllllmmnnnnnnooooooooppqrrrrrrssssttttttuuuuvvwwxyyz"
	b := Bag(MakeTiles(MakeWord("aaaaaaaaabbccddddeeeeeeeeeeeeffggghhiiiiiiiiijkllllmmnnnnnnooooooooppqrrrrrrssssttttttuuuuvvwwxyyzaa"), strings.Repeat("x", len(allLetters))+"  "))
	return &b
}

func (b Bag) Shuffle() {
	for i := len(b) - 1; i > 0; i-- {
		j := rand.Intn(i)
		b[i], b[j] = b[j], b[i]
	}
}

func (b *Bag) Draw(n int) []Tile {
	if n > len(*b) {
		n = len(*b)
	}
	out := (*b)[:n:n]
	*b = (*b)[n:]
	return out
}

type Rack []Tile

func (r *Rack) Remove(tiles []Tile) {
outer:
	for _, tile := range tiles {
		for i := range *r {
			if (*r)[i] == tile {
				(*r)[i] = (*r)[len(*r)-1]
				*r = (*r)[:len(*r)-1]
				continue outer
			}
		}
	}
}
