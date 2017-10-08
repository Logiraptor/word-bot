package smarter

import (
	"fmt"

	"github.com/Logiraptor/word-bot/ai"
	"github.com/Logiraptor/word-bot/core"

	"github.com/glemzurg/go-mcts"
)

type GameState struct {
	wordList    core.WordList
	searchSpace ai.WordTree
	moves       []core.ScoredMove
	rack        core.Rack
}

var _ mcts.GameState = &GameState{}

func (g *GameState) AvailableMoves() []mcts.Move {
	b := core.NewBoard()
	for _, m := range g.moves {
		b.PlaceTiles(m.PlacedTiles)
	}
	s := ai.NewSmartyAI(g.wordList, g.searchSpace)
	defer s.Kill()
	// Run smarty
	moves := []mcts.Move{}
	s.GenerateMoves(b, g.rack, func(turn core.Turn) bool {
		moves = append(moves, Move{
			Turn: turn,
		})
		return true
	})
	// return moves from smarty
	return moves
}

func (g *GameState) Clone() mcts.GameState {
	newMoves := make([]core.ScoredMove, len(g.moves))
	copy(newMoves, g.moves)
	return &GameState{
		moves:       newMoves,
		rack:        g.rack,
		wordList:    g.wordList,
		searchSpace: g.searchSpace,
	}
}

func (g *GameState) MakeMove(m mcts.Move) {
	switch v := m.(Move).Turn.(type) {
	case core.ScoredMove:
		g.moves = append(g.moves, v)
	default:
		panic(fmt.Sprintf("Smarter cannot make move %#v", v))
	}
}

func (g *GameState) RandomizeUnknowns() {
	// Randomize bag?
}

type Move struct {
	core.Turn
}

var _ mcts.Move = Move{}

func (m Move) Probability() float64 {
	// this is weird and always super low, so maybe keeping a small constant here will work...
	return 0.001
}

type MCTSAI struct {
}

func (m MCTSAI) FindMove(moves []core.ScoredMove, rack core.Rack) core.Turn {
	move := mcts.Uct(&GameState{moves: moves, rack: rack}, 10, 10, 0.1, 0, m.Score)
	return move.(Move).Turn
}

// Score the game from playerIds perspective
func (m MCTSAI) Score(playerId uint64, s mcts.GameState) float64 {
	// Simulate with smarty playout and return difference in score

	return 1
}
