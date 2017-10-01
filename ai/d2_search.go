package ai

import "github.com/Logiraptor/word-bot/core"
import "math"

type MoveGenerator interface {
	GenerateMoves(board *core.Board, rack LetterSet, numMoves int) []core.ScoredMove
}

type D2Searcher struct {
	moveGenerator MoveGenerator
}

func NewD2Searcher(moveGenerator MoveGenerator) *D2Searcher {
	return &D2Searcher{
		moveGenerator: moveGenerator,
	}
}

func (d *D2Searcher) FindMove(board *core.Board, rack LetterSet, bag LetterSet) core.ScoredMove {
	bestMoves := d.moveGenerator.GenerateMoves(board, rack, 10)

	var bestMove core.ScoredMove
	var bestDiff = math.Inf(-1)

	for _, move := range bestMoves {
		localBoard := board.Copy()
		bestOpponentMoves := d.moveGenerator.GenerateMoves(localBoard, bag, 10)
		avgScore := 0.0
		for _, opponentMove := range bestOpponentMoves {
			avgScore += float64(opponentMove.Score) / 10
		}

		diff := float64(move.Score) - avgScore
		if diff > bestDiff {
			bestMove = move
			bestDiff = diff
		}
	}

	return bestMove
}
