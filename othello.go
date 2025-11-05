package main

import (
	"fmt"
	"math"
	"math/rand"
)

type State struct {
	Boards    Board
	BlackTurn bool
}

// Board holds bitboards for black and white disks (disjoint).
type Board struct {
	Black uint64
	White uint64
}

// CellState represents the state of a single cell on the board.
type CellState int

const (
	CELL_BLACK CellState = iota
	CELL_WHITE
	CELL_EMPTY
)

// --- Core Board Methods ---

// CellState returns the state of the cell at (row, col).
func (b *Board) CellState(row, col int) CellState {
	if row < 0 || row > 7 || col < 0 || col > 7 {
		panic("row/col out of range")
	}

	mask := uint64(1) << (row*8 + col)
	// We set a 64bit integer to 1 and then we shift it to its position
	// in the abstract representation of the board. Imagine we start with 000000...1
	// and then we move it to 0000...1.....0 and then if you print this as a board it will
	// be in the calculated position

	switch {
	case b.Black&mask != 0:
		return CELL_BLACK
	case b.White&mask != 0:
		return CELL_WHITE
	default:
		return CELL_EMPTY
	}
}

// Count the number of pieces on the black board
func (b *Board) CountOfPieces(forBlack bool) int {
	var n uint64
	if forBlack {
		n = b.Black
	} else {
		n = b.White
	}
	counter := 0
	for n != 0 {
		n = n & (n - 1)
		counter++
	}
	return counter
}

// SetCellState sets a cell to black, white, or empty.
func (b *Board) SetCellState(row, col int, state CellState) {
	if row < 0 || row > 7 || col < 0 || col > 7 {
		panic("row/col out of range")
	}

	mask := uint64(1) << (row*8 + col)

	// Clear any existing disks at position mask
	b.Black = b.Black &^ mask
	b.White = b.White &^ mask

	// Set the new one
	switch state {
	case CELL_BLACK:
		b.Black = b.Black | mask
	case CELL_WHITE:
		b.White = b.White | mask
	}
}

// Init initializes the board with the standard starting position.
func (b *Board) Init() {
	b.Black = 0
	b.White = 0

	// Starting position:
	// Black on D5 (3,4) and E4 (4,3)
	// White on D4 (3,3) and E5 (4,4)
	b.SetCellState(3, 4, CELL_BLACK)
	b.SetCellState(4, 3, CELL_BLACK)
	b.SetCellState(3, 3, CELL_WHITE)
	b.SetCellState(4, 4, CELL_WHITE)
}

// --- Bit Shifting for Directions ---

const NUM_DIRS = 8

// shift moves all bits in `disks` one step in the given direction.
func shift(disks uint64, dir int) uint64 {
	if dir < 0 || dir >= NUM_DIRS {
		panic("invalid direction")
	}
	// We use the MASKS so that when shifting the bits do not wrap around
	var MASKS = [NUM_DIRS]uint64{
		0x7F7F7F7F7F7F7F7F, // Right
		0x007F7F7F7F7F7F7F, // Down-right
		0xFFFFFFFFFFFFFFFF, // Down
		0x00FEFEFEFEFEFEFE, // Down-left
		0xFEFEFEFEFEFEFEFE, // Left
		0xFEFEFEFEFEFEFE00, // Up-left
		0xFFFFFFFFFFFFFFFF, // Up
		0x7F7F7F7F7F7F7F00, // Up-right
	}

	var LSHIFTS = [NUM_DIRS]uint{
		0, // Right
		0, // Down-right
		0, // Down
		0, // Down-left
		1, // Left
		9, // Up-left
		8, // Up
		7, // Up-right
	}

	var RSHIFTS = [NUM_DIRS]uint{
		1, // Right
		9, // Down-right
		8, // Down
		7, // Down-left
		0, // Left
		0, // Up-left
		0, // Up
		0, // Up-right
	}

	if dir < NUM_DIRS/2 {
		// Shifting right
		return (disks >> RSHIFTS[dir]) & MASKS[dir]
	}
	// Shifting left
	return (disks << LSHIFTS[dir]) & MASKS[dir]
}

