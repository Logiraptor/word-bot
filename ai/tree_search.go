package ai

// // For this AI, we will use Monte Carlo Tree Search
// // In order to simplify the first attempt we will make some assumptions:
// // We will only consider the top 10 ai moves at each turn
// // We will assume that the opponent always simply chooses the highest scoring available move

// import "github.com/Logiraptor/word-bot/core"
// import "math"

// type MoveFinder interface {
// 	FindMoves(board *core.Board, rack core.ConsumableRack, numMoves int) []core.ScoredMove
// }

// type GamePlayer interface {
// 	PlayGame(board *core.Board, bag core.ConsumableBag, rack core.ConsumableRack) (win bool)
// }

// type TreeSearchAI struct {
// 	moveFinder MoveFinder
// 	root       *AINode
// 	bias       float64
// }

// func NewTreeSearchAI(moveFinder MoveFinder, bias float64) *TreeSearchAI {
// 	return &TreeSearchAI{
// 		moveFinder: moveFinder,
// 		bias:       bias,
// 		root:       nil,
// 	}
// }

// type Node struct {
// 	Wins, Simulations int
// 	NumVisits         int
// }

// func (n Node) EstimatedValue() float64 {
// 	return float64(n.Wins) / float64(n.Simulations)
// }

// func (n Node) UCB(bias float64) float64 {
// 	return upperConfidenceBound(n.EstimatedValue(), bias, float64(n.NumVisits), float64(n.Parent.NumVisits))
// }

// type AINode struct {
// 	Node
// 	Parent   *OpponentNode
// 	Move     core.ScoredMove
// 	Rack     core.ConsumableRack
// 	Children []*OpponentNode
// }

// func (a *AINode) UCB(bias float64) float64 {
// 	if a == nil {
// 		return 0
// 	}
// 	return a.Node.UCB(bias)
// }

// type OpponentNode struct {
// 	Node
// 	Parent   *AINode
// 	Move     core.ScoredMove
// 	Rack     core.ConsumableRack
// 	Children []*AINode
// }

// func (o *OpponentNode) UCB(bias float64) float64 {
// 	if o == nil {
// 		return 0
// 	}
// 	return o.Node.UCB(bias)
// }

// type PlayoutResult struct {
// 	Win bool
// }

// func upperConfidenceBound(estimatedValue, bias, nodeVisits, parentVisits float64) float64 {
// 	return estimatedValue + bias*(math.Sqrt(math.Log2(parentVisits)/nodeVisits))
// }

// func (t *TreeSearchAI) step(b *core.Board, r core.ConsumableRack) {
// 	l := t.selectLeaf(b)
// 	t.expandTree(l, b, r)
// 	for _, c := range l.Children {
// 		result := t.playoutNode(c)
// 		t.updateWeights(result, c)
// 	}
// }

// func (t *TreeSearchAI) selectLeaf(b *core.Board) *AINode {
// 	// Starting at root node R, recursively select optimal child nodes (explained below) until a leaf node L is reached.
// 	current := t.root

// 	for current != nil {
// 		// select child with highest UCB
// 		// Opportunity here to use a heap for faster move selection
// 		var bestOpponentMove = current.Children[0]
// 		{
// 			for _, c := range current.Children {
// 				if c.UCB(t.bias) > bestOpponentMove.UCB(t.bias) {
// 					bestOpponentMove = c
// 				}
// 			}
// 		}

// 		var bestAIMove = bestOpponentMove.Children[0]
// 		{
// 			for _, c := range bestOpponentMove.Children {
// 				if c.UCB(t.bias) > bestAIMove.UCB(t.bias) {
// 					bestAIMove = c
// 				}
// 			}
// 		}

// 		b.PlaceTiles(bestOpponentMove.Move.Word, bestOpponentMove.Move.Row, bestOpponentMove.Move.Col, bestOpponentMove.Move.Direction)
// 		b.PlaceTiles(bestAIMove.Move.Word, bestAIMove.Move.Row, bestAIMove.Move.Col, bestAIMove.Move.Direction)

// 		current = bestAIMove
// 	}
// 	return current
// }

// func (t *TreeSearchAI) expandTree(node *AINode, board *core.Board, rack core.ConsumableRack) {
// 	// If L is a not a terminal node (i.e. it does not end the game) then create one or more child nodes and select one C.
// 	moves := t.moveFinder.FindMoves(board, rack, 10)
// }

// func (t *TreeSearchAI) playoutNode(node *AINode) PlayoutResult {
// 	// Run a simulated playout from C until a result is achieved.
// 	return PlayoutResult{Win: true}
// }

// func (t *TreeSearchAI) updateWeights(result PlayoutResult, node *AINode) {
// 	// Update the current move sequence with the simulation result.
// }
