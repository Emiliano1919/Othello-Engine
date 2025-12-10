package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"runtime"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	ScreenWidth  = 420
	ScreenHeight = 600
	boardSize    = 8 // 8 by 8 is 64 as the othello board
)

type GamePhase int

const (
	StateStartScreen GamePhase = iota
	StatePlaying
	StateEndScreen
)

type Game struct {
	node           *PUCTNode
	boardImage     *ebiten.Image
	rng            *rand.Rand
	legalMoves     uint64 // We put it here because we calculate it at the end of the machine turn
	waitingForUser bool   // We use this to keep the asynchronous code out of the loop
	userIsBlack    bool
	state          GamePhase
}

func (g *Game) UpdateStartScreen() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if x > 50 && x < 200 && y > 200 && y < 250 { // Black button
			initialNode := InitialRootPUCTNode()
			g.node = initialNode
			g.userIsBlack = true
			g.state = StatePlaying
			g.waitingForUser = true // Black moves first
			// Calculate black's legal moves at start
			g.legalMoves = generateMoves(
				g.node.GameState.Boards.Black,
				g.node.GameState.Boards.White,
			)
		} else if x > 50 && x < 200 && y > 300 && y < 350 { // White button
			initialNode := InitialRootPUCTNode()
			g.node = initialNode
			g.userIsBlack = false
			g.state = StatePlaying
			g.waitingForUser = false
		}
	}
}

func (g *Game) UpdateEndScreen() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if x > 50 && x < 200 && y > 200 && y < 250 { // Restart as Black
			initialNode := InitialRootPUCTNode()
			g.node = initialNode
			g.userIsBlack = true
			g.state = StatePlaying
			g.waitingForUser = true // Black moves first
			// Calculate black's legal moves at start
			g.legalMoves = generateMoves(
				g.node.GameState.Boards.Black,
				g.node.GameState.Boards.White,
			)
		} else if x > 50 && x < 200 && y > 300 && y < 350 { // Restart as White
			initialNode := InitialRootPUCTNode()
			g.node = initialNode
			g.userIsBlack = false
			g.state = StatePlaying
			g.waitingForUser = false
		}
	}
}

func (g *Game) Update() error {
	switch g.state {
	case StateStartScreen:
		g.UpdateStartScreen()
	case StatePlaying:
		// original Update logic here
		if g.waitingForUser {
			// human turn
			if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
				x, y := ebiten.CursorPosition()
				col := uint8(x / (tileSize + tileMargin))
				row := uint8(y / (tileSize + tileMargin))
				maskIndex := row*8 + col
				mask := uint64(1) << (uint64(maskIndex))

				if g.legalMoves&mask != 0 {
					g.node = NextPUCTNodeFromInput(g.node, maskIndex)
					g.waitingForUser = false
				}
			}
		} else {
			if !g.userIsBlack {
				if g.node.GameState.BlackTurn {
					g.node = SingleRunParallelizationMCTSPUCT(g.node, 15000, g.rng)
				}
				// Calculate the possible moves of the opponent if you pass the turn to them
				if !g.node.GameState.BlackTurn {
					g.legalMoves = generateMoves(g.node.GameState.Boards.White, g.node.GameState.Boards.Black)
					g.waitingForUser = true
				}
			} else {
				if !g.node.GameState.BlackTurn {
					g.node = SingleRunParallelizationMCTSPUCT(g.node, 15000, g.rng)
				}
				// Calculate the possible moves of the opponent if you pass the turn to them
				if g.node.GameState.BlackTurn {
					g.legalMoves = generateMoves(g.node.GameState.Boards.Black, g.node.GameState.Boards.White)
					g.waitingForUser = true
				}
			}

		}

		// Check if game is over
		if g.node.IsTerminalPUCT() {
			g.state = StateEndScreen
		}

	case StateEndScreen:
		g.UpdateEndScreen()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case StateStartScreen:
		screen.Fill(color.RGBA{30, 30, 30, 255})
		ebitenutil.DebugPrintAt(screen, "Play as Black", 50, 200)
		ebitenutil.DebugPrintAt(screen, "Play as White", 50, 300)

	case StatePlaying:
		if g.boardImage == nil {
			size := boardSize*tileSize + (boardSize+1)*tileMargin
			g.boardImage = ebiten.NewImage(size, size)
		}
		g.node.GameState.Draw(g.boardImage)
		if g.waitingForUser {
			for i := 0; i < 64; i++ {
				mask := uint64(1) << i
				if g.legalMoves&mask != 0 {
					if !g.userIsBlack {
						drawDiskAtIndex(g.boardImage, i, whitePossibleDiskImg)
					} else {
						drawDiskAtIndex(g.boardImage, i, blackPossibleDiskImg)
					}
				}
			}
		}

		screen.Fill(color.RGBA{0, 0, 0, 255})
		op := &ebiten.DrawImageOptions{}
		screen.DrawImage(g.boardImage, op)

	case StateEndScreen:
		screen.Fill(color.RGBA{20, 20, 20, 255})
		score := CurrentStateScore(g.node.GameState)
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Game Over!\nBlack: %d\nWhite: %d", score[0], score[1]), 50, 100)
		ebitenutil.DebugPrintAt(screen, "Restart as Black", 50, 200)
		ebitenutil.DebugPrintAt(screen, "Restart as White", 50, 300)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	size := boardSize*tileSize + (boardSize+1)*tileMargin
	return size, size
}