/*
	Returns a bitboard where the 1s represent valid places to put the disk

Uses the Dumb7fill algorithm to use the bitboard to get the valid moves
*/
func generateMoves(myDisks, oppDisks uint64) uint64 {
	empty := ^(myDisks | oppDisks)
	var legalMoves uint64

	for dir := 0; dir < NUM_DIRS; dir++ {
		x := shift(myDisks, dir) & oppDisks // x is where there are my disk and next to them the opponents disk
		for i := 0; i < 6; i++ {            // repeat 6 more times to cover up to 7 squares
			// The x in the expression is to keep track of the previous sequences
			x = x | (shift(x, dir) & oppDisks) // Add to x the oponent disks adjacent to those
		}
		legalMoves = legalMoves | (shift(x, dir) & empty) // After all that if you find a white space it is a legal move
	}

	return legalMoves
}

// Get a sorted array of the legal move locations in the board
func ArrayOfMoves(legalMoves uint64) []int {
	var res []int
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			mask := uint64(1) << (row*8 + col)
			if mask&legalMoves != 0 {
				res = append(res, row*8+col)
			}
		}
	}
	return res
}

// Get a sorted array of the legal move locations in row column format in the board
func ArrayOfPositionalMoves(legalMoves []int) [][2]int {
	var res [][2]int
	for _, move := range legalMoves {
		row := move / 8
		col := move % 8
		res = append(res, [2]int{row, col})
	}
	return res
}

// --- Helpers to check if there are possible moves ---
// Input: Flag to indicate if we want to check for black or not (White)
func (b *Board) HasValidMove(forBlack bool) bool {
	if forBlack {
		return generateMoves(b.Black, b.White) != 0
	}
	return generateMoves(b.White, b.Black) != 0
}

/*
	Check if a particular move is valid, return true if yes, and 0 if not

// We achieve this using masks
*/
func (b *Board) IsValidMove(forBlack bool, row, col int) bool {
	mask := uint64(1) << (row*8 + col) // Generate a mask where the 1 is placed where the move is what we want to verify
	if forBlack {
		return generateMoves(b.Black, b.White)&mask != 0
	}
	return generateMoves(b.White, b.Black)&mask != 0
}

/*
Once a move is made we update the board (the sandwhiched disks need to change colors)
moveIndex should be a number that represents a position in the uint64 so it can range from 0 to 63
*/
func resolveMove(myDisks, oppDisks *uint64, moveIndex int) {
	newDisk := uint64(1) << moveIndex
	var captured uint64

	*myDisks |= newDisk
	// Use dumb7fill to find captured/sandwhiched disks
	for dir := 0; dir < NUM_DIRS; dir++ {
		x := shift(newDisk, dir) & *oppDisks // We mark with 1 only the disks sandwhiched
		for i := 0; i < 6; i++ {
			x = x | (shift(x, dir) & *oppDisks)
		}
		boundingDisk := shift(x, dir) & *myDisks
		// If you found captured disks then
		if boundingDisk != 0 {
			captured = captured | x
		}
	}

	*myDisks ^= captured  // We add the captured ones
	*oppDisks ^= captured // We substract the captured ones
}

/*
Given a position and a board execute the move.
*/

func (b *Board) MakeMove(forBlack bool, row, col int) {
	if !b.IsValidMove(forBlack, row, col) {
		panic("invalid move")
	}

	moveIndex := row*8 + col
	if forBlack {
		resolveMove(&b.Black, &b.White, moveIndex)
	} else {
		resolveMove(&b.White, &b.Black, moveIndex)
	}
}

// --- Debug / Display Helpers ---

