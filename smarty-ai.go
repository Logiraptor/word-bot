package main

import (
	"fmt"
	"time"
)

type SmartyAI struct {
	board *Board
}

func NewSmartyAI(board *Board) *SmartyAI {
	return &SmartyAI{
		board: board,
	}
}

type ConsumableRack struct {
	rack     []Tile
	consumed int
}

func (c ConsumableRack) Consume(i int) ConsumableRack {
	return ConsumableRack{
		rack:     c.rack,
		consumed: c.consumed | (1 << uint(i)),
	}
}

func (c ConsumableRack) CanConsume(i int) bool {
	return c.consumed&(1<<uint(i)) == 0
}

func (b *SmartyAI) FindMoves(tiles []Tile) []ScoredMove {

	start := time.Now()
	fmt.Println()
	numMoves := 0
	var bestMove ScoredMove

	rack := ConsumableRack{rack: tiles, consumed: 0}

	for i := 0; i < 15; i++ {
		for j := 0; j < 15; j++ {
			b.Search(i, j, Horizontal, rack, wordDB, nil, func(word []Tile) {
				if b.board.ValidateMove(word, i, j, Horizontal) {
					newWord := make([]Tile, len(word))
					copy(newWord, word)

					numMoves++
					current := ScoredMove{
						PlacedWord: PlacedWord{newWord, i, j, Horizontal},
						Score:      b.board.Score(newWord, i, j, Horizontal),
					}

					if current.Score > bestMove.Score {
						bestMove = current
					}
					fmt.Print("\rFound ", numMoves, " valid moves. High score: ", bestMove)
				}
			})

			b.Search(i, j, Vertical, rack, wordDB, nil, func(word []Tile) {
				if b.board.ValidateMove(word, i, j, Vertical) {
					newWord := make([]Tile, len(word))
					copy(newWord, word)

					numMoves++
					current := ScoredMove{
						PlacedWord: PlacedWord{newWord, i, j, Vertical},
						Score:      b.board.Score(newWord, i, j, Vertical),
					}

					if current.Score > bestMove.Score {
						bestMove = current
					}
					fmt.Print("\rFound ", numMoves, " valid moves. High score: ", bestMove)
				}
			})
		}
	}

	dur := time.Since(start)
	fmt.Println()

	fmt.Println("Finished in", dur)

	return []ScoredMove{bestMove}
}

func (s *SmartyAI) Search(i, j int, dir Direction, rack ConsumableRack, wordDB *Trie, prev []Tile, callback func([]Tile)) {
	dRow, dCol := dir.Offsets()
	if wordDB.terminal {
		callback(prev)
	}

	if s.board.outOfBounds(i, j) {
		return
	}
	if s.board.HasTile(i, j) {
		letter := s.board.cells[i][j].tile
		if next, ok := wordDB.CanBranch(letter); ok {
			s.Search(i+dRow, j+dCol, dir, rack, next, prev, callback)
		}
	} else {
		for i, letter := range rack.rack {
			if letter.IsBlank() {
				for r := 'a'; r <= 'z'; r++ {
					letter := rune2Tile(r, true)
					if next, ok := wordDB.CanBranch(letter); ok && rack.CanConsume(i) {
						s.Search(i+dRow, j+dCol, dir, rack.Consume(i), next, append(prev, letter), callback)
					}
				}
			} else {
				if next, ok := wordDB.CanBranch(letter); ok && rack.CanConsume(i) {
					s.Search(i+dRow, j+dCol, dir, rack.Consume(i), next, append(prev, letter), callback)
				}
			}
		}
	}
}