//Visual Main

// func main() {
// 	ebiten.SetWindowTitle("Othello Engine (Ebiten Board)")
// 	size := boardSize*tileSize + (boardSize+1)*tileMargin
// 	ebiten.SetWindowSize(size, size)
// 	// Create a new RNG
// 	start := time.Now()
// 	rng := rand.New(rand.NewSource(start.UnixNano()))

// 	// Pass rng to the Game instance
// 	game := &Game{
// 		rng: rng,
// 	}
// 	if err := ebiten.RunGame(game); err != nil {
// 		log.Fatal(err)
// 	}
// }

func Versus() {
	start := time.Now()
	rng := rand.New(rand.NewSource(start.UnixNano()))
	// Each AI needs its own tree, so that they do not share knowledge and influence the other
	// But they will update each other of their respective moves
	initialNodeP1 := InitialRootNode()
	initialNodeP2 := InitialRootNode()
	OpponentIsBlack := false // Is the opponent of baseline black?
	var nodeP1 *Node
	var nodeP2 *Node
	if !OpponentIsBlack {
		nodeP1 = OriginalMonteCarloTreeSearch(initialNodeP1, 500, rng)
		nodeP2 = NextNodeFromInput(initialNodeP2, nodeP1.Move)
		//nodeP2.GameState.Boards.PrintBoard()
	} else {
		//nodeP2 = InnacurateMonteCarloTreeSearch(initialNodeP2, 500, OPTIMIZE_FOR_BLACK)
		nodeP2 = SingleRunParallelizationMCTS(initialNodeP2, 50, rng)
		nodeP1 = NextNodeFromInput(initialNodeP1, nodeP2.Move)
		// nodeP1.GameState.Boards.PrintBoard()
	}
	for !nodeP1.IsTerminal() { // This works because both nodes update each other
		if !OpponentIsBlack {
			if !nodeP1.GameState.BlackTurn {
				//nodeP2 = InnacurateMonteCarloTreeSearch(nodeP2, 500, OPTIMIZE_FOR_WHITE)
				nodeP2 = SingleRunParallelizationMCTS(nodeP2, 50, rng)
				nodeP1 = NextNodeFromInput(nodeP1, nodeP2.Move)
				//nodeP1.GameState.Boards.PrintBoard()
			} else {
				nodeP1 = OriginalMonteCarloTreeSearch(nodeP1, 500, rng)
				nodeP2 = NextNodeFromInput(nodeP2, nodeP1.Move)
				//nodeP2.GameState.Boards.PrintBoard()
			}
		} else {
			if nodeP1.GameState.BlackTurn {
				//nodeP2 = InnacurateMonteCarloTreeSearch(nodeP2, 500, OPTIMIZE_FOR_BLACK)
				nodeP2 = SingleRunParallelizationMCTS(nodeP2, 50, rng)
				nodeP1 = NextNodeFromInput(nodeP1, nodeP2.Move)
				//nodeP1.GameState.Boards.PrintBoard()
			} else {
				nodeP1 = OriginalMonteCarloTreeSearch(nodeP1, 500, rng)
				nodeP2 = NextNodeFromInput(nodeP2, nodeP1.Move)
				//nodeP2.GameState.Boards.PrintBoard()
			}
		}
	}
}

