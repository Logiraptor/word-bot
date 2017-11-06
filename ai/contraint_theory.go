package ai

import "github.com/Logiraptor/word-bot/core"
import "github.com/Logiraptor/word-bot/wordlist"

type CT struct {
	wordDB *wordlist.Gaddag
}

func (c *CT) GenerateMoves(b *core.Board, rack core.Rack, onMove func(core.Turn) bool) {
	var boardConstraint [15][15]Constraint
	for i := 0; i < 15; i++ {
		for j := 0; j < 15; j++ {
			boardConstraint[i][j] = PermittedTiles(b, c.wordDB, i, j)
		}
	}
	for row := 0; row < 15; row++ {
		for col := 0; col < 15; col++ {
			mg := ctmg{
				board:            b,
				boardConstraints: boardConstraint,
				rack:             rack,
				onMove:           onMove,
				wordDB:           c.wordDB,
				anchorRow:        row,
				anchorCol:        col,
			}
			mg.generateHorizontalMoves(row, col, nil)
			mg.generateVerticalMoves(nil)
		}
	}
}

type ctmg struct {
	wordDB               *wordlist.Gaddag
	board                *core.Board
	boardConstraints     [15][15]Constraint
	rack                 core.Rack
	onMove               func(core.Turn) bool
	anchorRow, anchorCol int
}

func (c ctmg) generateHorizontalMoves(row, col int, prev []core.Tile) {

	if c.wordDB.CanReverse() && c.wordDB.Reverse().IsTerminal() {
		word := make([]core.Tile, len(prev))
		copy(word, prev)
		pt := core.PlacedTiles{
			Col:       c.anchorCol,
			Row:       c.anchorRow,
			Direction: core.Horizontal,
			Word:      word,
		}
		c.onMove(core.ScoredMove{
			PlacedTiles: pt,
			Score:       c.board.Score(pt),
		})
	}

	for rackIndex, tile := range c.rack.Rack {
		if !c.rack.CanConsume(rackIndex) {
			continue
		}

		lowerBound, upperBound := tile, tile
		if tile.IsBlank() {
			lowerBound = blankA
			upperBound = blankZ
		}
		for tile := lowerBound; tile <= upperBound; tile++ {
			if c.canPlace(row, col, tile) {
				c.place(rackIndex, tile).generateHorizontalMoves(row, col+1, append(prev, tile))
			}
		}
	}
}

func (c ctmg) canPlace(row, col int, t core.Tile) bool {
	if !c.wordDB.CanBranch(t) {
		return false
	}
	if !c.boardConstraints[row][col].AllowsTile(t) {
		return false
	}
	return true
}

func (c ctmg) place(rackIndex int, t core.Tile) ctmg {
	return ctmg{
		wordDB:           c.wordDB.Branch(t),
		board:            c.board,
		boardConstraints: c.boardConstraints,
		rack:             c.rack.Consume(rackIndex),
		onMove:           c.onMove,
		anchorRow:        c.anchorRow,
		anchorCol:        c.anchorCol,
	}
}

func (c ctmg) generateVerticalMoves(prev []core.Tile) {

}
