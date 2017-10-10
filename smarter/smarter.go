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

	moves := []mcts.Move{Move{Turn: core.Pass{}}, Move{Turn: core.Exchange{}}}
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
			return true
		}
		if _, ok := x.(core.Exchange); ok {
			return true
		}
		if _, ok := y.(core.Pass); ok {
			return true
		}
		if _, ok := y.(core.Exchange); ok {
			return true
		}
		smX := x.(core.ScoredMove)
		smY := y.(core.ScoredMove)
		return smX.Score > smY.Score
	})
	// return moves from smarty
	if len(moves) > 10 {
		moves = moves[:10]
	}
	return moves
}

func (g *GameState) Clone() mcts.GameState {
	return &GameState{
		opponentTurn: g.opponentTurn,
		moveGen:      g.moveGen,
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
				g.lastPlay = v
				g.board.PlaceTiles(v.PlacedTiles)
				g.rack = rack
			} else {
				fmt.Println("Cannot play the tiles")
			}
		} else {
			if rack, ok := g.opponentRack.Play(v.Word); ok {
				g.lastPlay = v
				g.board.PlaceTiles(v.PlacedTiles)
				g.opponentRack = rack
			} else {
				fmt.Println("Cannot play the tiles")
			}
		}

		return
	case core.Exchange:
		newRack := core.NewConsumableRack(nil)
		g.bag, newRack.Rack = g.bag.FillRack(newRack.Rack, 7-len(newRack.Rack))
		if g.opponentTurn {
			g.bag = g.bag.Replace(g.opponentRack.Rack)
			g.opponentRack = newRack
			g.bag, g.opponentRack.Rack = g.bag.FillRack(g.opponentRack.Rack, 7-len(g.opponentRack.Rack))
		} else {
			g.bag = g.bag.Replace(g.rack.Rack)
			g.rack = newRack
			g.bag, g.rack.Rack = g.bag.FillRack(g.rack.Rack, 7-len(g.rack.Rack))
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
	if g.opponentTurn {
		g.bag, g.rack.Rack = g.bag.FillRack(g.rack.Rack, 7-len(g.rack.Rack))
		g.bag, g.opponentRack.Rack = g.bag.FillRack(g.opponentRack.Rack, 7-len(g.opponentRack.Rack))
	} else {
		g.bag, g.opponentRack.Rack = g.bag.FillRack(g.opponentRack.Rack, 7-len(g.opponentRack.Rack))
		g.bag, g.rack.Rack = g.bag.FillRack(g.rack.Rack, 7-len(g.rack.Rack))
	}
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
	moveGen                 ai.MoveGenerator
	eval                    ai.BoardEvaluator
	iterations, simulations uint
	bias                    float64
}

func NewMCTSAI(moveGen ai.MoveGenerator, eval ai.BoardEvaluator, iterations, simulations uint, bias float64) *MCTSAI {
	return &MCTSAI{
		moveGen:     moveGen,
		eval:        eval,
		iterations:  iterations,
		simulations: simulations,
		bias:        bias,
	}
}

var _ ai.AI = &MCTSAI{}

func (m *MCTSAI) FindMove(board *core.Board, bag core.Bag, rack core.Rack, callback func(core.Turn) bool) {

	r := core.NewConsumableRack(nil)

	move := mcts.Uct(&GameState{
		board:        board,
		rack:         rack,
		moveGen:      m.moveGen,
		opponentRack: r,
		bag:          bag,
	}, m.iterations, m.simulations, m.bias, 0, m.Score)
	callback(move.(Move).Turn)
}

func (m *MCTSAI) Name() string {
	return fmt.Sprintf("monty %d %d %f", m.iterations, m.simulations, m.bias)
}

// Score the game from playerIds perspective
func (m *MCTSAI) Score(playerId uint64, s mcts.GameState) float64 {
	gs := s.(*GameState)
	score := 0.0
	if gs.opponentTurn {
		score = -m.eval.Evaluate(gs.board, gs.bag, gs.opponentRack, gs.rack)
	} else {
		score = m.eval.Evaluate(gs.board, gs.bag, gs.rack, gs.opponentRack)
	}
	return score
}
