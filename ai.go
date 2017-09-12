package main

import (
	"fmt"
	"sync"
)

type BruteForceAI struct {
	board *Board
}

func NewBruteForceAI(board *Board) *BruteForceAI {
	return &BruteForceAI{
		board: board,
	}
}

type ByScore []ScoredMove

func (b ByScore) Len() int           { return len(b) }
func (b ByScore) Less(i, j int) bool { return b[i].Score > b[j].Score }
func (b ByScore) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

type ScoredMove struct {
	PlacedWord
	Score
}

func (s ScoredMove) String() string {
	return fmt.Sprintf("(%s scores %d)", s.PlacedWord, s.Score)
}

func (b *BruteForceAI) FindMoves(rack []Tile) []ScoredMove {
	words := permute(rack)

	fmt.Println("Checking", len(words), "words")
	fmt.Println()

	validChan := make(chan PlacedWord, 100)
	wg := new(sync.WaitGroup)
	for i := 0; i < 15; i++ {
		wg.Add(1)
		go func(i int) {
			for j := 0; j < 15; j++ {
				for _, permutedWord := range words {
					if b.board.ValidateMove(permutedWord, i, j, Horizontal) {
						validChan <- PlacedWord{
							word:      permutedWord,
							row:       i,
							col:       j,
							direction: Horizontal,
						}
					}

					if b.board.ValidateMove(permutedWord, i, j, Vertical) {
						validChan <- PlacedWord{
							word:      permutedWord,
							row:       i,
							col:       j,
							direction: Vertical,
						}
					}
				}
			}
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
		close(validChan)
	}()

	var bestMove ScoredMove
	numMoves := 0
	for v := range validChan {
		numMoves++
		fmt.Print("\rFound ", numMoves, " valid moves")

		current := ScoredMove{
			PlacedWord: v,
			Score:      b.board.Score(v.word, v.row, v.col, v.direction),
		}

		if current.Score > bestMove.Score {
			bestMove = current
		}
	}

	fmt.Println()

	return []ScoredMove{bestMove}
}

func permute(rack []Tile) [][]Tile {
	if len(rack) == 0 {
		return [][]Tile{nil}
	}
	first := rack[0]
	rest := rack[1:]
	subPerm := permute(rest)
	output := make([][]Tile, len(subPerm), len(subPerm)*2)
	copy(output, subPerm)

	if first.IsBlank() {
		for option := Tile(0); option < 26; option++ {
			for _, perm := range subPerm {
				for i := range perm {
					output = append(output, append(append(perm[:i:i], option), perm[i:]...))
				}
				output = append(output, append(perm, option))
			}
		}
		return output
	}
	for _, perm := range subPerm {
		for i := range perm {
			output = append(output, append(append(perm[:i:i], first), perm[i:]...))
		}
		output = append(output, append(perm, first))
	}
	return output
}
