package main

import "fmt"

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

func RequestMove(userIsBlack bool) [2]int {
	var arr [2]int

	color := "white"
	if userIsBlack {
		color = "black"
	}

	fmt.Printf("Enter your move %s (e.g., 6 7): ", color)
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
