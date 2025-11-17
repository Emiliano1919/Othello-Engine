package main

import (
	"math"
	"math/rand"
	"time"
)

// New OptimizeFor... to increase Clarity
type OptimizeFor int

const (
	OPTIMIZE_FOR_BLACK OptimizeFor = iota
	OPTIMIZE_FOR_WHITE
)

// Return an unexplored child of the current node
func (node *Node) Expand() *Node {
	if len(node.UntriedMoves) == 0 {
		return nil
	}
	// Pop
	move := node.UntriedMoves[len(node.UntriedMoves)-1]
	node.UntriedMoves = node.UntriedMoves[:len(node.UntriedMoves)-1]

	// We generate a new node because make move does not generate a new board by default
	newBoards := node.GameState.Boards // Copy
	newBoards.MakeMove(node.GameState.BlackTurn, move[0], move[1])

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

// Traverse the Montecarlo tree using the best UCT, when you find a leaf node expand it
func AgressiveTraverse(node *Node) *Node {
	for node.IsFullyExpanded() && !node.IsTerminal() {
		node = AgressiveBestUCT(node, float64(2))
	}
	if node.IsTerminal() {
		return node
	}
	return node.Expand()
}

// Traverse the Montecarlo tree using the best UCT, when you find a leaf node expand it
func OriginalTraverse(node *Node) *Node {
	for node.IsFullyExpanded() && !node.IsTerminal() {
		node = OriginalBestUCT(node, float64(2))
	}
	if node.IsTerminal() {
		return node
	}
	return node.Expand()
}

// Choose best child to explore using UCT
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

// Choose best child to explore using UCT
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

// Simulate randomly from current node to the end of the game (choosing randomly at each step)
// The states explored here are not saved, only the result
func SimulateRollout(state State, seed int64) WinState {
	if seed == 0 {
		seed = time.Now().UnixNano() // use a random seed by default
	}
	random := rand.New(rand.NewSource(seed))
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

		moveArray := ArrayOfPositionalMoves(ArrayOfMoves(moves))
		move := moveArray[random.Intn(len(moveArray))] // Here is the rollout ppolicy  which is random

		current.Boards.MakeMove(current.BlackTurn, move[0], move[1])
		current.BlackTurn = !current.BlackTurn
	}
	// 1 = Black win, 0 = White win, 2 = draw
	return WinnerState(current)
}

// Update visits and wins (A tie also counts as a win)
// It is innacurate because it only updates the levels of the color it wants
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
					node.Wins += 1 // Otherwise optmize for draw
				case DRAW:
					node.Wins += 1
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
					node.Wins += 1
				}
			}
		}
		node = node.Parent
	}
}

// Update visits and wins (A tie also counts as a win)
func OriginalBackpropagate(node *Node, result WinState) {
	for node != nil {
		node.Visits++
		if node.Parent != nil && node.Parent.GameState.BlackTurn {
			switch result {
			case WHITE_WIN:
				node.Wins += 0
			case BLACK_WIN:
				node.Wins += 1 // Otherwise optmize for draw
			case DRAW:
				node.Wins += 1
			}
		}

		if node.Parent != nil && !node.Parent.GameState.BlackTurn {
			switch result {
			case WHITE_WIN:
				node.Wins += 1 // If the machine is white optimize for white
			case BLACK_WIN:
				node.Wins += 0
			case DRAW:
				node.Wins += 1
			}
		}

		node = node.Parent
	}
}

// Selection of the best node  (The one with most visits) once MCTS has Backpropagated
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

// Montecarlo Tree search algorithm with agressive UCT (travers), double expansion and incorrect backpropagation
// This model is more fun to play against than normal MCTS
func InnacurateMonteCarloTreeSearch(currentRoot *Node, iterations int, optimizeFor OptimizeFor) *Node {
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

			result := SimulateRollout(nodeToSimulateFrom.GameState, 0)
			InnacurateBackpropagate(nodeToSimulateFrom, result, optimizeFor)
		}
	} else {
		for i := 0; i < iterations; i++ {
			// It does not benefit from being agressive on white, so we use Original Traverseß
			nodeToSimulateFrom := OriginalTraverse(currentRoot) // Select and Expand are coded in Traverse.
			result := SimulateRollout(nodeToSimulateFrom.GameState, 0)
			InnacurateBackpropagate(nodeToSimulateFrom, result, optimizeFor)
		}
	}

	return BestNodeFromMCTS(currentRoot)
}

// Montecarlo Tree Search Implemented correctly
func OriginalMonteCarloTreeSearch(currentRoot *Node, iterations int) *Node {
	if currentRoot.IsTerminal() {
		return currentRoot
	}
	for i := 0; i < iterations; i++ {
		nodeToSimulateFrom := OriginalTraverse(currentRoot) // Select and Expand are coded in Traverse
		result := SimulateRollout(nodeToSimulateFrom.GameState, 0)
		OriginalBackpropagate(nodeToSimulateFrom, result)
	}
	return BestNodeFromMCTS(currentRoot)
}

// Montecarlo Tree Search Implemented correctly
func RootAfterOriginalMCTS(currentRoot *Node, iterations int) *Node {
	if currentRoot.IsTerminal() {
		return currentRoot
	}
	for i := 0; i < iterations; i++ {
		nodeToSimulateFrom := OriginalTraverse(currentRoot) // Select and Expand are coded in Traverse
		result := SimulateRollout(nodeToSimulateFrom.GameState, 0)
		OriginalBackpropagate(nodeToSimulateFrom, result)
	}
	return currentRoot
}

// Send back the number of games and wins per move from the root
func OriginalMCTSWinsPlayoutsByMove(currentRoot *Node, iterations int, seed int64) map[[2]uint8][2]int {
	if currentRoot.IsTerminal() {
		return nil
	}
	for i := 0; i < iterations; i++ {
		nodeToSimulateFrom := OriginalTraverse(currentRoot) // Select and Expand are coded in Traverse
		result := SimulateRollout(nodeToSimulateFrom.GameState, seed)
		OriginalBackpropagate(nodeToSimulateFrom, result)
	}
	movesWithRatio := make(map[[2]uint8][2]int, len(currentRoot.Children))
	for _, child := range currentRoot.Children {
		movesWithRatio[child.Move] = [2]int{(child.Wins), (child.Visits)}
	}
	return movesWithRatio // We return a map with the information needed
}

// // This is a root level parallelization
// It works by generating one master tree and at the same time running in parallel simulations that will be used
// to update the first level of the master tree (the children ) with statistics from the parallel simulations
// This method decreases the variance according to research
func SingleRunParallelizationMCTS(currentRoot *Node, iterationsPerRoutine int) *Node {
	maxProcesses := 9
	firstLayerRes := make(chan map[[2]uint8][2]int, maxProcesses)
	masterTree := make(chan *Node, 1)
	go func() {
		masterTree <- RootAfterOriginalMCTS(currentRoot, iterationsPerRoutine)
	}()
	for i := 0; i < maxProcesses; i++ {
		go func(seed int64) {
			broadcastedNode := NewNode(currentRoot.GameState, nil, [2]uint8{})
			// We just need the state of the root (the tree can be generated of it), we don't care about the parent of this one
			// We need to do this because if we share the original root there will be race conditions
			parallelResult := OriginalMCTSWinsPlayoutsByMove(broadcastedNode, iterationsPerRoutine, seed)
			firstLayerRes <- parallelResult
		}(int64(i))
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
				child.Wins += res[0]
				child.Visits += res[1]
			}
		}
	}
	return BestNodeFromMCTS(currentRoot)
}
