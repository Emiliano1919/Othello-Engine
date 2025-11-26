package main

// Node struct and methods for PUCT

// PUCTNode is a game Node with extra variables needed to implement PUCT version of MCTS.
type PUCTNode struct {
	Q            map[uint8]float64
	P            map[uint8]float64
	N            map[uint8]int
	Parent       *PUCTNode
	Children     []*PUCTNode
	UntriedMoves []uint8
	GameState    State
	Visits       int
	Move         uint8
}

// InitialRootPUCTNode returns the initial root, the start of the game in Node PUCT form.
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

// NextPUCTNodeFromInput returns the root node of the subtree resulting
// from selecting the given move from the current position.
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

// NewPUCTNode returns a new PUCT node.
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

// IsFullyExpandedPUCT returns true if node is fully expanded otherwise false.
// A node is considered fully expanded if no moves are left to try (you cannot generate more children of it).
func (node *PUCTNode) IsFullyExpandedPUCT() bool {
	return len(node.UntriedMoves) == 0
}

// IsTerminalPUCT checks if black and white have remaining moves.
// If the game can be continued after this node then it returns false, otherwise true.
func (node *PUCTNode) IsTerminalPUCT() bool {
	return IsTerminalState(node.GameState)
}

// CurrentScorePUCT returns the current score of the game at this node/state.
// Position 0 is black, position 1 is white
func (node *PUCTNode) CurrentScorePUCT() [2]int {
	blackScore := node.GameState.Boards.CountOfPieces(true)
	whiteScore := node.GameState.Boards.CountOfPieces(false)
	return [2]int{blackScore, whiteScore}
}

// WinnerPUCT returns the current WinState at this node.
// Who is winning.
func (node *PUCTNode) WinnerPUCT() WinState {
	return WinnerState(node.GameState)
}