// Versus main
func main() {
	fmt.Println("GOMAXPROCS:", runtime.GOMAXPROCS(0))
	start := time.Now()
	seed := start.UnixNano()
	rng1 := rand.New(rand.NewSource(seed))
	rng2 := rand.New(rand.NewSource(seed + 1))
	// Just one is needed because access is sequential between the 2 AIs  (they dont move at the same time)
	// And the parallelization already has the necessary protections for sharing rng (each goroutine has its own rng)
	OpponentWinCounter := 0
	DrawsCounter := 0
	Games := 100
	for i := 0; i < Games; i++ {
		// Each AI needs its own tree, so that they do not share knowledge and influence the other
		// But they will update each other of their respective moves
		initialNodeP1 := InitialRootNode()
		initialNodeP2 := InitialRootPUCTNode()
		OpponentIsBlack := false // Is the opponent of baseline black?
		var nodeP1 *Node
		var nodeP2 *PUCTNode
		if !OpponentIsBlack {
			nodeP1 = OriginalMonteCarloTreeSearch(initialNodeP1, 500, rng1)
			nodeP2 = NextPUCTNodeFromInput(initialNodeP2, nodeP1.Move)
			//nodeP2.GameState.Boards.PrintBoard()
		} else {
			nodeP2 = SingleRunParallelizationMCTSPUCT(initialNodeP2, 200, rng2)
			nodeP1 = NextNodeFromInput(initialNodeP1, nodeP2.Move)
			// nodeP1.GameState.Boards.PrintBoard()
		}
		for !nodeP1.IsTerminal() { // This works because both nodes update each other
			if !OpponentIsBlack {
				if !nodeP1.GameState.BlackTurn {
					nodeP2 = SingleRunParallelizationMCTSPUCT(nodeP2, 200, rng2)
					nodeP1 = NextNodeFromInput(nodeP1, nodeP2.Move)
					//nodeP1.GameState.Boards.PrintBoard()
				} else {
					nodeP1 = OriginalMonteCarloTreeSearch(nodeP1, 500, rng1)
					nodeP2 = NextPUCTNodeFromInput(nodeP2, nodeP1.Move)
					//nodeP2.GameState.Boards.PrintBoard()
				}
			} else {
				if nodeP1.GameState.BlackTurn {
					nodeP2 = SingleRunParallelizationMCTSPUCT(nodeP2, 200, rng2)
					nodeP1 = NextNodeFromInput(nodeP1, nodeP2.Move)
					//nodeP1.GameState.Boards.PrintBoard()
				} else {
					nodeP1 = OriginalMonteCarloTreeSearch(nodeP1, 500, rng1)
					nodeP2 = NextPUCTNodeFromInput(nodeP2, nodeP1.Move)
					//nodeP2.GameState.Boards.PrintBoard()
				}
			}
		}
		if nodeP1.IsTerminal() {
			//OutputResult(nodeP1)
			//OutputResult(nodeP2)
			switch nodeP1.Winner() {
			case BLACK_WIN:
				if OpponentIsBlack {
					OpponentWinCounter++
				}
			case WHITE_WIN:
				if !OpponentIsBlack {
					OpponentWinCounter++
				}
			case DRAW:
				DrawsCounter++
			}
		}
		fmt.Printf("Number of finalized games: %d\n", i+1)
	}
	elapsed := time.Since(start)
	fmt.Printf("Opponent Wins: %d\n", OpponentWinCounter)
	fmt.Printf("Draws: %d\n", DrawsCounter)
	fmt.Printf("Total Games ran: %d\n", Games)
	fmt.Printf("Total run time for all the games: %s", elapsed)

}

// Debugging Main
// func main() {
// 	start := time.Now()
// 	rng := rand.New(rand.NewSource(start.UnixNano()))
// 	initialNode := InitialRootNode()
// 	initialNode.GameState.Boards.PrintBoard()
// 	userIsBlack := RequestUserIsBlack()
// 	var node *Node
// 	if !userIsBlack {
// 		bestOpening := SingleRunParallelizationMCTS(initialNode, 5000, rng)
// 		bestOpening.GameState.Boards.PrintBoard()
// 		node = bestOpening
// 	} else {
// 		initialNode.GameState.PrintBoardWithMoves()
// 		blackMove := RequestMove(userIsBlack)
// 		node = NextNodeFromInput(initialNode, blackMove)
// 	}
// 	for !node.IsTerminal() {
// 		if !userIsBlack {
// 			if !node.GameState.BlackTurn {
// 				node.GameState.PrintBoardWithMoves()
// 				whiteMove := RequestMove(userIsBlack)
// 				node = NextNodeFromInput(node, whiteMove)
// 			} else {
// 				mctsNode := SingleRunParallelizationMCTS(node, 5000, rng)
// 				mctsNode.GameState.Boards.PrintBoard()
// 				node = mctsNode
// 			}
// 		} else {
// 			if node.GameState.BlackTurn {
// 				node.GameState.PrintBoardWithMoves()
// 				blackMove := RequestMove(userIsBlack)
// 				node = NextNodeFromInput(node, blackMove)
// 			} else {
// 				mctsNode := SingleRunParallelizationMCTS(node, 5000, rng)
// 				mctsNode.GameState.Boards.PrintBoard()
// 				node = mctsNode
// 			}
// 		}
// 	}
// 	if node.IsTerminal() {
// 		OutputResult(node)
// 	}
// }

// func main() {
// 	// Just a quick test to verify the output is the same
// 	initialNode := InitialRootNode()

// 	legal := generateMoves(initialNode.GameState.Boards.Black, initialNode.GameState.Boards.White)
// 	PrintBitboard(legal)

// 	v := (uint64(1) << 0) |
// 		(uint64(1) << 2) |
// 		(uint64(1) << 44) |
// 		(uint64(1) << 63)
// 	moves := FastArrayOfMoves(v)
// 	otherMoves := ArrayOfMoves(v)
// 	otherMovesPos := ArrayOfPositionalMoves(otherMoves)
// 	PrintBitboard(v)
// 	fmt.Println("ArrayOfMoves output:", otherMoves)
// 	fmt.Println("ArrayOfPositionalMoves output:", otherMovesPos)
// 	fmt.Println("FastArrayOfMoves output:", moves)

// 	fmt.Println("Moves (row, col):")
// 	for _, m := range moves {
// 		row := m >> 3
// 		col := m & 7
// 		fmt.Printf("  (%d, %d)\n", row, col)
// 	}
// 	initialNode.GameState.Boards.PrintBoard()
// 	initialNode.GameState.PrintBoardWithMoves()
// }
