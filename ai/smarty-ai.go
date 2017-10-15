package ai

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/Logiraptor/word-bot/core"
	"github.com/Logiraptor/word-bot/wordlist"
)

type SmartyAI struct {
	wordList    core.WordList
	jobs        chan<- job
	searchSpace *wordlist.Trie
}

var _ AI = &SmartyAI{}
var _ MoveGenerator = &SmartyAI{}

var blankA = core.Rune2Letter('a').ToTile(true)
var blankZ = core.Rune2Letter('z').ToTile(true)

func NewSmartyAI(wordList core.WordList, searchSpace *wordlist.Trie) *SmartyAI {

	jobs := make(chan job, 15*15*2)

	s := &SmartyAI{
		wordList:    wordList,
		searchSpace: searchSpace,
		jobs:        jobs,
	}

	for i := 0; i < runtime.NumCPU(); i++ {
		go searchWorker(s, jobs)
	}

	return s
}

type job struct {
	i, j       int
	board      *core.Board
	dir        core.Direction
	rack       core.Rack
	wordDB     *wordlist.Trie
	resultChan chan<- core.PlacedTiles
	wg         *sync.WaitGroup
}

func (s *SmartyAI) FindMove(b *core.Board, bag core.Bag, rack core.Rack, callback func(core.Turn) bool) {
	var bestMove core.ScoredMove
	s.GenerateMoves(b, rack, func(turn core.Turn) bool {
		switch x := turn.(type) {
		case core.ScoredMove:
			if x.Score > bestMove.Score {
				bestMove = x
				return callback(x)
			}
		}
		return true
	})
}

func (s *SmartyAI) GenerateMoves(b *core.Board, rack core.Rack, callback func(core.Turn) bool) {
	var wg = new(sync.WaitGroup)

	dirs := []core.Direction{core.Horizontal, core.Vertical}

	results := make(chan core.PlacedTiles, 10)

	go func() {
		for i := 0; i < 15; i++ {
			for j := 0; j < 15; j++ {
				if b.HasTile(i, j) {
					continue
				}
				for _, dir := range dirs {
					wg.Add(1)
					s.jobs <- job{
						board:      b,
						i:          i,
						j:          j,
						dir:        dir,
						rack:       rack,
						resultChan: results,
						wg:         wg,
						wordDB:     s.searchSpace,
					}
				}
			}
		}

		wg.Wait()
		close(results)
	}()

	for result := range results {
		score := b.Score(result)
		callback(core.ScoredMove{
			PlacedTiles: result,
			Score:       score,
		})
	}
}

func searchWorker(s *SmartyAI, jobs <-chan job) {
	var tiles = make([]core.Tile, 0, 15)
	for job := range jobs {
		s.Search(job.board, job.i, job.j, job.dir, job.rack, job.wordDB, tiles, func(word []core.Tile) {
			if len(word) == 0 {
				return
			}

			result := core.PlacedTiles{
				Word:      word,
				Row:       job.i,
				Col:       job.j,
				Direction: job.dir,
			}
			if job.board.ValidateMove(result, s.wordList) {
				_, canPlay := job.rack.Play(word)
				if canPlay {
					newWord := make([]core.Tile, len(result.Word))
					copy(newWord, result.Word)
					result.Word = newWord

					job.resultChan <- result
				}
			}
		})
		job.wg.Done()
	}
}

func (s *SmartyAI) Kill() {
	close(s.jobs)
}

func (s *SmartyAI) Name() string {
	return "Smarty"
}

func (s *SmartyAI) Search(board *core.Board, i, j int, dir core.Direction, rack core.Rack, wordDB *wordlist.Trie, prev []core.Tile, callback func([]core.Tile)) {
	dRow, dCol := dir.Offsets()
	if wordDB.IsTerminal() {
		fmt.Println("FOUND: ", prev)
		callback(prev)
	}

	if board.OutOfBounds(i, j) {
		fmt.Println("BAIL: Out of bounds")
		return
	}
	if board.HasTile(i, j) {
		letter := board.Cells[i][j].Tile
		if next, ok := wordDB.CanBranch(letter); ok {
			fmt.Println("CONT: Consuming tile from board", letter)
			s.Search(board, i+dRow, j+dCol, dir, rack, next, prev, callback)
		}
	} else {
		for index, letter := range rack.Rack {
			if !rack.CanConsume(index) {
				continue
			}
			if letter.IsBlank() {
				for r := blankA; r <= blankZ; r++ {
					fmt.Println("CONT: Consuming tile from rack", r)
					if next, ok := wordDB.CanBranch(r); ok {
						s.stepForward(board, i, j, r, dir, rack.Consume(index), next, append(prev, r), callback)
					}
				}
			} else {
				if next, ok := wordDB.CanBranch(letter); ok {
					fmt.Println("CONT: Consuming tile from rack", letter)
					s.stepForward(board, i, j, letter, dir, rack.Consume(index), next, append(prev, letter), callback)
				}
			}
		}
	}
}

func (s *SmartyAI) stepForward(board *core.Board, i, j int, placed core.Tile, dir core.Direction, rack core.Rack, wordDB *wordlist.Trie, prev []core.Tile, callback func([]core.Tile)) {
	// back up perpendicular to advancing direction until I hit a blank
	var (
		ok                 bool
		perpDRow, perpDCol = (!dir).Offsets()
		perpI, perpJ       = i - perpDRow, j - perpDCol
	)
	for board.HasTile(perpI, perpJ) {
		perpI -= perpDRow
		perpJ -= perpDCol
	}
	// go back to the last tile I was on
	perpI += perpDRow
	perpJ += perpDCol
	startI, startJ := perpI, perpJ

	// validate continuous string of tiles is a word
	wordRoot := s.searchSpace
	for {
		if board.HasTile(perpI, perpJ) {
			t := board.Cells[perpI][perpJ].Tile
			if wordRoot, ok = wordRoot.CanBranch(t); !ok {
				// This is not a word, bail out
				return
			}
		} else if perpI == i && perpJ == j {
			if wordRoot, ok = wordRoot.CanBranch(placed); !ok {
				return
			}
		} else {
			break
		}
		perpI += perpDRow
		perpJ += perpDCol
	}

	l := (perpI - startI) + (perpJ - startJ)

	// if so, recurse on search
	if wordRoot.IsTerminal() || l == 1 {
		dRow, dCol := dir.Offsets()
		s.Search(board, i+dRow, j+dCol, dir, rack, wordDB, prev, callback)
	}
	// if not, return
}
