package ai

import (
	"github.com/Logiraptor/word-bot/core"
	"github.com/Logiraptor/word-bot/wordlist"
)

type Constraint int

func (c Constraint) AllowsTile(t core.Tile) bool {
	return c&(1<<uint(t.ToLetter())) > 0
}

const allLetters = 0x3FFFFFF

func PermittedTiles(board *core.Board, wordDB *wordlist.Gaddag, row, col int) Constraint {
	return permittedTilesByDir(board, wordDB, row, col, core.Horizontal) &
		permittedTilesByDir(board, wordDB, row, col, core.Vertical)
}

func permittedTilesByDir(board *core.Board, wordDB *wordlist.Gaddag, row, col int, dir core.Direction) Constraint {
	dRow, dCol := dir.Offsets()

	if !board.HasTile(row+dRow, col+dCol) && !board.HasTile(row-dRow, col-dCol) {
		return allLetters
	}

	// scan forward
	r, c := row+dRow, col+dCol
	for board.HasTile(r, c) {
		t := board.Cells[r][c].Tile
		if !wordDB.CanBranch(t) {
			return 0
		}
		r += dRow
		c += dCol
		wordDB = wordDB.Branch(t)
	}
	// reverse
	if !wordDB.CanReverse() {
		return 0
	}
	wordDB = wordDB.Reverse()

	constr := Constraint(0)
	// for each tile t that can branch
outer:
	for i := blankA; i <= blankZ; i++ {
		tmpWordDB := wordDB
		if !tmpWordDB.CanBranch(i) {
			continue
		}

		tmpWordDB = tmpWordDB.Branch(i)

		// scan backward
		r, c := row-dRow, col-dCol
		for board.HasTile(r, c) {
			t := board.Cells[r][c].Tile
			if !tmpWordDB.CanBranch(t) {
				continue outer
			}
			r -= dRow
			c -= dCol
			tmpWordDB = tmpWordDB.Branch(t)
		}
		// if terminal, add that t to the set
		if tmpWordDB.IsTerminal() {
			constr |= 1 << uint(i.ToLetter())
		}
	}
	return constr
}
