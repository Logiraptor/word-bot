package ai

import "github.com/Logiraptor/word-bot/core"
import "github.com/Logiraptor/word-bot/persist"

type LeaveWeighter struct {
	weights map[string]float64
}

func NewLeaveWeighter(db *persist.DB) *LeaveWeighter {
	records, _ := db.LoadLeaveWeights()
	weights := make(map[string]float64)
	for _, r := range records {
		weights[r.Leave] = r.Weight
	}
	return &LeaveWeighter{
		weights: weights,
	}
}

var _ MoveEvaluator = &LeaveWeighter{}

func (l *LeaveWeighter) Evaluate(b *core.Board, rack core.Rack, move core.ScoredMove) float64 {
	leave, _ := rack.Play(move.Word)
	if leaveScore, ok := l.weights[core.Tiles2String(leave.Rack)]; ok {
		return float64(move.Score) + leaveScore*10
	}
	return float64(move.Score) + 11
}
