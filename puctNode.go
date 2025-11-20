package main

// Node struct and methods for PUCT

type PUCTNode struct {
	Q map[uint8]float64 // Average reward for each action
	P map[uint8]float64 // Prior probabilities for each action
	N map[uint8]int     // Visit count for each action (Apparently other implementations write it like this)
	// But i could get the same information of N by iterating over the child and getting the visits
	Visits       int
	Children     []*PUCTNode
	Parent       *PUCTNode
	UntriedMoves []uint8
	Move         uint8 // The move that led us here
	GameState    State // Current boards with whose turn is it to move
}

func InitialRootPUCTNode() *PUCTNode {
	var node *PUCTNode
	var boards Board
	boards.Init()
	var state State
	state.Boards = boards
	state.BlackTurn = true
	var emptyMove uint8
	node = NewPUCTNode(state, nil, emptyMove)
	return node
}

func NextPUCTNodeFromInput(parent *PUCTNode, moveIndex uint8) *PUCTNode {
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
	return NewPUCTNode(newState, parent, moveIndex)
}

func NewPUCTNode(state State, parent *PUCTNode, move uint8) *PUCTNode {
	var legalMoves uint64
	if state.BlackTurn {
		legalMoves = generateMoves(state.Boards.Black, state.Boards.White)
	} else {
		legalMoves = generateMoves(state.Boards.White, state.Boards.Black)
	}
	movesFromCurrent := FastArrayOfMoves(legalMoves)

	priors := make(map[uint8]float64, len(movesFromCurrent))
	uniformPrior := 1.0 / float64(len(movesFromCurrent))
	for _, m := range movesFromCurrent {
		priors[m] = uniformPrior
	}

	return &PUCTNode{
		Parent:       parent,
		GameState:    state,
		Move:         move,
		UntriedMoves: movesFromCurrent,
		Children:     []*PUCTNode{},
		N:            make(map[uint8]int),
		Q:            make(map[uint8]float64),
		P:            priors,
	}
}

// Return true if node is fully expanded otherwise false
func (node *PUCTNode) IsFullyExpandedPUCT() bool {
	return len(node.UntriedMoves) == 0
}

// Check if black and white have remaining moves
func (node *PUCTNode) IsTerminalPUCT() bool {
	return IsTerminalState(node.GameState)
}

// Return the current score of the node. Position 0 is black, position 1 is white
func (node *PUCTNode) CurrentScorePUCT() [2]int {
	blackScore := node.GameState.Boards.CountOfPieces(true)
	whiteScore := node.GameState.Boards.CountOfPieces(false)
	return [2]int{blackScore, whiteScore}
}

func (node *PUCTNode) WinnerPUCT() WinState {
	return WinnerState(node.GameState)
}
