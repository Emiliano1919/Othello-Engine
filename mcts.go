package main

import (
	"math"
	"math/rand"
	"time"
)

// OptimizeFor declares for who you need to optimize the moves for.
type OptimizeFor int

const (
	OPTIMIZE_FOR_BLACK OptimizeFor = iota
	OPTIMIZE_FOR_WHITE
)

// Expand returns an unexplored child of the current node.
// If the node has been completely explored it returns nil.
func (node *Node) Expand() *Node {
	if len(node.UntriedMoves) == 0 {
		return nil
	}
	// Pop
	move := node.UntriedMoves[len(node.UntriedMoves)-1]
	node.UntriedMoves = node.UntriedMoves[:len(node.UntriedMoves)-1]

	// We generate a new node because make move does not generate a new board by default
	newBoards := node.GameState.Boards // Copy
	newBoards.MakeMoveIndex(node.GameState.BlackTurn, move)

	nextState := State{
		Boards:    newBoards,
		BlackTurn: !node.GameState.BlackTurn, // SWITCH TURN
	}

	// Edge case: (If in the next turn the player cannot make a move switch the turn again)

	if !nextState.Boards.HasValidMove(nextState.BlackTurn) && nextState.Boards.HasValidMove(!nextState.BlackTurn) {
		nextState.BlackTurn = !nextState.BlackTurn
	}
	// Generate the child with the new values and add it to the list of children of the node
	child := NewNode(nextState, node, move)
	node.Children = append(node.Children, child)
	return child
}

// AgressiveTraverse returns a new node using a modified MCTS traversal.
// The AgressiveBestUCT is biased for explotation and we return not the leaf but the expansion of it.
func AgressiveTraverse(node *Node) *Node {
	for node.IsFullyExpanded() && !node.IsTerminal() {
		node = AgressiveBestUCT(node, float64(2))
	}
	if node.IsTerminal() {
		return node
	}
	return node.Expand()
}

// Select traverses until reaching a leaf using OriginalBestUCT.
func Select(node *Node, c float64) *Node {
	for node.IsFullyExpanded() && !node.IsTerminal() {
		node = OriginalBestUCT(node, c)
	}
	return node
}

// ExpandLeaf expands the node if there are moves left to try, by creating new children.
func ExpandLeaf(node *Node) *Node {
	if node.IsTerminal() {
		return node
	}
	return node.Expand()
}

// OriginalTraverse is selection and expansion of MCTS in one.
// This function was split for clarity (DO NOT USE).
// Traverse the Montecarlo tree using the best UCT, when you find a leaf node expand it.
func OriginalTraverse(node *Node) *Node {
	for node.IsFullyExpanded() && !node.IsTerminal() {
		node = OriginalBestUCT(node, float64(2))
	}
	if node.IsTerminal() {
		return node
	}
	return node.Expand()
}

// AgressiveBestUCT chooses the best child to explore using a modified UCT.
// The modified UCT gives more weight to explotation resulting in more simulation for early good finds.
func AgressiveBestUCT(node *Node, c float64) *Node {
	var best *Node
	bestUCT := float64(-1 << 63)
	for _, child := range node.Children {
		explotationTerm := float64(child.Wins) / float64(child.Visits)
		explorationTerm := math.Sqrt(math.Log(float64(node.Visits)) / float64(child.Visits))
		C := math.Sqrt(c) // Theoretical value, will try to find a better one through self play
		//UCTValue := C*explorationTerm + explotationTerm // This is the correct formula
		UCTValue := explorationTerm + C*explotationTerm // Formula that produces better results but is not the usual

		if UCTValue > bestUCT {
			bestUCT = UCTValue
			best = child
		}
	}
	return best
}

// OriginalBestUCT chooses the best child to explore using UCT.
func OriginalBestUCT(node *Node, c float64) *Node {
	var best *Node
	bestUCT := float64(-1 << 63)
	for _, child := range node.Children {
		explotationTerm := float64(child.Wins) / float64(child.Visits)
		explorationTerm := math.Sqrt(math.Log(float64(node.Visits)) / float64(child.Visits))
		C := math.Sqrt(c)                               // Theoretical value, will try to find a better one through self play
		UCTValue := C*explorationTerm + explotationTerm // This is the correct formula

		if UCTValue > bestUCT {
			bestUCT = UCTValue
			best = child
		}
	}
	return best
}

