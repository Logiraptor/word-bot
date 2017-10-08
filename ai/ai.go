package ai

import "github.com/Logiraptor/word-bot/core"

// AI can automate gameplay
type AI interface {
	// onMove will be called with increasingly valuable moves until it returns false or
	// all moves have been generated
	FindMove(b *core.Board, rack core.Rack, onMove func(core.Turn) bool)
	Name() string
}

// MoveGenerator will generate moves, calling onMove until it returns false
// Or all moves have been generated.
type MoveGenerator interface {
	GenerateMoves(b *core.Board, rack core.Rack, onMove func(core.Turn) bool)
}

// A WordTree enables efficient, progressive building of turns
type WordTree interface {
	IsTerminal() bool
	CanBranch(t core.Tile) (WordTree, bool)
}
