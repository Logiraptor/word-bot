package ai

import (
	"github.com/Logiraptor/word-bot/core"
)

func BruteForce(b *core.Board, rack core.Rack, wordDB core.WordList, callback func(core.Turn)) {
	dirs := []core.Direction{core.Horizontal, core.Vertical}
	perms := permute(rack.Rack)
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

func permute(rack []core.Tile) [][]core.Tile {
	if len(rack) == 0 {
		return [][]core.Tile{nil}
	}
	first := rack[0]
	rest := rack[1:]
	subPerm := permute(rest)
	output := make([][]core.Tile, len(subPerm), len(subPerm)*2)
	copy(output, subPerm)

	if first.IsBlank() {
		for option := 'a'; option <= 'z'; option++ {
			letter := core.Rune2Letter(option).ToTile(true)
			for _, perm := range subPerm {
				for i := range perm {
					output = append(output, append(append(perm[:i:i], letter), perm[i:]...))
				}
				output = append(output, append(perm, letter))
			}
		}
		return output
	}
	for _, perm := range subPerm {
		for i := range perm {
			output = append(output, append(append(perm[:i:i], first), perm[i:]...))
		}
		output = append(output, append(perm, first))
	}
	return output
}
