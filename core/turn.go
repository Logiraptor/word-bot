package core

import "encoding/json"

type Turn interface {
	isTurn()
}

type ScoredMove struct {
	Move  PlacedWord
	Score Score
}

var _ Turn = ScoredMove{}

func (ScoredMove) isTurn() {}

var _ json.Marshaler = ScoredMove{}

func (p ScoredMove) MarshalJSON() ([]byte, error) {
	type ScoredMoveJSON struct {
		ScoredMove
		Type string `json:"type"`
	}
	return json.Marshal(ScoredMoveJSON{
		ScoredMove: p,
		Type:       "scored_move",
	})
}

type Pass struct{}

var _ Turn = Pass{}

func (Pass) isTurn() {}

var _ json.Marshaler = Pass{}

func (p Pass) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type string `json:"type"`
	}{Type: "pass"})
}

type ChallengeWord struct {
	Move PlacedWord
}

var _ Turn = ChallengeWord{}

func (ChallengeWord) isTurn() {}

var _ json.Marshaler = ChallengeWord{}

func (c ChallengeWord) MarshalJSON() ([]byte, error) {
	type ChallengeMoveJSON struct {
		ChallengeWord
		Type string `json:"type"`
	}
	return json.Marshal(ChallengeMoveJSON{
		ChallengeWord: c,
		Type:          "challenge_word",
	})
}