// PrintBoard prints the board in a readable 8×8 grid.
func (b *Board) PrintBoard() {
	fmt.Println()
	fmt.Println("    a b c d e f g h")
	fmt.Println("   -----------------")
	for row := 0; row < 8; row++ {
		fmt.Printf("%d | ", row+1)
		for col := 0; col < 8; col++ {
			switch b.CellState(row, col) {
			case CELL_BLACK:
				fmt.Print("@ ")
			case CELL_WHITE:
				fmt.Print("o ")
			default:
				fmt.Print(". ")
			}
		}
		fmt.Printf("| %d\n", row+1)
	}
	fmt.Println("   -----------------")
	fmt.Println("    a b c d e f g h")
	fmt.Println()
}

// PrintBoardWithMoves prints the board and shows all possible legal moves
// for the current player, marking them with '*'. It also prints a list of
// usable move inputs like "d3".
func (s *State) PrintBoardWithMoves() {
	b := s.Boards

	var legalMoves uint64
	if s.BlackTurn {
		legalMoves = generateMoves(b.Black, b.White)
	} else {
		legalMoves = generateMoves(b.White, b.Black)
	}

	fmt.Println()
	if s.BlackTurn {
		fmt.Println("Turn: Black (B)")
	} else {
		fmt.Println("Turn: White (W)")
	}

	fmt.Println("    a b c d e f g h")
	fmt.Println("   -----------------")
	for row := 0; row < 8; row++ {
		fmt.Printf("%d | ", row+1)
		for col := 0; col < 8; col++ {
			mask := uint64(1) << (row*8 + col)
			switch {
			case b.Black&mask != 0:
				fmt.Print("@ ")
			case b.White&mask != 0:
				fmt.Print("o ")
			case legalMoves&mask != 0:
				fmt.Print("* ") // possible move
			default:
				fmt.Print(". ")
			}
		}
		fmt.Printf("| %d\n", row+1)
	}
	fmt.Println("   -----------------")
	fmt.Println("    a b c d e f g h")

	// Now print the moves in input format
	arr := ArrayOfPositionalMoves(ArrayOfMoves(legalMoves))
	if len(arr) == 0 {
		fmt.Println("\nNo legal moves available.")
		return
	}

	fmt.Println("\nPossible moves:")
	for _, m := range arr {
		row := m[0]
		col := m[1]
		fmt.Printf("  %c%d -> enter as: %c%d or \"%d %d\"\n",
			'a'+rune(col), row+1, 'a'+rune(col), row+1, row, col)
	}
	fmt.Println()
}

func RequestMove() [2]int {
	var arr [2]int

	fmt.Print("Enter your move white (e.g., 3 7): ")
	_, err := fmt.Scanf("%d %d", &arr[0], &arr[1])
	if err != nil {
		fmt.Println("Error:", err)
		panic(nil)
	}
	fmt.Println("You entered:", arr)
	return arr
}

func RequestUserIsBlack() bool {
	for {
		var choice string
		fmt.Print("Enter B if you want to play as Black or W if you want to play as White: ")
		_, err := fmt.Scanln(&choice)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		if len(choice) == 0 {
			fmt.Println("No input received, defaulting to White.")
			return false
		}

		switch choice[0] {
		case 'B', 'b':
			fmt.Println("You will play as Black.")
			return true
		case 'W', 'w':
			fmt.Println("You will play as White.")
			return false
		default:
			fmt.Println("Invalid input, please enter B or W.")
		}
	}
}

func OutputResult(node *Node) {
	if node.IsTerminal() {
		fmt.Println("Game finished:")
	}
	black := node.GameState.Boards.CountOfPieces(true)
	white := node.GameState.Boards.CountOfPieces(false)
	fmt.Printf("Black: %d\n", black)
	fmt.Printf("White: %d\n", white)

	if black > white {
		fmt.Println("Winner: Black")
	} else if white > black {
		fmt.Println("Winner: White")
	} else {
		fmt.Println("It's a tie!")
	}
}

