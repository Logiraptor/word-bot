package ai

import (
	"fmt"
	"math/rand"

	"github.com/Logiraptor/word-bot/core"
	"github.com/Logiraptor/word-bot/persist"
)

type Player struct {
	ai    AI
	name  string
	rack  core.Rack
	score core.Score
}

func NewPlayer(ai AI) *Player {
	return &Player{
		ai:   ai,
		name: ai.Name(),
		rack: core.NewConsumableRack(nil),
	}
}

func (p *Player) takeTurn(wordDB core.WordList, board *core.Board, bag core.Bag) (core.Bag, []core.Tile, core.ScoredMove, bool) {
	var turn core.Turn
	p.ai.FindMove(board, bag, p.rack, func(t core.Turn) bool {
		turn = t
		return true
	})
	if turn == nil {
		return bag, nil, core.ScoredMove{}, false
	}

	switch move := turn.(type) {
	case core.ScoredMove:
		if !board.ValidateMove(move.PlacedTiles, wordDB) {
			fmt.Printf("%s played an invalid move: %v!\n", p.name, move)
			return bag, nil, core.ScoredMove{}, false
		}

		newRack, ok := p.rack.Play(move.Word)
		if !ok {
			return bag, nil, core.ScoredMove{}, false
		}

		p.rack = newRack

		score := board.Score(move.PlacedTiles)
		board.PlaceTiles(move.PlacedTiles)

		leave := newRack.Rack
		bag, p.rack.Rack = bag.FillRack(p.rack.Rack, 7-len(p.rack.Rack))

		p.score += score

		return bag, leave, move, true
	case core.Pass:
		return bag, nil, core.ScoredMove{}, false
	case core.Exchange:
		newRack := core.NewConsumableRack(nil)
		bag, newRack.Rack = bag.FillRack(newRack.Rack, 7-len(newRack.Rack))
		bag = bag.Replace(p.rack.Rack)
		p.rack = newRack
		bag, p.rack.Rack = bag.FillRack(p.rack.Rack, 7-len(p.rack.Rack))
		return bag, nil, core.ScoredMove{}, false
	default:
		panic(fmt.Sprintf("%s played unknown turn type: %#v", p.ai.Name(), move))
	}
}

func PlayGame(wordDB core.WordList, a, b func(board *core.Board) *Player) persist.Game {
	game := persist.Game{}

	swapped := false
	if rand.Intn(2) == 0 {
		swapped = true
		a, b = b, a
	}

	board := core.NewBoard()
	p1 := a(board)
	p2 := b(board)

	bag := core.NewConsumableBag().Shuffle()

	bag, p1.rack.Rack = bag.FillRack(p1.rack.Rack, 7)
	bag, p2.rack.Rack = bag.FillRack(p2.rack.Rack, 7)

	var (
		p1Ok, p2Ok = true, true
		move       core.ScoredMove
		leave      []core.Tile
	)

	for bag.Count() > 0 && (p1Ok || p2Ok) {
		if bag, leave, move, p1Ok = p1.takeTurn(wordDB, board, bag); p1Ok {
			game.AddMove(p1.name, leave, move)
		}
		if bag, leave, move, p2Ok = p2.takeTurn(wordDB, board, bag); p2Ok {
			game.AddMove(p2.name, leave, move)
		}
	}

	if swapped {
		p1, p2 = p2, p1
	}

	return game
}
