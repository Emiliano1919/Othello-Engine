package main

import "math/bits"

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

// New Winstate to avoid incorrect comparisons and modifications
type WinState int

const (
	WHITE_WIN WinState = iota
	BLACK_WIN
	DRAW
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

// Count the number of pieces on any board
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

// Cardinal directions you can shift to (horizontal, vertical, and diagonals)
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

func FastArrayOfMoves(legal uint64) []uint8 {
	// According to a paper I found the maximum size ever is 33
	res := make([]uint8, 0, 33)

	for m := legal; m != 0; m &= m - 1 { // Use Kernighan algorithm O(log n)
		index := bits.TrailingZeros64(m)
		res = append(res, uint8(index))
	}
	return res
}

// Get a sorted array of the legal move locations in the board
func ArrayOfMoves(legalMoves uint64) []uint8 {
	var res []uint8
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			mask := uint64(1) << (row*8 + col)
			if mask&legalMoves != 0 {
				res = append(res, uint8(row*8+col))
			}
		}
	}
	return res
}

// Get a sorted array of the legal move locations in row column format in the board
func ArrayOfPositionalMoves(legalMoves []uint8) [][2]uint8 {
	res := make([][2]uint8, 0, len(legalMoves))
	for _, move := range legalMoves {
		row := uint8(move >> 3) // Faster division by 8
		col := uint8(move & 7)  // Faster modulo 8
		res = append(res, [2]uint8{row, col})
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
func (b *Board) IsValidMovePositional(forBlack bool, row, col uint8) bool {
	// Check valid board coordinates
	if row >= 8 || col >= 8 {
		return false
	}

	index := row*8 + col // Safe because row,col < 8 → index ∈ [0,63]
	mask := uint64(1) << index

	if forBlack {
		return generateMoves(b.Black, b.White)&mask != 0
	}
	return generateMoves(b.White, b.Black)&mask != 0
}

/*
	Check if a particular move is valid, return true if yes, and 0 if not

// We achieve this using masks
*/
func (b *Board) IsValidMoveIndex(forBlack bool, index uint8) bool {
	if index >= 64 {
		return false
	}
	mask := uint64(1) << index // Generate a mask where the 1 is placed where the move is what we want to verify
	if forBlack {
		return generateMoves(b.Black, b.White)&mask != 0
	}
	return generateMoves(b.White, b.Black)&mask != 0
}

/*
Once a move is made we update the board (the sandwhiched disks need to change colors)
moveIndex should be a number that represents a position in the uint64 so it can range from 0 to 63
*/
func resolveMove(myDisks, oppDisks *uint64, moveIndex uint8) {
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
func (b *Board) MakeMovePositional(forBlack bool, row, col uint8) {
	if !b.IsValidMovePositional(forBlack, row, col) {
		panic("invalid move")
	}

	moveIndex := row*8 + col
	if forBlack {
		resolveMove(&b.Black, &b.White, moveIndex)
	} else {
		resolveMove(&b.White, &b.Black, moveIndex)
	}
}

/*
Given a position and a board execute the move.
*/
func (b *Board) MakeMoveIndex(forBlack bool, index uint8) {
	if !b.IsValidMoveIndex(forBlack, index) {
		panic("invalid move")
	}

	if forBlack {
		resolveMove(&b.Black, &b.White, index)
	} else {
		resolveMove(&b.White, &b.Black, index)
	}
}

// If at least one  (black or white) has a possible move to make then it is a non terminal state
// Only if both have exhausted their moves will it be false
func IsTerminalState(state State) bool {
	return !state.Boards.HasValidMove(true) && !state.Boards.HasValidMove(false)
}

// There are only 3 states Black win, white win and draw
func WinnerState(state State) WinState {
	if state.Boards.CountOfPieces(true) > state.Boards.CountOfPieces(false) {
		return BLACK_WIN // Black wins
	} else if state.Boards.CountOfPieces(true) < state.Boards.CountOfPieces(false) {
		return WHITE_WIN // White wins
	} else {
		return DRAW // Draw
	}
}

// Return the current score of the node. Position 0 is black, position 1 is white
func CurrentStateScore(state State) [2]int {
	blackScore := state.Boards.CountOfPieces(true)
	whiteScore := state.Boards.CountOfPieces(false)
	return [2]int{blackScore, whiteScore}
}
