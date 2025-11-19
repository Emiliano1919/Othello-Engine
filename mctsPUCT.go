package main

import (
	"math"
	"math/rand"
)

// Traverse until a leaf using PUCT
func SelectPUCT(node *PUCTNode, c float64) *PUCTNode {
	for node.IsFullyExpandedPUCT() && !node.IsTerminalPUCT() {
		node = BestPUCT(node, c)
	}
	return node
}

// Expand if there are moves left to try, creating new children
func ExpandLeafPUCT(node *PUCTNode) *PUCTNode {
	if node.IsTerminalPUCT() {
		return node
	}
	return node.ExpandPUCT()
}

// Return an unexplored child of the current node
func (node *PUCTNode) ExpandPUCT() *PUCTNode {
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
	child := NewPUCTNode(nextState, node, move)
	node.Children = append(node.Children, child)
	return child
}

// Update visits and wins (A tie counts as 0.5)
func BackpropagatePUCT(node *PUCTNode, result WinState) {
	var reward float64
	for node != nil {
		node.Visits++
		if node.Parent != nil && node.Parent.GameState.BlackTurn {
			node.Parent.N[node.Move]++
			switch result {
			case WHITE_WIN:
				reward = 0
			case BLACK_WIN:
				reward = 1 // Otherwise optmize for draw
			case DRAW:
				reward = 0.5
			}
			node.Parent.Q[node.Move] += (reward - node.Parent.Q[node.Move]) / float64(node.Visits)
		}

		if node.Parent != nil && !node.Parent.GameState.BlackTurn {
			node.Parent.N[node.Move]++
			switch result {
			case WHITE_WIN:
				reward = 1 // If the machine is white optimize for white
			case BLACK_WIN:
				reward = 0
			case DRAW:
				reward = 0.5
			}
			node.Parent.Q[node.Move] += (reward - node.Parent.Q[node.Move]) / float64(node.Visits)
		}

		node = node.Parent
	}
}

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

// Selection of the best node  (The one with most visits) once MCTS has Backpropagated
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

// Montecarlo Tree Search Implemented correctly
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
