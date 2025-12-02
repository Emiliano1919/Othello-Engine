package main

import (
	"math"
	"math/rand"
	"time"
)

// SelectPUCT traverses tree until a leaf node is found using PUCT.
func SelectPUCT(node *PUCTNode, c float64) *PUCTNode {
	for node.IsFullyExpandedPUCT() && !node.IsTerminalPUCT() {
		node = BestPUCT(node, c)
	}
	return node
}

// ExpandLeafPUCT expands the  leaf node if there are moves left to try, creating new children.
func ExpandLeafPUCT(node *PUCTNode) *PUCTNode {
	if node.IsTerminalPUCT() {
		return node
	}
	return node.ExpandPUCT()
}

// ExpandPUCT returns an unexplored child of the current node.
func (node *PUCTNode) ExpandPUCT() *PUCTNode {
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
	child := NewPUCTNode(nextState, node, move)
	node.Children = append(node.Children, child)
	return child
}

// rewardFor is an aid function that calculates the reward for the current PUCTNode
func rewardFor(parentBlackTurn bool, result WinState) float64 {
	switch result {
	case DRAW:
		return 0.5
	case BLACK_WIN:
		if parentBlackTurn {
			return 1
		}
		return 0
	case WHITE_WIN:
		if !parentBlackTurn {
			return 1
		}
		return 0
	}
	return 0
}

// BackpropagatePUCT updates visits and wins (A tie counts as 0.5)
func BackpropagatePUCT(node *PUCTNode, result WinState) {
	for n := node; n != nil; n = n.Parent {
		n.Visits++
		p := n.Parent
		if p == nil {
			continue // Skip to next iteration where n will be nil
			// End BackpropagatePUCT
		}
		p.N[n.Move]++
		reward := rewardFor(p.GameState.BlackTurn, result)
		q := p.Q[n.Move]
		p.Q[n.Move] += (reward - q) / float64(n.Visits) // Increment the running average
	}
}

// BestPUCT returns the best child node, of the curent node, according to the PUCT equation.
func BestPUCT(node *PUCTNode, c float64) *PUCTNode {
	var bestChildNode *PUCTNode
	bestPUCT := -math.MaxFloat64
	for _, child := range node.Children {
		move := child.Move
		behaviorPolicy := node.P[move]
		estimatedValueOfAction := node.Q[move]
		totalVisitsOfNode := float64(node.Visits)
		totalVisitsOfAction := float64(node.N[move])

		childPUCT := estimatedValueOfAction + (c*behaviorPolicy)*math.Sqrt(float64(totalVisitsOfNode))/(1+totalVisitsOfAction)

		if childPUCT > bestPUCT {
			bestPUCT = childPUCT
			bestChildNode = child
		}
	}
	return bestChildNode
}

// BestNodeFromMCTSPUCT selects the best child node, move, (The one with most visits) once MCTS has Backpropagated.
func BestNodeFromMCTSPUCT(node *PUCTNode) *PUCTNode {
	var bestNode *PUCTNode
	maxVisits := -1
	for _, child := range node.Children {
		if child.Visits > maxVisits {
			maxVisits = child.Visits
			bestNode = child
		}
	}
	return bestNode
}

// MonteCarloTreeSearchPUCT determines the best move, from the current state/node, using MCTS with PUCT equation.
func MonteCarloTreeSearchPUCT(currentRoot *PUCTNode, iterations int, rng *rand.Rand) *PUCTNode {
	if currentRoot.IsTerminalPUCT() {
		return currentRoot
	}
	for i := 0; i < iterations; i++ {
		selected := SelectPUCT(currentRoot, 2.0)
		nodeToSimulateFrom := ExpandLeafPUCT(selected)
		result := SimulateRollout(nodeToSimulateFrom.GameState, rng)
		BackpropagatePUCT(nodeToSimulateFrom, result)
	}
	return BestNodeFromMCTSPUCT(currentRoot)
}

// OriginalMCTSWinsPlayoutsByMove returns back the number of visits by move after MCTS PUCT
// Returns the updated statistics after doing MCTS PUCT of the moves from the current position.
// It is used for Single run parallelization MCTS PUCT.
func MCTSPUCTWinsPlayoutsByMove(currentRoot *PUCTNode, iterations int, rng *rand.Rand) map[uint8]int {
	if currentRoot.IsTerminalPUCT() {
		return nil
	}
	for i := 0; i < iterations; i++ {
		selected := SelectPUCT(currentRoot, 2.0)
		nodeToSimulateFrom := ExpandLeafPUCT(selected)
		result := SimulateRollout(nodeToSimulateFrom.GameState, rng)
		BackpropagatePUCT(nodeToSimulateFrom, result)
	}
	return currentRoot.N // We return a map with the information needed to make the final decision
	// We just need to return the visits
}

// RootAfterMCTSPUCT returns the root, with the updated info, instead of the best move.
func RootAfterMCTSPUCT(currentRoot *PUCTNode, iterations int, rng *rand.Rand) *PUCTNode {
	if currentRoot.IsTerminalPUCT() {
		return currentRoot
	}
	for i := 0; i < iterations; i++ {
		selected := SelectPUCT(currentRoot, 2.0)
		nodeToSimulateFrom := ExpandLeafPUCT(selected)
		result := SimulateRollout(nodeToSimulateFrom.GameState, rng)
		BackpropagatePUCT(nodeToSimulateFrom, result)
	}
	return currentRoot
}

// SingleRunParallelizationMCTSPUCT is a root level parallelization of MCTS PUCT.
// It works by generating one master tree and at the same time running in parallel simulations from a root copy
// that will be used to update the first level of the master tree (the children )
// with statistics from the parallel simulations.
// This method decreases the variance according to research.
func SingleRunParallelizationMCTSPUCT(currentRoot *PUCTNode, iterationsPerRoutine int, baseRNG *rand.Rand) *PUCTNode {
	maxProcesses := 9
	firstLayerRes := make(chan map[uint8]int, maxProcesses)
	masterTree := make(chan *PUCTNode, 1)
	go func() {
		masterTree <- RootAfterMCTSPUCT(currentRoot, iterationsPerRoutine, baseRNG)
	}()
	for i := 0; i < maxProcesses; i++ {
		go func(id int) {
			// We provide a different rng per parallelization to avoid having the same results.
			parallelRNGi := rand.New(rand.NewSource(time.Now().UnixNano() + int64(id)))
			var emptyMove uint8
			broadcastedNode := NewPUCTNode(currentRoot.GameState, nil, emptyMove)
			// We just need the state of the root (the tree can be generated of it), we don't care about the parent of this one
			// We need to do this because if we share the original root there will be race conditions
			parallelResult := MCTSPUCTWinsPlayoutsByMove(broadcastedNode, iterationsPerRoutine, parallelRNGi)
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
				// We just add visits by move because the selection of the best move just takes that into account
				// The wins are not added because they are just used to decide the exploration/explotation in the process of estimation
			}
		}
	}
	return BestNodeFromMCTSPUCT(currentRoot)
}
