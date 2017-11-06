package ai

import (
	"runtime"
	"sync"

	"github.com/Logiraptor/word-bot/wordlist"

	"github.com/Logiraptor/word-bot/core"
)

type SpeedyAI struct {
	wordList    core.WordList
	jobs        chan<- speedyJob
	searchSpace *wordlist.Gaddag
}

var _ AI = &SpeedyAI{}
var _ MoveGenerator = &SpeedyAI{}

func NewSpeedyAI(wordList core.WordList, searchSpace *wordlist.Gaddag) *SpeedyAI {

	jobs := make(chan speedyJob, 15*15*2)

	s := &SpeedyAI{
		wordList:    wordList,
		searchSpace: searchSpace,
		jobs:        jobs,
	}

	for i := 0; i < runtime.NumCPU(); i++ {
		go speedySearchWorker(s, jobs)
	}

	return s
}

type speedyJob struct {
	i, j             int
	board            *core.Board
	boardConstraints [15][15]Constraint
	dir              core.Direction
	rack             core.Rack
	wordDB           *wordlist.Gaddag
	resultChan       chan<- core.PlacedTiles
	wg               *sync.WaitGroup
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

	var boardConstraint [15][15]Constraint
	for i := 0; i < 15; i++ {
		for j := 0; j < 15; j++ {
			boardConstraint[i][j] = PermittedTiles(b, s.searchSpace, i, j)
		}
	}

	go func() {
		for i := 0; i < 15; i++ {
			for j := 0; j < 15; j++ {
				for _, dir := range dirs {

					if !b.HasTile(i, j) && ((i == 7 && j == 7) ||
						b.HasTile(i-1, j) || b.HasTile(i, j-1) ||
						b.HasTile(i+1, j) || b.HasTile(i, j+1)) {
						wg.Add(1)
						s.jobs <- speedyJob{
							board:            b,
							boardConstraints: boardConstraint,
							i:                i,
							j:                j,
							dir:              dir,
							rack:             rack,
							resultChan:       results,
							wg:               wg,
							wordDB:           s.searchSpace,
						}
						// fmt.Println("Pursuing", i, j)
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

func speedySearchWorker(s *SpeedyAI, jobs <-chan speedyJob) {
	var tiles = make([]core.Tile, 0, 15)
	for job := range jobs {
		s.Search(job.board, job.boardConstraints, job.i, job.j, job.dir, job.rack, job.wordDB, tiles, func(i, j int, reversePrefix, rest []core.Tile) {
			// fmt.Println("RECEIVED:", reversePrefix, rest)
			if len(reversePrefix)+len(rest) == 0 {
				return
			}

			word := make([]core.Tile, len(reversePrefix)+len(rest))
			p := 0
			for i := len(reversePrefix) - 1; i >= 0; i-- {
				word[p] = reversePrefix[i]
				p++
			}
			for _, x := range rest {
				word[p] = x
				p++
			}

			result := core.PlacedTiles{
				Word:      word,
				Row:       i,
				Col:       j,
				Direction: job.dir,
			}
			// fmt.Println("RECONSTRUCTED", result)
			if job.board.ValidateMove(result, s.wordList) {
				_, canPlay := job.rack.Play(word)
				if canPlay {
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

func (s *SpeedyAI) Search(board *core.Board, boardConstraint [15][15]Constraint, i, j int, dir core.Direction, rack core.Rack, wordDB *wordlist.Gaddag, prev []core.Tile, callback func(int, int, []core.Tile, []core.Tile)) {
	// fmt.Println("CONT: Starting search at ", i, j, dir)
	// fmt.Println("CONT: With words", wordDB.DumpOptions())
	s.searchForward(board, boardConstraint, i, j, i, j, dir, rack, wordDB, prev, callback)
}

func (s *SpeedyAI) searchForward(board *core.Board, boardConstraint [15][15]Constraint, startI, startJ, row, col int, dir core.Direction, rack core.Rack, wordDB *wordlist.Gaddag, prev []core.Tile, callback func(int, int, []core.Tile, []core.Tile)) {
	// fmt.Println("CONT: searchForward", row, col, prev)
	// fmt.Println("CONT: With words", wordDB.DumpOptions())

	dRow, dCol := dir.Offsets()
	if board.HasTile(row, col) {
		// fmt.Println("CONT: Attempting to consume board tile")
		letter := board.Cells[row][col].Tile
		if !wordDB.CanBranch(letter) {
			// fmt.Println("BAIL: cannot branch on board tile")
			return
		}
		// fmt.Println("CONT: consuming board tile", letter)
		s.searchForward(board, boardConstraint, startI, startJ, row+dRow, col+dCol, dir, rack, wordDB.Branch(letter), prev, callback)
		return
	}

	if wordDB.CanReverse() {
		// fmt.Println("CONT: can reverse, walking backwards from ", startI-dRow, startJ-dCol)
		s.searchBackward(board, boardConstraint, startI-dRow, startJ-dCol, dir, rack, wordDB.Reverse(), nil, prev, callback)
	}

	if board.OutOfBounds(row, col) {
		// fmt.Println("BAIL: out of bounds")
		return
	}

	for i, letter := range rack.Rack {
		if !rack.CanConsume(i) {
			continue
		}
		if letter.IsBlank() {
			// fmt.Println("CONT: attempting to consume blank from rack")
			for r := blankA; r <= blankZ; r++ {
				// fmt.Println("CONT: assigning blank as", r)
				if !wordDB.CanBranch(r) {
					// fmt.Println("BAIL: Cannot branch on", r)
					continue
				}
				if !boardConstraint[row][col].AllowsTile(r) {
					continue
				}
				s.searchForward(board, boardConstraint, startI, startJ, row+dRow, col+dCol, dir, rack.Consume(i), wordDB.Branch(r), append(prev, r), callback)
			}
			continue
		}

		// fmt.Println("CONT: consuming rack tile", letter)

		if !wordDB.CanBranch(letter) {
			// fmt.Println("BAIL: cannot branch on rack tile", letter)
			continue
		}
		if !boardConstraint[row][col].AllowsTile(letter) {
			continue
		}

		s.searchForward(board, boardConstraint, startI, startJ, row+dRow, col+dCol, dir, rack.Consume(i), wordDB.Branch(letter), append(prev, letter), callback)
	}
}

func (s *SpeedyAI) searchBackward(board *core.Board, boardConstraint [15][15]Constraint, row, col int, dir core.Direction, rack core.Rack, wordDB *wordlist.Gaddag, prefix []core.Tile, rest []core.Tile, callback func(int, int, []core.Tile, []core.Tile)) {
	// fmt.Println("CONT: searchBackward", row, col, rest, "#", prefix)
	// fmt.Println("CONT: With words", wordDB.DumpOptions(), rack)
	dRow, dCol := dir.Offsets()
	dRow *= -1
	dCol *= -1

	if board.HasTile(row, col) {
		// fmt.Println("BACK: attempting to consume board tile")
		letter := board.Cells[row][col].Tile
		if !wordDB.CanBranch(letter) {
			// fmt.Println("BAIL: cannot branch on board tile", letter)
			return
		}
		// fmt.Println("CONT: consuming board tile", letter)
		s.searchBackward(board, boardConstraint, row+dRow, col+dCol, dir, rack, wordDB.Branch(letter), prefix, rest, callback)
		return
	}
	if wordDB.IsTerminal() {
		// fmt.Println("TERM: ", row, col, prefix, rest)
		callback(row-dRow, col-dCol, prefix, rest)
	}
	if board.OutOfBounds(row, col) {
		// fmt.Println("BAIL: out of bounds")
		return
	}

	for i, letter := range rack.Rack {
		// fmt.Println("CONT: attempting to consume rack tile", letter, rack, i, len(rack.Rack))
		if !rack.CanConsume(i) {
			continue
		}
		if letter.IsBlank() {
			// fmt.Println("CONT: attempting to consume blank from rack")
			for r := blankA; r <= blankZ; r++ {
				// fmt.Println("CONT: attempting to assign blank to", r)
				if !wordDB.CanBranch(r) {
					// fmt.Println("BAIL: cannot branch on ", r)
					continue
				}
				if !boardConstraint[row][col].AllowsTile(r) {
					continue
				}
				// fmt.Println("CONT: assigning blank to ", r)
				s.searchBackward(board, boardConstraint, row+dRow, col+dCol, dir, rack.Consume(i), wordDB.Branch(r), append(prefix, r), rest, callback)
			}
			continue
		}

		if !wordDB.CanBranch(letter) {
			// fmt.Println("BAIL: cannot branch on ", letter)
			continue
		}
		if !boardConstraint[row][col].AllowsTile(letter) {
			continue
		}

		// fmt.Println("CONT: consuming rack tile ", letter)
		s.searchBackward(board, boardConstraint, row+dRow, col+dCol, dir, rack.Consume(i), wordDB.Branch(letter), append(prefix, letter), rest, callback)
	}
}

// for each square on the board:
//   if square is an anchor:
//     if terminal, trigger callback
//     if oob, bail
//     for each unconsumed tile in rack:
//     attempt to play
//       gaddag can branch
//       board cross-word constraints allow the tile
//       recurse on each unconsumed tile at position++
