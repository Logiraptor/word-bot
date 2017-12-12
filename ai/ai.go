package ai

import "github.com/Logiraptor/word-bot/core"

// AI can automate gameplay
type AI interface {
	// onMove will be called with increasingly valuable moves until it returns false or
	// all moves have been generated
	FindMove(b *core.Board, bag core.Bag, rack core.Rack, onMove func(core.Turn) bool)
	Name() string
}

// MoveGenerator will generate moves, calling onMove until it returns false
// Or all moves have been generated.
type MoveGenerator interface {
	GenerateMoves(b *core.Board, rack core.Rack, onMove func(core.Turn) bool)
}

// A BoardEvaluator determines a heuristic score for a board position
type BoardEvaluator interface {
	Evaluate(b *core.Board, bag core.Bag, p1, p2 core.Rack) float64
}

type MoveEvaluator interface {
	Evaluate(b *core.Board, rack core.Rack, move core.ScoredMove) float64
}
