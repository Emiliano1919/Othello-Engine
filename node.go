package main

// Node struct and methods

type Node struct {
	Visits       int
	Wins         int
	Parent       *Node
	Children     []*Node
	GameState    State  // Current boards with whose turn is it to move
	Move         [2]int // The move that led us here
	UntriedMoves [][2]int
}

func InitialRootNode() *Node {
	var node *Node
	var boards Board
	boards.Init()
	var state State
	state.Boards = boards
	state.BlackTurn = true
	var empty [2]int
	node = NewNode(state, nil, empty)
	return node
}

func NextNodeFromInput(parent *Node, move [2]int) *Node {
	newBoards := parent.GameState.Boards                             // Copy
	newBoards.MakeMove(parent.GameState.BlackTurn, move[0], move[1]) // We make the move on the current player
	newState := State{
		Boards:    newBoards,
		BlackTurn: !parent.GameState.BlackTurn, // switch turn immediately
	}
	// Switch turn but conditionally to deal with edgecase of the opponent having no moves.
	if !newState.Boards.HasValidMove(newState.BlackTurn) && newState.Boards.HasValidMove(!newState.BlackTurn) {
		//fmt.Println("No valid moves for next player — passing turn back.")
		newState.BlackTurn = !newState.BlackTurn
	}
	return NewNode(newState, parent, move)
}

func NewNode(state State, parent *Node, move [2]int) *Node {
	var legalMoves uint64
	if state.BlackTurn {
		legalMoves = generateMoves(state.Boards.Black, state.Boards.White)
	} else {
		legalMoves = generateMoves(state.Boards.White, state.Boards.Black)
	}
	movesFromCurrent := ArrayOfPositionalMoves(ArrayOfMoves(legalMoves))

	return &Node{
		Parent:       parent,
		GameState:    state,
		Move:         move,
		UntriedMoves: movesFromCurrent,
		Children:     []*Node{},
	}
}

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
