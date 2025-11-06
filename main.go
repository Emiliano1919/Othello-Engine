package main

// Game loop
func main() {
	initialNode := InitialRootNode()
	initialNode.GameState.Boards.PrintBoard()
	userIsBlack := RequestUserIsBlack()
	var node *Node
	if !userIsBlack {
		bestOpening := MonteCarloTreeSearch(initialNode, 5000)
		bestOpening.GameState.Boards.PrintBoard()
		node = bestOpening
	} else {
		initialNode.GameState.PrintBoardWithMoves()
		blackMove := RequestMove(userIsBlack)
		node = NextNodeFromInput(initialNode, blackMove)
	}
	for !node.IsTerminal() {
		if !userIsBlack {
			if !node.GameState.BlackTurn {
				node.GameState.PrintBoardWithMoves()
				whiteMove := RequestMove(userIsBlack)
				node = NextNodeFromInput(node, whiteMove)
			} else {
				mctsNode := MonteCarloTreeSearch(node, 5000)
				mctsNode.GameState.Boards.PrintBoard()
				node = mctsNode
			}
		} else {
			if node.GameState.BlackTurn {
				node.GameState.PrintBoardWithMoves()
				blackMove := RequestMove(userIsBlack)
				node = NextNodeFromInput(node, blackMove)
			} else {
				mctsNode := MonteCarloTreeSearch(node, 5000)
				mctsNode.GameState.Boards.PrintBoard()
				node = mctsNode
			}
		}
	}
	if node.IsTerminal() {
		OutputResult(node)
	}
}
