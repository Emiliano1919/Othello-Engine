package main

// Node struct and methods

type Node struct {
	Parent       *Node
	Children     []*Node
	UntriedMoves []uint8
	GameState    State
	Visits       int
	Wins         int
	Move         uint8
}

func InitialRootNode() *Node {
	var node *Node
	var boards Board
	boards.Init()
	var state State
	state.Boards = boards
	state.BlackTurn = true
	var emptyMove uint8
	node = NewNode(state, nil, emptyMove)
	return node
}

func NextNodeFromInput(parent *Node, moveIndex uint8) *Node {

	// If the move exists in a subtree cut that subtree and preserve it, to preserve the information of previous simulations
	for _, child := range parent.Children {
		if child.Move == moveIndex {
			child.Parent = nil
			return child
		}
	}
	newBoards := parent.GameState.Boards                           // Copy
	newBoards.MakeMoveIndex(parent.GameState.BlackTurn, moveIndex) // We make the move on the current player
	newState := State{
		Boards:    newBoards,
		BlackTurn: !parent.GameState.BlackTurn, // switch turn immediately
	}
	// Switch turn but conditionally to deal with edgecase of the opponent having no moves.
	if !newState.Boards.HasValidMove(newState.BlackTurn) && newState.Boards.HasValidMove(!newState.BlackTurn) {
		//fmt.Println("No valid moves for next player — passing turn back.")
		newState.BlackTurn = !newState.BlackTurn
	}
	return NewNode(newState, parent, moveIndex)
}

func NewNode(state State, parent *Node, move uint8) *Node {
	var legalMoves uint64
	if state.BlackTurn {
		legalMoves = generateMoves(state.Boards.Black, state.Boards.White)
	} else {
		legalMoves = generateMoves(state.Boards.White, state.Boards.Black)
	}
	movesFromCurrent := FastArrayOfMoves(legalMoves)

	return &Node{
		Parent:       parent,
		GameState:    state,
		Move:         move,
		UntriedMoves: movesFromCurrent,
		Children:     []*Node{},
	}
}

// Return true if node is fully expanded otherwise false
func (node *Node) IsFullyExpanded() bool {
	return len(node.UntriedMoves) == 0
}

func (node *Node) IsTerminal() bool {
	// Check if black and white have remaining moves
	// NOTE: CHECK HERE I THINK THE LOGIC MIGHT NOT BE FULLY CORRECT
	return IsTerminalState(node.GameState)
}

// Return the current score of the node. Position 0 is black, position 1 is white
func (node *Node) CurrentScore() [2]int {
	blackScore := node.GameState.Boards.CountOfPieces(true)
	whiteScore := node.GameState.Boards.CountOfPieces(false)
	return [2]int{blackScore, whiteScore}
}

func (node *Node) Winner() WinState {
	return WinnerState(node.GameState)
}
