package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const MODIFIER = 15

type Game struct {
	emulator *Chip8
}

func (g *Game) Update() error {
	g.emulator.Cycle()

	if ebiten.IsKeyPressed(ebiten.Key1) {
		g.emulator.Key(0x1, true)
	} else {
		g.emulator.Key(0x1, false)
	}
	if ebiten.IsKeyPressed(ebiten.Key2) {
		g.emulator.Key(0x2, true)
	} else {
		g.emulator.Key(0x2, false)
	}
	if ebiten.IsKeyPressed(ebiten.Key3) {
		g.emulator.Key(0x3, true)
	} else {
		g.emulator.Key(0x3, false)
	}
	if ebiten.IsKeyPressed(ebiten.Key4) {
		g.emulator.Key(0xC, true)
	} else {
		g.emulator.Key(0xC, false)
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.emulator.Key(0x4, true)
	} else {
		g.emulator.Key(0x4, false)
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.emulator.Key(0x5, true)
	} else {
		g.emulator.Key(0x5, false)
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		g.emulator.Key(0x6, true)
	} else {
		g.emulator.Key(0x6, false)
	}
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		g.emulator.Key(0xD, true)
	} else {
		g.emulator.Key(0xD, false)
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.emulator.Key(0x7, true)
	} else {
		g.emulator.Key(0x7, false)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.emulator.Key(0x8, true)
	} else {
		g.emulator.Key(0x8, false)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.emulator.Key(0x9, true)
	} else {
		g.emulator.Key(0x9, false)
	}
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		g.emulator.Key(0xE, true)
	} else {
		g.emulator.Key(0xE, false)
	}
	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		g.emulator.Key(0xA, true)
	} else {
		g.emulator.Key(0xA, false)
	}
	if ebiten.IsKeyPressed(ebiten.KeyX) {
		g.emulator.Key(0x0, true)
	} else {
		g.emulator.Key(0x0, false)
	}
	if ebiten.IsKeyPressed(ebiten.KeyC) {
		g.emulator.Key(0xB, true)
	} else {
		g.emulator.Key(0xB, false)
	}
	if ebiten.IsKeyPressed(ebiten.KeyV) {
		g.emulator.Key(0xF, true)
	} else {
		g.emulator.Key(0xF, false)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.emulator.Draw() {
		buffer := g.emulator.Buffer()
		for j := 0; j < len(buffer); j++ {
			for i := 0; i < len(buffer[j]); i++ {
				if buffer[j][i] != 0 {
					ebitenutil.DrawRect(screen, float64(i*MODIFIER), float64(j*MODIFIER), MODIFIER, MODIFIER, color.RGBA{51, 255, 102, 255})
				} else {
					ebitenutil.DrawRect(screen, float64(i*MODIFIER), float64(j*MODIFIER), MODIFIER, MODIFIER, color.RGBA{0, 0, 0, 255})
				}
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {
	emulator := Init()
	emulator.Beeper(func() { fmt.Print("\a") })
	err := emulator.LoadProgram("./brick.ch8")
	if err != nil {
		panic(err)
	}

	ebiten.SetWindowSize(64*MODIFIER, 32*MODIFIER)
	ebiten.SetWindowTitle("Chip8")
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetMaxTPS(260)
	game := &Game{
		emulator,
	}
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
