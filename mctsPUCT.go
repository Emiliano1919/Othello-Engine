package main

import (
	"math"
	"math/rand"
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
