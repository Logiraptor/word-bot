package web

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderWithTags(t *testing.T) {
	board := Render(MoveRequest{
		Moves: []Move{
			{
				Col: 7, Row: 7, Dir: "horizontal", Tiles: []TileJS{
					{
						Blank:  true,
						Bonus:  "",
						Flags:  []uint{0},
						Letter: "",
						Value:  0,
					},
					{
						Blank:  true,
						Bonus:  "",
						Flags:  []uint{1},
						Letter: "",
						Value:  0,
					},
				},
			},
		},
		Rack: []TileJS{},
	})

	assert.EqualValues(t, board.Board[7][7].Flags, []uint{0})
	assert.EqualValues(t, board.Board[7][8].Flags, []uint{1})
}
