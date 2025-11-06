package main

import (
	"fmt"
	"image/color"
	"log"

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
	node           *Node
	userIsBlack    bool
	boardImage     *ebiten.Image
	waitingForUser bool   // We use this to keep the asynchronous code out of the loop
	legalMoves     uint64 // We put it here because we calculate it at the end of the machine turn
	state          GamePhase
}

func NewGame() *Game {
	initialNode := InitialRootNode()
	return &Game{
		node:  initialNode,
		state: StateStartScreen, // Start with the start screen
		// userIsBlack is not set yet; will be set by user's choice
	}
}

func (g *Game) UpdateStartScreen() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if x > 50 && x < 200 && y > 200 && y < 250 { // Black button
			initialNode := InitialRootNode()
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
			initialNode := InitialRootNode()
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
			*g = *NewGame()
			g.userIsBlack = true
			g.waitingForUser = true
			g.state = StatePlaying
		} else if x > 50 && x < 200 && y > 300 && y < 350 { // Restart as White
			*g = *NewGame()
			g.userIsBlack = false
			g.waitingForUser = false
			g.state = StatePlaying
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
				col := x / (tileSize + tileMargin)
				row := y / (tileSize + tileMargin)
				mask := uint64(1) << (row*8 + col)

				if g.legalMoves&mask != 0 {
					g.node = NextNodeFromInput(g.node, [2]int{row, col})
					g.waitingForUser = false
				}
			}
		} else {
			if !g.userIsBlack {
				if g.node.GameState.BlackTurn {
					g.node = MonteCarloTreeSearch(g.node, 5000)
				}
				// Calculate the possible moves of the opponent if you pass the turn to them
				if !g.node.GameState.BlackTurn {
					g.legalMoves = generateMoves(g.node.GameState.Boards.White, g.node.GameState.Boards.Black)
					g.waitingForUser = true
				}
			} else {
				if !g.node.GameState.BlackTurn {
					g.node = MonteCarloTreeSearch(g.node, 5000)
				}
				// Calculate the possible moves of the opponent if you pass the turn to them
				if g.node.GameState.BlackTurn {
					g.legalMoves = generateMoves(g.node.GameState.Boards.Black, g.node.GameState.Boards.White)
					g.waitingForUser = true
				}
			}

		}

		// Check if game is over
		if g.node.IsTerminal() {
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
		score := g.node.CurrentScore()
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Game Over!\nBlack: %d\nWhite: %d", score[0], score[1]), 50, 100)
		ebitenutil.DebugPrintAt(screen, "Restart as Black", 50, 200)
		ebitenutil.DebugPrintAt(screen, "Restart as White", 50, 300)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	size := boardSize*tileSize + (boardSize+1)*tileMargin
	return size, size
}

func main() {
	ebiten.SetWindowTitle("Othello Engine (Ebiten Board)")
	size := boardSize*tileSize + (boardSize+1)*tileMargin
	ebiten.SetWindowSize(size, size)
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}

// // Game loop
// func main() {
// 	initialNode := InitialRootNode()
// 	initialNode.GameState.Boards.PrintBoard()
// 	userIsBlack := RequestUserIsBlack()
// 	var node *Node
// 	if !userIsBlack {
// 		bestOpening := MonteCarloTreeSearch(initialNode, 5000)
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
// 				mctsNode := MonteCarloTreeSearch(node, 5000)
// 				mctsNode.GameState.Boards.PrintBoard()
// 				node = mctsNode
// 			}
// 		} else {
// 			if node.GameState.BlackTurn {
// 				node.GameState.PrintBoardWithMoves()
// 				blackMove := RequestMove(userIsBlack)
// 				node = NextNodeFromInput(node, blackMove)
// 			} else {
// 				mctsNode := MonteCarloTreeSearch(node, 5000)
// 				mctsNode.GameState.Boards.PrintBoard()
// 				node = mctsNode
// 			}
// 		}
// 	}
// 	if node.IsTerminal() {
// 		OutputResult(node)
// 	}
// }
