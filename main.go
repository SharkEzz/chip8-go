package main

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const MODIFIER = 15

type Game struct {
	emulator *Chip8
}

func (g *Game) Update() error {
	g.emulator.Cycle()

	time.Sleep(1000 / 60)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	buffer := g.emulator.Buffer()
	for j := 0; j < len(buffer); j++ {
		for i := 0; i < len(buffer[j]); i++ {
			if buffer[j][i] != 0 {
				ebitenutil.DrawRect(screen, float64(i*MODIFIER), float64(j*MODIFIER), MODIFIER, MODIFIER, color.RGBA{255, 255, 0, 255})
			} else {
				ebitenutil.DrawRect(screen, float64(i*MODIFIER), float64(j*MODIFIER), MODIFIER, MODIFIER, color.RGBA{255, 0, 0, 255})
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {
	emulator := Init()
	err := emulator.LoadProgram("./stars.ch8")
	if err != nil {
		panic(err)
	}

	ebiten.SetWindowSize(64*MODIFIER, 32*MODIFIER)
	ebiten.SetWindowTitle("Chip8")
	game := &Game{
		emulator,
	}
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
