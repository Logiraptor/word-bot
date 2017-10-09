package ai

import (
	"github.com/Logiraptor/word-bot/core"
)

type Playout struct {
	ai AI
}

func NewPlayout(ai AI) *Playout {
	return &Playout{
		ai: ai,
	}
}

var _ BoardEvaluator = &Playout{}

func (p *Playout) Evaluate(b *core.Board, bag core.Bag, rack core.Rack) float64 {
	b = b.Clone()

	opponentRack := core.NewConsumableRack(nil)
	bag, opponentRack.Rack = bag.FillRack(opponentRack.Rack, 7)

	var (
		p1Ok, p2Ok       bool = true, true
		p1Score, p2Score core.Score
		pt               core.ScoredMove
	)

	for (bag.Count() > 0 || len(opponentRack.Rack) > 0 || len(rack.Rack) > 0) && (p1Ok || p2Ok) {
		p1Ok, p2Ok = false, false
		var move core.Turn
		p.ai.FindMove(b, bag, opponentRack, func(turn core.Turn) bool {
			move = turn
			return true
		})

		pt, p2Ok = move.(core.ScoredMove)
		if p2Ok {
			if opponentRack, p2Ok = opponentRack.Play(pt.Word); p2Ok {
				b.PlaceTiles(pt.PlacedTiles)
				bag, opponentRack.Rack = bag.FillRack(opponentRack.Rack, 7-len(opponentRack.Rack))
				p2Score += pt.Score
			}
		}

		move = nil
		p.ai.FindMove(b, bag, rack, func(turn core.Turn) bool {
			move = turn
			return true
		})

		pt, p1Ok = move.(core.ScoredMove)
		if p1Ok {
			if rack, p1Ok = rack.Play(pt.Word); p1Ok {
				b.PlaceTiles(pt.PlacedTiles)
				bag, rack.Rack = bag.FillRack(rack.Rack, 7-len(rack.Rack))
				p1Score += pt.Score
			}
		}
	}

	return float64(p1Score) - float64(p2Score)
}
