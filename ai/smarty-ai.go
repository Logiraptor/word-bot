package ai

import (
	"runtime"
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

	jobs := make(chan Job, 15*15*2)

	s := &SmartyAI{
		board:       board,
		wordList:    wordList,
		searchSpace: searchSpace,
		jobs:        jobs,
	}

	for i := 0; i < runtime.NumCPU(); i++ {
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
	wg         *sync.WaitGroup
}

type Result struct {
	word     []core.Tile
	row, col int
	dir      core.Direction
}

func searchWorker(s *SmartyAI, jobs <-chan Job) {
	var tiles = make([]core.Tile, 0, 15)
	for job := range jobs {
		s.Search(job.i, job.j, job.dir, job.rack, job.wordDB, tiles, func(word []core.Tile) {
			job.resultChan <- Result{
				word: word,
				row:  job.i,
				col:  job.j,
				dir:  job.dir,
			}
		})
		job.wg.Done()
	}
}

func (b *SmartyAI) FindMoves(tiles []core.Tile) []ScoredMove {
	var bestMove ScoredMove
	var wg = new(sync.WaitGroup)

	var badMoves uint64

	rack := core.NewConsumableRack(tiles)

	dirs := []core.Direction{core.Horizontal, core.Vertical}

	results := make(chan Result, 10)

	go func() {
		for i := 0; i < 15; i++ {
			for j := 0; j < 15; j++ {
				for _, dir := range dirs {
					wg.Add(1)
					b.jobs <- Job{
						i:          i,
						j:          j,
						dir:        dir,
						rack:       rack,
						resultChan: results,
						wg:         wg,
						wordDB:     b.searchSpace,
					}
				}
			}
		}

		wg.Wait()
		close(results)
	}()

	for result := range results {
		if len(result.word) == 0 {
			continue
		}

		if b.board.ValidateMove(result.word, result.row, result.col, result.dir, b.wordList) {
			score := b.board.Score(result.word, result.row, result.col, result.dir)

			if score > bestMove.Score {
				newWord := make([]core.Tile, len(result.word))
				copy(newWord, result.word)

				current := ScoredMove{
					PlacedWord: core.PlacedWord{
						Word:      newWord,
						Row:       result.row,
						Col:       result.col,
						Direction: result.dir,
					},
					Score: score,
				}

				bestMove = current
			}
		} else {
			atomic.AddUint64(&badMoves, 1)
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
