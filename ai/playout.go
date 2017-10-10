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

func (p *Playout) Evaluate(b *core.Board, bag core.Bag, p1, p2 core.Rack) float64 {
	b = b.Clone()

	bag, p2.Rack = bag.FillRack(p2.Rack, 7-len(p2.Rack))
	bag, p1.Rack = bag.FillRack(p1.Rack, 7-len(p1.Rack))

	var (
		p1Ok, p2Ok       bool = true, true
		p1Score, p2Score core.Score
		pt               core.ScoredMove
	)

	for (bag.Count() > 0 || len(p2.Rack) > 0 || len(p1.Rack) > 0) && (p1Ok || p2Ok) {
		p1Ok, p2Ok = false, false
		var move core.Turn
		p.ai.FindMove(b, bag, p2, func(turn core.Turn) bool {
			move = turn
			return true
		})

		pt, p2Ok = move.(core.ScoredMove)
		if p2Ok {
			if p2, p2Ok = p2.Play(pt.Word); p2Ok {
				b.PlaceTiles(pt.PlacedTiles)
				bag, p2.Rack = bag.FillRack(p2.Rack, 7-len(p2.Rack))
				p2Score += pt.Score
			}
		}

		move = nil
		p.ai.FindMove(b, bag, p1, func(turn core.Turn) bool {
			move = turn
			return true
		})

		pt, p1Ok = move.(core.ScoredMove)
		if p1Ok {
			if p1, p1Ok = p1.Play(pt.Word); p1Ok {
				b.PlaceTiles(pt.PlacedTiles)
				bag, p1.Rack = bag.FillRack(p1.Rack, 7-len(p1.Rack))
				p1Score += pt.Score
			}
		}
	}

	return float64(p1Score) - float64(p2Score)
}
