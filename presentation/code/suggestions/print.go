package suggestions

import (
	"fmt"
	"time"

	"github.com/Logiraptor/word-bot/ai"
	"github.com/Logiraptor/word-bot/core"
)

func PrintSuggestions(player ai.AI) {
	b := core.NewBoard()
	b.PlaceTiles(core.PlacedTiles{
		Col:  7, Row: 7, Direction: core.Horizontal,
		Word: core.MakeTiles(core.MakeWord("hello"), "xxxxx"),
	})
	b.PlaceTiles(core.PlacedTiles{
		Col:  11, Row: 6, Direction: core.Vertical,
		Word: core.MakeTiles(core.MakeWord("wrld"), "xxxx"),
	})
	bag := core.NewConsumableBag()
	rack := core.NewConsumableRack(core.MakeTiles(core.MakeWord("abc"), "xxx"))
	player.FindMove(b, bag, rack, func(turn core.Turn) bool {
		return true
	})
	numIterations := 10
	fmt.Printf("Running %s ai for %d iterations\n", player.Name(), numIterations)

	for _, x := range b.ValidatedMoves[:numIterations] {
		tempBoard := b.Clone()
		tempBoard.PlaceTiles(x)
		tempBoard.Print()
		fmt.Println("------------------------")
		time.Sleep(time.Second)
	}
	fmt.Printf("%s ran for %d total iterations\n", player.Name(), len(b.ValidatedMoves))
}
