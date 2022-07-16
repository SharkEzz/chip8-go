package main

import (
	"flag"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/SharkEzz/chip8-go/pkg/emulator"
)

const MODIFIER = 15

type Game struct {
	emulator *emulator.Chip8
	ticker   *time.Ticker
}

func (g *Game) Update() error {
	select {
	default:
		return nil
	case <-g.ticker.C:
	}

	g.emulator.Cycle()

	g.processKeyPress()

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

func (g *Game) processKeyPress() {
	if ebiten.IsKeyPressed(ebiten.Key1) {
		g.emulator.SetKeyState(0x1, true)
	} else {
		g.emulator.SetKeyState(0x1, false)
	}
	if ebiten.IsKeyPressed(ebiten.Key2) {
		g.emulator.SetKeyState(0x2, true)
	} else {
		g.emulator.SetKeyState(0x2, false)
	}
	if ebiten.IsKeyPressed(ebiten.Key3) {
		g.emulator.SetKeyState(0x3, true)
	} else {
		g.emulator.SetKeyState(0x3, false)
	}
	if ebiten.IsKeyPressed(ebiten.Key4) {
		g.emulator.SetKeyState(0xC, true)
	} else {
		g.emulator.SetKeyState(0xC, false)
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.emulator.SetKeyState(0x4, true)
	} else {
		g.emulator.SetKeyState(0x4, false)
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.emulator.SetKeyState(0x5, true)
	} else {
		g.emulator.SetKeyState(0x5, false)
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		g.emulator.SetKeyState(0x6, true)
	} else {
		g.emulator.SetKeyState(0x6, false)
	}
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		g.emulator.SetKeyState(0xD, true)
	} else {
		g.emulator.SetKeyState(0xD, false)
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.emulator.SetKeyState(0x7, true)
	} else {
		g.emulator.SetKeyState(0x7, false)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.emulator.SetKeyState(0x8, true)
	} else {
		g.emulator.SetKeyState(0x8, false)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.emulator.SetKeyState(0x9, true)
	} else {
		g.emulator.SetKeyState(0x9, false)
	}
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		g.emulator.SetKeyState(0xE, true)
	} else {
		g.emulator.SetKeyState(0xE, false)
	}
	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		g.emulator.SetKeyState(0xA, true)
	} else {
		g.emulator.SetKeyState(0xA, false)
	}
	if ebiten.IsKeyPressed(ebiten.KeyX) {
		g.emulator.SetKeyState(0x0, true)
	} else {
		g.emulator.SetKeyState(0x0, false)
	}
	if ebiten.IsKeyPressed(ebiten.KeyC) {
		g.emulator.SetKeyState(0xB, true)
	} else {
		g.emulator.SetKeyState(0xB, false)
	}
	if ebiten.IsKeyPressed(ebiten.KeyV) {
		g.emulator.SetKeyState(0xF, true)
	} else {
		g.emulator.SetKeyState(0xF, false)
	}
}

func main() {
	file := flag.String("file", "", "The program to load")

	flag.Parse()

	chip8 := emulator.Init()

	if *file != "" {
		err := chip8.LoadProgram(*file)
		if err != nil {
			panic(err)
		}
	}

	ebiten.SetWindowSize(64*MODIFIER, 32*MODIFIER)
	ebiten.SetWindowTitle("Chip8")
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMaximum)
	ebiten.SetMaxTPS(-1)
	game := &Game{
		chip8,
		time.NewTicker(time.Second / 60),
	}
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
