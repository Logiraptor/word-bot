package ai

import (
	"fmt"

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
		p1Skip, p2Skip   bool
		p1Score, p2Score core.Score
	)

	for (bag.Count() > 0 || len(opponentRack.Rack) > 0 || len(rack.Rack) > 0) || (p1Skip && p2Skip) {
		var move core.Turn
		p.ai.FindMove(b, bag, opponentRack, func(turn core.Turn) bool {
			move = turn
			return true
		})

		if pt, ok := move.(core.ScoredMove); ok {
			p2Skip = false
			if opponentRack, ok = opponentRack.Play(pt.Word); ok {
				fmt.Println("Play", pt.PlacedTiles)
				b.PlaceTiles(pt.PlacedTiles)
				bag, opponentRack.Rack = bag.FillRack(opponentRack.Rack, 7-len(opponentRack.Rack))
				p2Score += pt.Score
			}
		} else {
			p2Skip = true
		}

		p.ai.FindMove(b, bag, rack, func(turn core.Turn) bool {
			move = turn
			return true
		})

		if pt, ok := move.(core.ScoredMove); ok {
			p1Skip = false
			if rack, ok = rack.Play(pt.Word); ok {
				fmt.Println("Counter", pt.PlacedTiles)
				b.PlaceTiles(pt.PlacedTiles)
				bag, rack.Rack = bag.FillRack(rack.Rack, 7-len(rack.Rack))
				p1Score += pt.Score
			}
		} else {
			p1Skip = true
		}
	}

	return float64(p1Score) - float64(p2Score)
}
