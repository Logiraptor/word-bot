package smarter

import (
	"fmt"
	"sort"

	"github.com/Logiraptor/word-bot/ai"
	"github.com/Logiraptor/word-bot/core"
	mcts "github.com/glemzurg/go-mcts"
)

type GameState struct {
	opponentTurn bool
	moveGen      ai.MoveGenerator
	eval         ai.BoardEvaluator
	lastPlay     core.ScoredMove
	board        *core.Board
	rack         core.Rack
	opponentRack core.Rack
	bag          core.Bag
}

func (g *GameState) String() string {
	rack := g.rack
	if g.opponentTurn {
		rack = g.opponentRack
	}
	return fmt.Sprintf("%t => %s => %s", g.opponentTurn, g.lastPlay, core.Tiles2String(rack.Rack))
}

var _ mcts.GameState = &GameState{}

func (g *GameState) AvailableMoves() []mcts.Move {
	// Run smarty
	var rack = g.rack
	if g.opponentTurn {
		rack = g.opponentRack
	}

	moves := []mcts.Move{Move{Turn: core.Pass{}}}
	g.moveGen.GenerateMoves(g.board, rack, func(turn core.Turn) bool {
		moves = append(moves, Move{
			Turn: turn,
		})
		return true
	})

	sort.Slice(moves, func(i, j int) bool {
		x := moves[i].(Move).Turn
		y := moves[j].(Move).Turn
		if _, ok := x.(core.Pass); ok {
			return false
		}
		if _, ok := y.(core.Pass); ok {
			return false
		}
		smX := x.(core.ScoredMove)
		smY := y.(core.ScoredMove)
		return smX.Score > smY.Score
	})
	// return moves from smarty
	if len(moves) > 4 {
		moves = moves[:4]
	}
	return moves
}

func (g *GameState) Clone() mcts.GameState {
	return &GameState{
		opponentTurn: g.opponentTurn,
		moveGen:      g.moveGen,
		eval:         g.eval,
		board:        g.board.Clone(),
		rack:         g.rack,
		opponentRack: g.opponentRack,
		bag:          g.bag,
	}
}

func (g *GameState) MakeMove(m mcts.Move) {
	g.opponentTurn = !g.opponentTurn
	switch v := m.(Move).Turn.(type) {
	case core.ScoredMove:
		if g.opponentTurn {
			if rack, ok := g.rack.Play(v.Word); ok {
				// fmt.Printf("PLAY %s FROM %s\n", core.Tiles2String(v.Word), core.Tiles2String(g.rack.Rack))

				// fmt.Printf("REMOVE %s GIVES %s\n", core.Tiles2String(v.Word), core.Tiles2String(rack.Rack))

				g.lastPlay = v
				g.board.PlaceTiles(v.PlacedTiles)
				g.bag, rack.Rack = g.bag.FillRack(rack.Rack, 7-len(rack.Rack))
				g.rack = rack

				// fmt.Printf("FILL FROM BAG WITH %d TILES GIVES %s\n", g.bag.Count(), core.Tiles2String(g.rack.Rack))
			} else {
				fmt.Println("Cannot play the tiles")
			}
		} else {
			if rack, ok := g.opponentRack.Play(v.Word); ok {
				g.lastPlay = v
				g.board.PlaceTiles(v.PlacedTiles)
				g.bag, rack.Rack = g.bag.FillRack(rack.Rack, 7-len(rack.Rack))
				g.opponentRack = rack
			} else {
				fmt.Println("Cannot play the tiles")
			}
		}

		return
	case core.Pass:
		return
	default:
		panic(fmt.Sprintf("Smarter cannot make move %#v", v))
	}
}

func (g *GameState) RandomizeUnknowns() {
	g.bag = g.bag.Shuffle()
}

type Move struct {
	core.Turn
}

var _ mcts.Move = Move{}

func (m Move) Probability() float64 {
	// this is weird and always super low, so maybe keeping a small constant here will work...
	return 1
}

type MCTSAI struct {
	moveGen ai.MoveGenerator
	eval    ai.BoardEvaluator
}

func NewMCTSAI(moveGen ai.MoveGenerator, eval ai.BoardEvaluator) *MCTSAI {
	return &MCTSAI{
		moveGen: moveGen,
		eval:    eval,
	}
}

var _ ai.AI = &MCTSAI{}

func (m *MCTSAI) FindMove(board *core.Board, bag core.Bag, rack core.Rack, callback func(core.Turn) bool) {

	r := core.NewConsumableRack(nil)
	bag, r.Rack = bag.FillRack(r.Rack, 7)

	move := mcts.Uct(&GameState{
		board:        board,
		rack:         rack,
		moveGen:      m.moveGen,
		opponentRack: r,
		bag:          bag,
	}, 100, 10, 10, 0, m.Score)
	callback(move.(Move).Turn)
}

func (m *MCTSAI) Name() string {
	return "monte 11"
}

// Score the game from playerIds perspective
func (m *MCTSAI) Score(playerId uint64, s mcts.GameState) float64 {
	// Simulate with smarty playout and return difference in score
	gs := s.(*GameState)
	if gs.opponentTurn {
		return -m.eval.Evaluate(gs.board, gs.bag, gs.opponentRack)
	}
	return m.eval.Evaluate(gs.board, gs.bag, gs.rack)
}