// SimulateRollout simulates games from the current node to end game.
// Choosing randomly at each move.
// The states explored here are done inline (see README.md) so it is inexpensive in memory.
func SimulateRollout(state State, random *rand.Rand) WinState {
	current := state

	for !IsTerminalState(current) {
		var moves uint64
		if current.BlackTurn {
			moves = generateMoves(current.Boards.Black, current.Boards.White)
		} else {
			moves = generateMoves(current.Boards.White, current.Boards.Black)
		}

		if moves == 0 {
			// no moves, then pass turn
			current.BlackTurn = !current.BlackTurn
			continue
		}

		moveArray := FastArrayOfMoves(moves)
		move := moveArray[random.Intn(len(moveArray))] // Here is the rollout ppolicy  which is random

		current.Boards.MakeMoveIndex(current.BlackTurn, move)
		current.BlackTurn = !current.BlackTurn
	}
	// 1 = Black win, 0 = White win, 2 = draw
	return WinnerState(current)
}

// InnacurateBackpropagate modified backpropagate that backpropagates a result to the root of the current tree.
// Update visits and wins (A tie also counts as a win).
// It is innacurate because it only updates the levels of the color it wants.
func InnacurateBackpropagate(node *Node, result WinState, optimizeFor OptimizeFor) {
	for node != nil {
		node.Visits++
		switch optimizeFor {
		case OPTIMIZE_FOR_BLACK:
			if node.Parent != nil && node.Parent.GameState.BlackTurn {
				switch result {
				case WHITE_WIN:
					node.Wins += 0
				case BLACK_WIN:
					node.Wins++ // Otherwise optmize for draw
				case DRAW:
					node.Wins++
				}
			}
		case OPTIMIZE_FOR_WHITE:
			if node.Parent != nil && !node.Parent.GameState.BlackTurn {
				switch result {
				case WHITE_WIN:
					node.Wins += 2 // If the machine is white optimize for white
				case BLACK_WIN:
					node.Wins += 0
				case DRAW:
					node.Wins++
				}
			}
		}
		node = node.Parent
	}
}

// OriginalBackpropagate backpropagates the results to each node until the root of the given tree is reached.
// Update visits and wins (A tie also counts as a win).
func OriginalBackpropagate(node *Node, result WinState) {
	for n := node; n != nil; n = n.Parent {
		n.Visits++
		if p := n.Parent; p != nil {
			isBlackTurn := p.GameState.BlackTurn
			isBlackWin := (result == BLACK_WIN || result == DRAW)
			isWhiteWin := (result == WHITE_WIN || result == DRAW)
			if (isBlackTurn && isBlackWin) || (!isBlackTurn && isWhiteWin) {
				n.Wins++
			}
		}
	}
}

// BestNodeFromMCTS returns the child node, of the current node, with the most visits.
// Used once MCTS has Backpropagated the results updating the statistics.
func BestNodeFromMCTS(node *Node) *Node {
	var bestNode *Node
	maxVisits := -1
	for _, child := range node.Children {
		if child.Visits > maxVisits {
			maxVisits = child.Visits
			bestNode = child
		}
	}
	return bestNode
}

// InnacurateMonteCarloTreeSearch is a modified version of MCTS, fun to play against.
// Montecarlo Tree search algorithm with agressive UCT (traverse), double expansion and incorrect backpropagation.
// This model is more fun to play against than normal MCTS.
func InnacurateMonteCarloTreeSearch(currentRoot *Node, iterations int, optimizeFor OptimizeFor, rng *rand.Rand) *Node {
	if currentRoot.IsTerminal() {
		return currentRoot
	}
	if optimizeFor == OPTIMIZE_FOR_BLACK {
		for i := 0; i < iterations; i++ {
			leaf := AgressiveTraverse(currentRoot) // Select and Expand are coded in Traverse
			var nodeToSimulateFrom *Node
			if len(leaf.UntriedMoves) > 0 {
				child := leaf.Expand() // This is a double expansion which is not in the normal implementation
				nodeToSimulateFrom = child
			} else {
				nodeToSimulateFrom = leaf
			}

			result := SimulateRollout(nodeToSimulateFrom.GameState, rng)
			InnacurateBackpropagate(nodeToSimulateFrom, result, optimizeFor)
		}
	} else {
		for i := 0; i < iterations; i++ {
			// It does not benefit from being agressive on white, so we use Original Traverseß
			selected := Select(currentRoot, 2.0)
			nodeToSimulateFrom := ExpandLeaf(selected)
			result := SimulateRollout(nodeToSimulateFrom.GameState, rng)
			InnacurateBackpropagate(nodeToSimulateFrom, result, optimizeFor)
		}
	}

	return BestNodeFromMCTS(currentRoot)
}

