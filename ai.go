package main

import (
	"fmt"
	"sort"
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

	validMoves := []PlacedWord{}

	cells := 15.0 * 15.0

	fmt.Println("Checking %d words", len(words))
	fmt.Println()

	for i := 0; i < 15; i++ {
		for j := 0; j < 15; j++ {
			for _, permutedWord := range words {
				if b.board.ValidateMove(permutedWord, i, j, Horizontal) {
					validMoves = append(validMoves, PlacedWord{
						word:      permutedWord,
						row:       i,
						col:       j,
						direction: Horizontal,
					})
				}

				if b.board.ValidateMove(permutedWord, i, j, Vertical) {
					validMoves = append(validMoves, PlacedWord{
						word:      permutedWord,
						row:       i,
						col:       j,
						direction: Vertical,
					})
				}
			}
			fmt.Print("\rFinished ", int((float64(i*15+j)/cells)*100), "%")
		}
	}

	fmt.Println()

	scoredMoves := []ScoredMove{}

	for _, validMove := range validMoves {
		scoredMoves = append(scoredMoves, ScoredMove{
			PlacedWord: validMove,
			Score:      b.board.Score(validMove.word, validMove.row, validMove.col, validMove.direction),
		})
	}

	sort.Sort(ByScore(scoredMoves))

	return scoredMoves
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