func PrintBitboard(bits uint64) {
	fmt.Println("  A B C D E F G H")
	for row := 0; row < 8; row++ {
		fmt.Printf("%d ", row+1)
		for col := 0; col < 8; col++ {
			mask := uint64(1) << (row*8 + col)
			if bits&mask != 0 {
				fmt.Print("x ")
			} else {
				fmt.Print(". ")
			}
		}
		fmt.Println()
	}
}

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
		fmt.Println("No valid moves for next player — passing turn back.")
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

func IsTerminalState(state State) bool {
	return state.Boards.HasValidMove(true) == false && state.Boards.HasValidMove(false) == false
}

func (node *Node) IsTerminal() bool {
	// Check if black and white have remaining moves
	// NOTE: CHECK HERE I THINK THE LOGIC MIGHT NOT BE FULLY CORRECT
	return IsTerminalState(node.GameState)
}

func WinnerState(state State) int {
	if state.Boards.CountOfPieces(true) > state.Boards.CountOfPieces(false) {
		return 1 // Black wins
	} else if state.Boards.CountOfPieces(true) < state.Boards.CountOfPieces(false) {
		return 0 // White wins
	} else {
		return 2 // Draw
	}
}

func (node *Node) Winner() int {
	return WinnerState(node.GameState)
}

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

func Traverse(node *Node) *Node {
	for node.IsFullyExpanded() && !node.IsTerminal() {
		node = BestUCT(node, float64(2))
	}
	if node.IsTerminal() {
		return node
	}
	return node.Expand()
}

func BestUCT(node *Node, c float64) *Node {
	var best *Node
	bestUCT := float64(-1 << 63)
	for _, child := range node.Children {
		explotationTerm := float64(child.Wins) / float64(child.Visits)
		explorationTerm := math.Sqrt(math.Log(float64(node.Visits)) / float64(child.Visits))
		C := math.Sqrt(c) // Theoretical value, will try to find a better one through self play
		UCTValue := explorationTerm + C*explotationTerm

		if UCTValue > bestUCT {
			bestUCT = UCTValue
			best = child
		}
	}
	return best
}

func SimulateRollout(state State) int {
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
		move := moveArray[rand.Intn(len(moveArray))] // Here is the rollout ppolicy  which is random

		current.Boards.MakeMove(current.BlackTurn, move[0], move[1])
		current.BlackTurn = !current.BlackTurn
	}
	// 1 = Black win, 0 = White win, 2 = draw
	return WinnerState(current)
}

func backpropagate(node *Node, result int) {
	for node != nil {
		node.Visits++
		// Will have to fix this here,a dn take into account the turn correctly
		if result == 2 {
			node.Wins += 1
		} else {
			node.Wins += result // We are technically speaking just accounting for black
		}
		node = node.Parent
	}
}

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

func MonteCarloTreeSearch(currentRoot *Node, iterations int) *Node {
	if currentRoot.IsTerminal() {
		return currentRoot
	}
	for i := 0; i < iterations; i++ {
		leaf := Traverse(currentRoot)
		var nodeToSimulateFrom *Node
		if len(leaf.UntriedMoves) > 0 {
			child := leaf.Expand()
			nodeToSimulateFrom = child
		} else {
			nodeToSimulateFrom = leaf
		}

		result := SimulateRollout(nodeToSimulateFrom.GameState)
		backpropagate(nodeToSimulateFrom, result)
	}
	return BestNodeFromMCTS(currentRoot)
}

// --- Example usage ---

func main() {
	initialNode := InitialRootNode()
	initialNode.GameState.Boards.PrintBoard()
	bestOpening := MonteCarloTreeSearch(initialNode, 5000)
	bestOpening.GameState.Boards.PrintBoard()
	node := bestOpening
	for !node.IsTerminal() {
		if !node.GameState.BlackTurn {
			node.GameState.PrintBoardWithMoves()
			whiteMove := RequestMove()
			node = NextNodeFromInput(node, whiteMove)
		} else {
			mctsNode := MonteCarloTreeSearch(node, 5000)
			mctsNode.GameState.Boards.PrintBoard()
			node = mctsNode
		}
	}
	if node.IsTerminal() {
		OutputResult(node)
	}
}
