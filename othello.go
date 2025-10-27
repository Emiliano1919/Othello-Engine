package main

import (
	"fmt"
)

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

// --- Helpers to check moves ---
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
	fmt.Println("  A B C D E F G H")
	for row := 0; row < 8; row++ {
		fmt.Printf("%d ", row+1)
		for col := 0; col < 8; col++ {
			switch b.CellState(row, col) {
			case CELL_BLACK:
				fmt.Print("B ")
			case CELL_WHITE:
				fmt.Print("W ")
			default:
				fmt.Print(". ")
			}
		}
		fmt.Println()
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

// --- Example usage ---

func main() {
	var board Board
	board.Init()
	board.PrintBoard()
	possibles := generateMoves(board.Black, board.White)
	PrintBitboard(possibles)
	fmt.Println(ArrayOfMoves(possibles))
	board.MakeMove(true, 2, 3)
	board.PrintBoard()
	x := board.CountOfPieces(true)
	fmt.Println(x)
	y := board.CountOfPieces(false)
	fmt.Println(y)
}
