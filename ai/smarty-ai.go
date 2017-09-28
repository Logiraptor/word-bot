package ai

import (
	"sync"
	"sync/atomic"

	"github.com/Logiraptor/word-bot/core"
)

type SmartyAI struct {
	board       *core.Board
	wordList    core.WordList
	jobs        chan<- Job
	searchSpace WordTree
}

func NewSmartyAI(board *core.Board, wordList core.WordList, searchSpace WordTree) *SmartyAI {

	jobs := make(chan Job)

	s := &SmartyAI{
		board:       board,
		wordList:    wordList,
		searchSpace: searchSpace,
		jobs:        jobs,
	}

	for i := 0; i < 10; i++ {
		go searchWorker(s, jobs)
	}

	return s
}

type Job struct {
	i, j       int
	dir        core.Direction
	rack       core.ConsumableRack
	wordDB     WordTree
	resultChan chan<- Result
}

type Result struct {
	word []core.Tile
}

func searchWorker(s *SmartyAI, jobs <-chan Job) {
	for job := range jobs {
		s.Search(job.i, job.j, job.dir, job.rack, job.wordDB, nil, func(word []core.Tile) {
			job.resultChan <- Result{
				word: word,
			}
		})
		close(job.resultChan)
	}
}

func (b *SmartyAI) FindMoves(tiles []core.Tile) []ScoredMove {
	var moves = make(chan ScoredMove)
	var bestMove ScoredMove
	var wg sync.WaitGroup

	var badMoves uint64

	rack := core.NewConsumableRack(tiles)

	dirs := []core.Direction{core.Horizontal, core.Vertical}

	for i := 0; i < 15; i++ {
		wg.Add(1)
		go func(i int) {
			var localBestMove ScoredMove
			for j := 0; j < 15; j++ {
				for _, dir := range dirs {

					b.Search(i, j, dir, rack, b.searchSpace, nil, func(word []core.Tile) {
						if len(word) == 0 {
							return
						}

						if b.board.ValidateMove(word, i, j, dir, b.wordList) {
							score := b.board.Score(word, i, j, dir)

							if score > localBestMove.Score {
								newWord := make([]core.Tile, len(word))
								copy(newWord, word)

								current := ScoredMove{
									PlacedWord: core.PlacedWord{Word: newWord, Row: i, Col: j, Direction: dir},
									Score:      score,
								}

								localBestMove = current
							}
						} else {
							atomic.AddUint64(&badMoves, 1)
						}
					})

				}
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

	// fmt.Println(badMoves, "invalid moves played")

	if bestMove.Word == nil {
		return nil
	}
	return []ScoredMove{bestMove}
}

type WordTree interface {
	IsTerminal() bool
	CanBranch(t core.Tile) (WordTree, bool)
}

var blankA = core.Rune2Letter('a').ToTile(true)
var blankZ = core.Rune2Letter('z').ToTile(true)

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
			s.stepForward(i+dRow, j+dCol, dir, rack, next, prev, callback)
		}
	} else {
		for i, letter := range rack.Rack {
			if !rack.CanConsume(i) {
				continue
			}
			if letter.IsBlank() {
				for r := blankA; r <= blankZ; r++ {
					if next, ok := wordDB.CanBranch(r); ok {
						s.stepForward(i+dRow, j+dCol, dir, rack.Consume(i), next, append(prev, r), callback)
					}
				}
			} else {
				if next, ok := wordDB.CanBranch(letter); ok {
					s.stepForward(i+dRow, j+dCol, dir, rack.Consume(i), next, append(prev, letter), callback)
				}
			}
		}
	}
}

func (s *SmartyAI) stepForward(i, j int, dir core.Direction, rack core.ConsumableRack, wordDB WordTree, prev []core.Tile, callback func([]core.Tile)) {
	// back up perpendicular to advancing direction until I hit a blank
	var (
		ok                 bool
		perpI, perpJ       = i, j
		perpDRow, perpDCol = (!dir).Offsets()
	)
	for s.board.HasTile(perpI, perpJ) {
		perpI -= perpDRow
		perpJ -= perpDCol
	}
	// go back to the last tile I was on
	perpI += perpDRow
	perpJ += perpDRow

	// validate continuous string of tiles is a word
	wordRoot := s.searchSpace
	for s.board.HasTile(perpI, perpJ) {
		t := s.board.Cells[perpI][perpJ].Tile
		if wordRoot, ok = wordRoot.CanBranch(t); !ok {
			// This is not a word, bail out
			return
		}
		perpI += perpDRow
		perpJ += perpDRow
	}

	l := (perpI - i) + (perpJ - j)

	// if so, recurse on search
	if wordRoot.IsTerminal() || l == 0 {
		s.Search(i, j, dir, rack, wordDB, prev, callback)
	}
	// if not, return
}
