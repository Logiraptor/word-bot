package ai

import (
	"sync"
	"word-bot/core"
)

type SmartyAI struct {
	board       *core.Board
	wordList    core.WordList
	searchSpace WordTree
}

func NewSmartyAI(board *core.Board, wordList core.WordList, searchSpace WordTree) *SmartyAI {
	return &SmartyAI{
		board:       board,
		wordList:    wordList,
		searchSpace: searchSpace,
	}
}

func (b *SmartyAI) FindMoves(tiles []core.Tile) []ScoredMove {
	var moves = make(chan ScoredMove)
	var bestMove ScoredMove
	var wg sync.WaitGroup

	rack := core.NewConsumableRack(tiles)

	for i := 0; i < 15; i++ {
		wg.Add(1)
		go func(i int) {
			var localBestMove ScoredMove
			for j := 0; j < 15; j++ {

				b.Search(i, j, core.Horizontal, rack, b.searchSpace, nil, func(word []core.Tile) {
					if len(word) == 0 {
						return
					}

					if b.board.ValidateMove(word, i, j, core.Horizontal, b.wordList) {

						score := b.board.Score(word, i, j, core.Horizontal)

						if score > localBestMove.Score {
							newWord := make([]core.Tile, len(word))
							copy(newWord, word)

							current := ScoredMove{
								PlacedWord: core.PlacedWord{Word: newWord, Row: i, Col: j, Direction: core.Horizontal},
								Score:      score,
							}

							localBestMove = current
						}
					}
				})

				b.Search(i, j, core.Vertical, rack, b.searchSpace, nil, func(word []core.Tile) {
					if len(word) == 0 {
						return
					}

					if b.board.ValidateMove(word, i, j, core.Vertical, b.wordList) {
						newWord := make([]core.Tile, len(word))
						copy(newWord, word)

						current := ScoredMove{
							PlacedWord: core.PlacedWord{Word: newWord, Row: i, Col: j, Direction: core.Vertical},
							Score:      b.board.Score(newWord, i, j, core.Vertical),
						}

						if current.Score > localBestMove.Score {
							localBestMove = current
						}
					}
				})
			}

			if localBestMove.Word != nil {
				moves <- localBestMove
			}
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
		close(moves)
	}()

	for current := range moves {
		if current.Score > bestMove.Score {
			bestMove = current
		}
	}

	if bestMove.Word == nil {
		return nil
	}
	return []ScoredMove{bestMove}
}

type WordTree interface {
	IsTerminal() bool
	CanBranch(t core.Tile) (WordTree, bool)
}

func (s *SmartyAI) Search(i, j int, dir core.Direction, rack core.ConsumableRack, wordDB WordTree, prev []core.Tile, callback func([]core.Tile)) {
	dRow, dCol := dir.Offsets()
	if wordDB.IsTerminal() {
		callback(prev)
	}

	if s.board.OutOfBounds(i, j) {
		return
	}
	if s.board.HasTile(i, j) {
		letter := s.board.Cells[i][j].Tile
		if next, ok := wordDB.CanBranch(letter); ok {
			s.Search(i+dRow, j+dCol, dir, rack, next, prev, callback)
		}
	} else {
		for i, letter := range rack.Rack {
			if letter.IsBlank() {
				for r := 'a'; r <= 'z'; r++ {
					letter := core.Rune2Letter(r).ToTile(true)
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
