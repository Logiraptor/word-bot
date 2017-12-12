package ai

import "github.com/Logiraptor/word-bot/core"

type MoveChooser struct {
	name      string
	generator MoveGenerator
	evaluator MoveEvaluator
}

func NewMoveChooser(name string, gen MoveGenerator, eval MoveEvaluator) *MoveChooser {
	return &MoveChooser{
		name:      name,
		generator: gen,
		evaluator: eval,
	}
}

var _ AI = &MoveChooser{}

func (m *MoveChooser) FindMove(b *core.Board, bag core.Bag, rack core.Rack, onMove func(core.Turn) bool) {
	var bestMove core.ScoredMove
	var bestScore float64
	m.generator.GenerateMoves(b, rack, func(t core.Turn) bool {
		if sm, ok := t.(core.ScoredMove); ok {
			score := m.evaluator.Evaluate(b, rack, sm)
			if score > bestScore {
				bestScore = score
				bestMove = sm
				return onMove(sm)
			}
		}
		if bestScore == 0 {
			return onMove(t)
		}
		return true
	})
}

func (m *MoveChooser) Name() string {
	return m.name
}
