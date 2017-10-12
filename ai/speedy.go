package ai

import (
	"sync"

	"github.com/Logiraptor/word-bot/wordlist"

	"github.com/Logiraptor/word-bot/core"
)

type SpeedyAI struct {
	wordList    core.WordList
	jobs        chan<- job
	searchSpace *wordlist.Trie
}

var _ AI = &SpeedyAI{}
var _ MoveGenerator = &SpeedyAI{}

func NewSpeedyAI(wordList core.WordList, searchSpace *wordlist.Trie) *SpeedyAI {

	jobs := make(chan job, 15*15*2)

	s := &SpeedyAI{
		wordList:    wordList,
		searchSpace: searchSpace,
		jobs:        jobs,
	}

	for i := 0; i < 1; i++ {
		go speedySearchWorker(s, jobs)
	}

	return s
}

func (s *SpeedyAI) FindMove(b *core.Board, bag core.Bag, rack core.Rack, callback func(core.Turn) bool) {
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

func (s *SpeedyAI) GenerateMoves(b *core.Board, rack core.Rack, callback func(core.Turn) bool) {
	var wg = new(sync.WaitGroup)

	dirs := []core.Direction{core.Horizontal, core.Vertical}

	results := make(chan core.PlacedTiles, 10)

	go func() {
		for i := 0; i < 15; i++ {
			for j := 0; j < 15; j++ {
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

func speedySearchWorker(s *SpeedyAI, jobs <-chan job) {
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

func (s *SpeedyAI) Kill() {
	close(s.jobs)
}

func (s *SpeedyAI) Name() string {
	return "Speedy"
}

func (s *SpeedyAI) Search(board *core.Board, i, j int, dir core.Direction, rack core.Rack, wordDB *wordlist.Trie, prev []core.Tile, callback func([]core.Tile)) {
	dRow, dCol := dir.Offsets()
	if wordDB.IsTerminal() {
		callback(prev)
	}

	if board.OutOfBounds(i, j) {
		return
	}
	if board.HasTile(i, j) {
		letter := board.Cells[i][j].Tile
		if next, ok := wordDB.CanBranch(letter); ok {
			s.stepForward(board, i+dRow, j+dCol, dir, rack, next, prev, callback)
		}
	} else {
		for i, letter := range rack.Rack {
			if !rack.CanConsume(i) {
				continue
			}
			if letter.IsBlank() {
				for r := blankA; r <= blankZ; r++ {
					if next, ok := wordDB.CanBranch(r); ok {
						s.stepForward(board, i+dRow, j+dCol, dir, rack.Consume(i), next, append(prev, r), callback)
					}
				}
			} else {
				if next, ok := wordDB.CanBranch(letter); ok {
					s.stepForward(board, i+dRow, j+dCol, dir, rack.Consume(i), next, append(prev, letter), callback)
				}
			}
		}
	}
}

func (s *SpeedyAI) stepForward(board *core.Board, i, j int, dir core.Direction, rack core.Rack, wordDB *wordlist.Trie, prev []core.Tile, callback func([]core.Tile)) {
	// back up perpendicular to advancing direction until I hit a blank
	var (
		ok                 bool
		perpI, perpJ       = i, j
		perpDRow, perpDCol = (!dir).Offsets()
	)
	for board.HasTile(perpI, perpJ) {
		perpI -= perpDRow
		perpJ -= perpDCol
	}
	// go back to the last tile I was on
	perpI += perpDRow
	perpJ += perpDRow

	// validate continuous string of tiles is a word
	wordRoot := s.searchSpace
	for board.HasTile(perpI, perpJ) {
		t := board.Cells[perpI][perpJ].Tile
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
		s.Search(board, i, j, dir, rack, wordDB, prev, callback)
	}
	// if not, return
}