// OriginalMonteCarloTreeSearch implemented as usual.
// Returns the best move determined by MCTS with UCT.
func OriginalMonteCarloTreeSearch(currentRoot *Node, iterations int, rng *rand.Rand) *Node {
	if currentRoot.IsTerminal() {
		return currentRoot
	}
	for i := 0; i < iterations; i++ {
		selected := Select(currentRoot, 2.0)
		nodeToSimulateFrom := ExpandLeaf(selected)
		result := SimulateRollout(nodeToSimulateFrom.GameState, rng)
		OriginalBackpropagate(nodeToSimulateFrom, result)
	}
	return BestNodeFromMCTS(currentRoot)
}

// RootAfterOriginalMCTS returns the root, with the updated info, instead of the best move.
func RootAfterOriginalMCTS(currentRoot *Node, iterations int, rng *rand.Rand) *Node {
	if currentRoot.IsTerminal() {
		return currentRoot
	}
	for i := 0; i < iterations; i++ {
		selected := Select(currentRoot, 2.0)
		nodeToSimulateFrom := ExpandLeaf(selected)
		result := SimulateRollout(nodeToSimulateFrom.GameState, rng)
		OriginalBackpropagate(nodeToSimulateFrom, result)
	}
	return currentRoot
}

// OriginalMCTSWinsPlayoutsByMove returns back the number of games and wins per move from the given node.
// Returns the updated statistics after doing MCTS of the moves from the current position.
// It is used for Single run parallelization MCTS.
func OriginalMCTSWinsPlayoutsByMove(currentRoot *Node, iterations int, rng *rand.Rand) map[uint8]int {
	if currentRoot.IsTerminal() {
		return nil
	}
	for i := 0; i < iterations; i++ {
		selected := Select(currentRoot, 2.0)
		nodeToSimulateFrom := ExpandLeaf(selected)
		result := SimulateRollout(nodeToSimulateFrom.GameState, rng)
		OriginalBackpropagate(nodeToSimulateFrom, result)
	}
	movesWithRatio := make(map[uint8]int, len(currentRoot.Children))
	for _, child := range currentRoot.Children {
		movesWithRatio[child.Move] = child.Visits
	}
	return movesWithRatio // We return a map with the information needed
}

// SingleRunParallelizationMCTS is a root level parallelization of MCTS.
// It works by generating one master tree and at the same time running in parallel simulations that will be used
// to update the first level of the master tree (the children ) with statistics from the parallel simulations.
// This method decreases the variance according to research.
func SingleRunParallelizationMCTS(currentRoot *Node, iterationsPerRoutine int, baseRNG *rand.Rand) *Node {
	maxProcesses := 9
	firstLayerRes := make(chan map[uint8]int, maxProcesses)
	masterTree := make(chan *Node, 1)
	go func() {
		masterTree <- RootAfterOriginalMCTS(currentRoot, iterationsPerRoutine, baseRNG)
	}()
	for i := 0; i < maxProcesses; i++ {
		go func(id int) {
			parallelRNGi := rand.New(rand.NewSource(time.Now().UnixNano() + int64(id)))
			var emptyMove uint8
			broadcastedNode := NewNode(currentRoot.GameState, nil, emptyMove)
			// We just need the state of the root (the tree can be generated of it), we don't care about the parent of this one
			// We need to do this because if we share the original root there will be race conditions
			parallelResult := OriginalMCTSWinsPlayoutsByMove(broadcastedNode, iterationsPerRoutine, parallelRNGi)
			firstLayerRes <- parallelResult
		}(i)
	}
	// Wait for all of them to complete
	currentRoot = <-masterTree // update the tree with the result of 1 simulation (this will be the base)
	close(masterTree)
	// Collect exactly maxProcesses results from worker goroutines
	for i := 0; i < maxProcesses; i++ {
		dict := <-firstLayerRes
		// Accumulate worker results into master's children
		for _, child := range currentRoot.Children {
			if res, exists := dict[child.Move]; exists {
				child.Visits += res
			}
		}
	}
	return BestNodeFromMCTS(currentRoot)
}

// func leafParallelizationMCTS(currentRoot *Node, iterationsPerRoutine int) *Node {
// Run multiple go routines once you get to a leaf
// }
