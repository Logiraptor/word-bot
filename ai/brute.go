package ai

import (
	"github.com/Logiraptor/word-bot/core"
)

type BruteForceGenerator struct {
	wordDB core.WordList
}

func NewBrute(wordDB core.WordList) BruteForceGenerator {
	return BruteForceGenerator{wordDB}
}

func (b BruteForceGenerator) GenerateMoves(board *core.Board, rack core.Rack, onMove func(core.Turn) bool) {
	shouldEmit := true
	BruteForce(board, rack, b.wordDB, func(t core.Turn) {
		if shouldEmit {
			shouldEmit = onMove(t)
		}
	})
}

func BruteForce(b *core.Board, rack core.Rack, wordDB core.WordList, callback func(core.Turn)) {
	dirs := []core.Direction{core.Horizontal, core.Vertical}
	perms := Permute(rack.Rack)
	for i := 0; i < 15; i++ {
		for j := 0; j < 15; j++ {
			if b.HasTile(i, j) {
				continue
			}

			for _, dir := range dirs {
				for _, p := range perms {

					pt := core.PlacedTiles{Word: p, Row: i, Col: j, Direction: dir}
					if b.ValidateMove(pt, wordDB) {
						_, canPlay := rack.Play(p)
						if canPlay {
							callback(core.ScoredMove{
								PlacedTiles: pt,
								Score:       b.Score(pt),
							})
						}
					}
				}
			}
		}
	}
}

func Permute(str []core.Tile) [][]core.Tile {
	result := [][]core.Tile{}
	if len(str) > 0 {
		for i, c := range str {

			if c.IsBlank() {
				for c := blankA; c <= blankZ; c++ {
					s := append(str[:i:i], str[i+1:]...)
					result = append(result, []core.Tile{c})
					if len(s) > 0 {
						e := Permute(s)
						for j := range e {
							result = append(result, append([]core.Tile{c}, e[j]...))
						}
					}
				}
			} else {
				s := append(str[:i:i], str[i+1:]...)
				result = append(result, []core.Tile{c})
				if len(s) > 0 {
					e := Permute(s)
					for j := range e {
						result = append(result, append([]core.Tile{c}, e[j]...))
					}
				}
			}
		}
	}
	return result
}
