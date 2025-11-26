package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	tileSize   = 50
	tileMargin = 2
)

var (
	frameColor = color.RGBA{0, 100, 0, 255}
	emptyColor = color.RGBA{0, 130, 0, 255}

	blackDiskColor     = color.RGBA{0, 0, 0, 255}
	whiteDiskColor     = color.RGBA{255, 255, 255, 255}
	blackPossibleColor = color.RGBA{10, 10, 10, 150}
	whitePossibleColor = color.RGBA{240, 240, 240, 150}

	boardImage           *ebiten.Image
	blackDiskImg         *ebiten.Image
	whiteDiskImg         *ebiten.Image
	blackPossibleDiskImg *ebiten.Image
	whitePossibleDiskImg *ebiten.Image
)

func init() {
	// Static pre-rendered disk
	boardImage = makeBoardImage()
	blackDiskImg = makeDiskImage(blackDiskColor)
	whiteDiskImg = makeDiskImage(whiteDiskColor)
	blackPossibleDiskImg = makeDiskImage(blackPossibleColor)
	whitePossibleDiskImg = makeDiskImage(whitePossibleColor)
}

func makeBoardImage() *ebiten.Image {
	img := ebiten.NewImage(boardSize*(tileSize+tileMargin)+tileMargin, boardSize*(tileSize+tileMargin)+tileMargin)
	img.Fill(frameColor)

	for row := 0; row < boardSize; row++ {
		for col := 0; col < boardSize; col++ {
			x := col*tileSize + (col+1)*tileMargin
			y := row*tileSize + (row+1)*tileMargin
			tile := ebiten.NewImage(tileSize, tileSize)
			tile.Fill(emptyColor)
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x), float64(y))
			img.DrawImage(tile, op)
		}
	}
	return img
}

// makeDiskImage pre renders the disks to avoid re rendering.
func makeDiskImage(c color.RGBA) *ebiten.Image {
	img := ebiten.NewImage(tileSize, tileSize)
	img.Fill(color.Transparent) // fully transparent background

	pixels := make([]byte, tileSize*tileSize*4) // We set the pixels directly for efficiency

	cx, cy := float64(tileSize)/2, float64(tileSize)/2
	radius := float64(tileSize)/2 - 2
	a := float64(c.A) / 255.0

	for y := 0; y < tileSize; y++ {
		for x := 0; x < tileSize; x++ {
			dx := float64(x) - cx
			dy := float64(y) - cy
			dist := dx*dx + dy*dy
			idx := (y*tileSize + x) * 4
			if dist <= radius*radius {
				// The pixels need to be multiplied by alpha(transparency for possible values)
				pixels[idx+0] = uint8(float64(c.R) * a)
				pixels[idx+1] = uint8(float64(c.G) * a)
				pixels[idx+2] = uint8(float64(c.B) * a)
				pixels[idx+3] = c.A
			} else {
				pixels[idx+0] = 0
				pixels[idx+1] = 0
				pixels[idx+2] = 0
				pixels[idx+3] = 0
			}
		}
	}

	img.WritePixels(pixels)
	return img
}

// Draw places the prerendered disks on the image of the board (also pre rendered).
func (s *State) Draw(screen *ebiten.Image) {
	// Draw the static board background
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(boardImage, op)

	// Then draw only the disks that are set in your bitboards
	for row := 0; row < boardSize; row++ {
		for col := 0; col < boardSize; col++ {
			switch s.Boards.CellState(row, col) {
			case CELL_BLACK: // black
				drawDisk(screen, row, col, blackDiskImg)
			case CELL_WHITE: // white
				drawDisk(screen, row, col, whiteDiskImg)
			}
		}
	}
}

// drawDisk places the disk where desired row col format.
func drawDisk(screen *ebiten.Image, row, col int, img *ebiten.Image) {
	x := col*tileSize + (col+1)*tileMargin
	y := row*tileSize + (row+1)*tileMargin
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(img, op)
}

// drawDiskAtIndex places the disk where desired index format.
func drawDiskAtIndex(screen *ebiten.Image, i int, diskImg *ebiten.Image) {
	row := i / boardSize
	col := i % boardSize

	x := col*tileSize + (col+1)*tileMargin
	y := row*tileSize + (row+1)*tileMargin

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(diskImg, op)
}
