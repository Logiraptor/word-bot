package ai

import (
	"fmt"
	"sync"
	"word-bot/core"
)

type BruteForceAI struct {
	board    *core.Board
	wordList core.WordList
}

func NewBruteForceAI(board *core.Board, wordList core.WordList) *BruteForceAI {
	return &BruteForceAI{
		board:    board,
		wordList: wordList,
	}
}

type ByScore []ScoredMove

func (b ByScore) Len() int           { return len(b) }
func (b ByScore) Less(i, j int) bool { return b[i].Score > b[j].Score }
func (b ByScore) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }

type ScoredMove struct {
	core.PlacedWord
	core.Score
}

func (s ScoredMove) String() string {
	return fmt.Sprintf("(%s scores %d)", s.PlacedWord, s.Score)
}

func (b *BruteForceAI) FindMoves(rack []core.Tile) []ScoredMove {
	words := permute(rack)

	// fmt.Println("Checking", len(words), "words")
	// start := time.Now()
	fmt.Println()

	validChan := make(chan core.PlacedWord, 100)
	wg := new(sync.WaitGroup)
	for i := 0; i < 15; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			for j := 0; j < 15; j++ {
				if b.board.HasTile(i, j) {
					return
				}
				for _, permutedWord := range words {
					if len(permutedWord) == 0 {
						continue
					}

					if b.board.ValidateMove(permutedWord, i, j, core.Horizontal, b.wordList) {
						validChan <- core.PlacedWord{
							Word:      permutedWord,
							Row:       i,
							Col:       j,
							Direction: core.Horizontal,
						}
					}

					if b.board.ValidateMove(permutedWord, i, j, core.Vertical, b.wordList) {
						validChan <- core.PlacedWord{
							Word:      permutedWord,
							Row:       i,
							Col:       j,
							Direction: core.Vertical,
						}
					}
				}
			}
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

		current := ScoredMove{
			PlacedWord: v,
			Score:      b.board.Score(v.Word, v.Row, v.Col, v.Direction),
		}

		if current.Score > bestMove.Score {
			bestMove = current
		}
		// fmt.Print("\rFound ", numMoves, " valid moves. High score: ", bestMove)
	}

	// dur := time.Since(start)
	// fmt.Println()

	// fmt.Println("Finished in", dur)

	return []ScoredMove{bestMove}
}

func permute(rack []core.Tile) [][]core.Tile {
	if len(rack) == 0 {
		return [][]core.Tile{nil}
	}
	first := rack[0]
	rest := rack[1:]
	subPerm := permute(rest)
	output := make([][]core.Tile, len(subPerm), len(subPerm)*2)
	copy(output, subPerm)

	if first.IsBlank() {
		for option := 'a'; option <= 'z'; option++ {
			letter := core.Rune2Letter(option).ToTile(true)
			for _, perm := range subPerm {
				for i := range perm {
					output = append(output, append(append(perm[:i:i], letter), perm[i:]...))
				}
				output = append(output, append(perm, letter))
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
