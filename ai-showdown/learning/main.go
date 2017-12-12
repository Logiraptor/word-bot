package main

//#include <stdlib.h>
import "C"
import (
	"fmt"
	"sync"
	"unsafe"

	"github.com/Logiraptor/word-bot/ai"
	"github.com/Logiraptor/word-bot/definitions"
	"github.com/Logiraptor/word-bot/wordlist"

	"github.com/Logiraptor/word-bot/core"
)

var sm *ai.SmartyAI

func init() {
	builder := wordlist.NewTrieBuilder(151434)
	err := definitions.LoadWords("../../words.txt", builder)
	if err != nil {
		panic(err)
	}

	wordDB := builder.Build()
	sm = ai.NewSmartyAI(wordDB, wordDB)
}

type GameContext struct {
	board    *core.Board
	nextRack core.Rack
	lastRack core.Rack
	leave    []core.Tile
	bag      core.Bag
}

var allCtxs = make(map[int]*GameContext)
var ctxCounter = 0
var ctxLock sync.RWMutex

func putContext(gc *GameContext) int {
	ctxLock.Lock()
	ctxCounter++
	allCtxs[ctxCounter] = gc
	ctxLock.Unlock()
	return ctxCounter
}

func getContext(ctx int) *GameContext {
	ctxLock.RLock()
	context := allCtxs[ctx]
	ctxLock.RUnlock()
	return context
}

//export MakeContext
func MakeContext() int {
	fullBag := core.NewConsumableBag()
	fullBag = fullBag.Shuffle()
	firstBag, firstRack := fullBag.FillRack(nil, 7)
	firstBag, secondRack := firstBag.FillRack(nil, 7)

	gc := &GameContext{
		board:    core.NewBoard(),
		nextRack: core.NewConsumableRack(firstRack),
		lastRack: core.NewConsumableRack(secondRack),
		bag:      firstBag,
	}
	return putContext(gc)
}

//export FreeContext
func FreeContext(ctx int) {
	ctxLock.Lock()
	delete(allCtxs, ctx)
	ctxLock.Unlock()
}

//export PrintContext
func PrintContext(ctx int) {
	context := getContext(ctx)
	context.board.Print()
	fmt.Println("Next Rack:", core.Tiles2String(context.nextRack.Rack))
	fmt.Println("Last Rack:", core.Tiles2String(context.lastRack.Rack))
	fmt.Println("Leave:", core.Tiles2String(context.leave))
}

//export GenerateMoves
func GenerateMoves(ctx int, elements **int, numElements *int) {
	context := getContext(ctx)

	outgoingContexts := make([]int, 0, 10)
	sm.GenerateMoves(context.board, context.nextRack, func(turn core.Turn) bool {
		switch v := turn.(type) {
		case core.ScoredMove:
			var (
				newBoard *core.Board
				newRack  core.Rack
				newBag   core.Bag
			)

			newBoard = context.board.Clone()
			newBoard.PlaceTiles(v.PlacedTiles)

			newRack, _ = context.nextRack.Play(v.Word)
			leave := newRack.Rack
			newBag, newRack.Rack = context.bag.FillRack(newRack.Rack, 7-len(newRack.Rack))

			outgoingContexts = append(outgoingContexts, putContext(&GameContext{
				board:    newBoard,
				bag:      newBag,
				nextRack: context.lastRack,
				lastRack: newRack,
				leave:    leave,
			}))
		case core.Pass:
			outgoingContexts = append(outgoingContexts, putContext(&GameContext{
				bag:      context.bag,
				board:    context.board,
				nextRack: context.lastRack,
				lastRack: context.nextRack,
				leave:    context.nextRack.Rack,
			}))
		case core.Exchange:

		}
		return true
	})

	fmt.Println("Generated", len(outgoingContexts), "Moves")
	fmt.Println(outgoingContexts)

	*numElements = len(outgoingContexts)
	elemPtr := C.malloc(C.sizeof_longlong * C.size_t(len(outgoingContexts)))
	elemSlice := (*[1 << 30]int)(elemPtr)
	copy(elemSlice[:], outgoingContexts)
	*elements = (*int)(elemPtr)
}

//export ConvertToTensor
func ConvertToTensor(ctx int, output **float64, length *int) {
	context := getContext(ctx)
	// board tiles, bonuses, leave + bag
	var tensor [15*15*3 + 27*2]float64
	for row := 0; row < 15; row++ {
		for col := 0; col < 15; col++ {
			cell := context.board.Cells[row][col]
			if !cell.Tile.IsNoTile() {
				tensor[row*15+col] = float64(cell.Tile.ToLetter())
				tensor[row*15+col+15*15+15*15] = float64(cell.Tile.PointValue())
			} else {
				tensor[row*15+col+15*15] = float64(context.board.Cells[row][col].Bonus)
			}
		}
	}

	for _, t := range context.leave {
		if t.IsBlank() {
			tensor[15*15*3+26]++
		} else {
			tensor[15*15*3+t.ToLetter()]++
		}
	}

	for _, t := range context.bag.Remaining() {
		if t.IsBlank() {
			tensor[15*15*3+27+26]++
		} else {
			tensor[15*15*3+27+t.ToLetter()]++
		}
	}

	*length = len(tensor)
	elemPtr := C.malloc(C.sizeof_double * C.size_t(len(tensor)))
	elemSlice := (*[1 << 30]float64)(elemPtr)
	copy(elemSlice[:], tensor[:])
	*output = (*float64)(elemPtr)
}

//export FreeContextBuffer
func FreeContextBuffer(elements *int) {
	C.free(unsafe.Pointer(elements))
}

//export FreeTensorBuffer
func FreeTensorBuffer(elements *float64) {
	C.free(unsafe.Pointer(elements))
}

// I need:
// X Some representation of a board, rack, and bag (game context?)
// X a way to allocate / deallocate one from C code.
// X a way to generate all valid moves in that context
// a way to turn the game context into a tensor

func main() {}
