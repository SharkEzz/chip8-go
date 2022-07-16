package emulator

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

var fontSet = []uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, //0
	0x20, 0x60, 0x20, 0x20, 0x70, //1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, //2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, //3
	0x90, 0x90, 0xF0, 0x10, 0x10, //4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, //5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, //6
	0xF0, 0x10, 0x20, 0x40, 0x40, //7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, //8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, //9
	0xF0, 0x90, 0xF0, 0x90, 0x90, //A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, //B
	0xF0, 0x80, 0x80, 0x80, 0xF0, //C
	0xE0, 0x90, 0x90, 0x90, 0xE0, //D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, //E
	0xF0, 0x80, 0xF0, 0x80, 0x80, //F
}

type Chip8 struct {
	Display    [32][64]uint8 // display
	Memory     [4096]uint8   // memory
	V          [16]uint8     // registers -> V0 -> VF
	I          uint16        // index register
	DT         uint8         // delay timer
	ST         uint8         // sound timer
	PC         uint16        // program counter, points to the next instruction to be executed
	SP         uint16        // stack pointer
	Key        [16]uint8     // keypad -> true if pressed
	Stack      [16]uint16    // stack
	ShouldDraw bool
	Beeper     func() // beeper function
	Clock      *time.Ticker
}

func Init(clock uint) *Chip8 {
	c := &Chip8{
		PC:     0x200, // Program start at 0x200
		Beeper: func() { fmt.Print("\a") },
		Clock:  time.NewTicker(time.Second / time.Duration(clock)),
	}

	// Copy fontset to memory, starting at 0x000
	copy(c.Memory[0:], fontSet)

	return c
}

// Reset restart the program currently in the memory (at address 0x200)
func (c *Chip8) Reset() {
	c.Display = [32][64]uint8{}
	c.I = 0x0
	c.DT = 0x0
	c.ST = 0x0
	c.PC = 0x200
	c.SP = 0x0
	c.Key = [16]uint8{}
	c.Stack = [16]uint16{}
	c.ShouldDraw = false

	for i := 0; i < len(c.V); i++ {
		c.V[i] = 0x0
	}
}

func (c *Chip8) Buffer() [32][64]uint8 {
	return c.Display
}

func (c *Chip8) Draw() bool {
	sd := c.ShouldDraw
	c.ShouldDraw = false
	return sd
}

func (c *Chip8) SetKeyState(num uint8, down bool) {
	if down {
		c.Key[num] = 1
	} else {
		c.Key[num] = 0
	}
}

// Cycle represent a CPU cycle.
//
// By default the Chi8 CPU run at 60Hz.
func (c *Chip8) Cycle() uint16 {
	select {
	default:
		return 0
	case <-c.Clock.C:
	}

	op := uint16(c.Memory[c.PC])<<8 | uint16(c.Memory[c.PC+1]) // Combine the 2 bytes of the opcode

	c.processOP(op)

	if c.DT > 0 {
		c.DT -= 1
	}
	if c.ST > 0 {
		c.Beeper()
		c.ST -= 1
	}

	return op
}

// processOP take a opcode and process it
func (c *Chip8) processOP(op uint16) {
	x := (op & 0x0F00) >> 8
	y := (op & 0x00F0) >> 4
	nnn := op & 0x0FFF
	kk := op & 0x00FF

	c.nextInstruction()

	switch op & 0xF000 { // 0xF000 because we need to get only the first nibble
	case 0x0000:
		switch op & 0x000F { // 0x000F because we need to get only the last nibble
		case 0x0000:
			for i := 0; i < 32; i++ {
				for j := 0; j < 64; j++ {
					c.Display[i][j] = 0
				}
			}
			c.ShouldDraw = true
		case 0x000E: // Return from subroutine
			c.PC = c.Stack[c.SP] // Set program counter to the address at the top of the stack
			c.SP--               // Subtract 1 from stack pointer
		default:
			fmt.Printf("invalid opcode: %X\n", op)
		}
	case 0x1000: // Jump to address NNN
		c.PC = nnn
	case 0x2000: // Call subroutine at NNN
		c.SP++
		c.Stack[c.SP] = c.PC // Save current program counter to stack
		c.PC = nnn
	case 0x3000: // Skip next instruction if VX equals KK
		if uint16(c.V[x]) == kk {
			c.nextInstruction()
		}
	case 0x4000: // Skip next instruction if VX doesn't equal KK
		if uint16(c.V[x]) != kk {
			c.nextInstruction()
		}
	case 0x5000: // Skip next instruction if Vx = Vy.
		if c.V[x] == c.V[y] {
			c.nextInstruction()
		}
	case 0x6000: // Set Vx = kk.
		c.V[x] = uint8(kk)
	case 0x7000: // Set Vx = Vx + kk.
		c.V[x] = c.V[x] + uint8(kk)
	case 0x8000:
		switch op & 0x000F {
		case 0x0000: // Set Vx = Vy.
			c.V[x] = c.V[y]
		case 0x0001: // Set Vx = Vx OR Vy.
			c.V[x] = c.V[x] | c.V[y]
		case 0x0002: // Set Vx = Vx AND Vy.
			c.V[x] = c.V[x] & c.V[y]
		case 0x0003: // Set Vx = Vx XOR Vy.
			c.V[x] = c.V[x] ^ c.V[y]
		case 0x0004: // Set Vx = Vx + Vy, set VF = carry.
			r := uint16(c.V[x]) + uint16(c.V[y])
			var cf byte
			if r > 0xFF {
				cf = 1
			}
			c.V[0xF] = cf
			c.V[x] = uint8(r)
		case 0x0005: // Set Vx = Vx - Vy, set VF = NOT borrow.
			var cf byte
			if c.V[x] > c.V[y] {
				cf = 1
			}
			c.V[0xF] = cf
			c.V[x] = c.V[x] - c.V[y]
		case 0x0006: // Set Vx = Vx SHR 1.
			var cf byte
			if (c.V[x] & 0x01) == 0x01 {
				cf = 1
			}
			c.V[0xF] = cf
			c.V[x] = c.V[x] / 2
		case 0x0007: // Set Vx = Vy - Vx, set VF = NOT borrow.
			var cf byte
			if c.V[y] > c.V[x] {
				cf = 1
			}
			c.V[0xF] = cf
			c.V[x] = c.V[y] - c.V[x]
		case 0x000E: // Set Vx = Vx SHL 1.
			var cf byte
			if (c.V[x] & 0x80) == 0x80 {
				cf = 1
			}
			c.V[0xF] = cf
			c.V[x] = c.V[x] * 2
		default:
			fmt.Printf("invalid opcode: %X\n", op)
		}
	case 0x9000: // Skip next instruction if Vx != Vy.
		if c.V[x] != c.V[y] {
			c.nextInstruction()
		}
	case 0xA000: // Set I = nnn.
		c.I = nnn
	case 0xB000: // Jump to location nnn + V0.
		c.PC = nnn + uint16(c.V[0x0])
	case 0xC000: // Set Vx = random byte AND kk.
		c.V[x] = uint8(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(256)) & uint8(kk)
	case 0xD000: // Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision.
		x := c.V[x]
		y := c.V[y]
		n := op & 0x000F
		c.V[0xF] = 0
		var j uint16
		var i uint16
		for j = 0; j < n; j++ {
			pixel := c.Memory[c.I+j]
			for i = 0; i < 8; i++ {
				if (pixel & (0x80 >> i)) != 0 {
					posY := uint8(y) + uint8(j)
					posX := uint8(x) + uint8(i)

					// TODO: fix
					if posX > 63 {
						posX = 63
					}
					if posY > 31 {
						posY = 31
					}

					if c.Display[posY][posX] == 1 {
						c.V[0xF] = 1
					}
					c.Display[posY][posX] ^= 1
				}
			}
		}
		c.ShouldDraw = true
	case 0xE000:
		switch op & 0x00FF {
		case 0x009E: // Skip next instruction if key with the value of Vx is pressed.
			if c.Key[c.V[x]] == 1 {
				c.nextInstruction()
			}
		case 0x00A1: // Skip next instruction if key with the value of Vx is not pressed.
			if c.Key[c.V[x]] == 0 {
				c.nextInstruction()
			}
		default:
			fmt.Printf("invalid opcode: %X\n", op)
		}
	case 0xF000:
		switch op & 0x00FF {
		case 0x0007: // Set Vx = delay timer value.
			c.V[x] = c.DT
		case 0x000A: // Wait for a key press, store the value of the key in Vx.
			pressed := false
			for i := 0; i < len(c.Key); i++ {
				if c.Key[i] != 0 {
					c.V[x] = uint8(i)
					pressed = true
				}
			}
			if !pressed {
				return
			}
		case 0x0015: // Set delay timer = Vx.
			c.DT = c.V[x]
		case 0x0018: // Set sound timer = Vx.
			c.ST = c.V[x]
		case 0x001E: // Set I = I + Vx.
			c.I = c.I + uint16(c.V[x])
		case 0x0029: // Set I = location of sprite for digit Vx.
			c.I = uint16(c.V[x]) * 0x5
		case 0x0033: // 0xFX33 Stores the binary-coded decimal representation of VX, with the most significant of three digits at the address in I, the middle digit at I plus 1, and the least significant digit at I plus 2
			c.Memory[c.I] = c.V[x] / 100
			c.Memory[c.I+1] = (c.V[x] / 10) % 10
			c.Memory[c.I+2] = (c.V[x] % 100) % 10
		case 0x0055: // 0xFX55 Stores V0 to VX (including VX) in memory starting at address I. I is increased by 1 for each value written
			for i := 0; i < int(x)+1; i++ {
				c.Memory[uint16(i)+c.I] = c.V[i]
			}
			c.I = x + 1
		case 0x0065: // 0xFX65 Fills V0 to VX (including VX) with values from memory starting at address I. I is increased by 1 for each value written
			for i := 0; i < int(x)+1; i++ {
				c.V[i] = c.Memory[c.I+uint16(i)]
			}
			c.I = x + 1
		}
	default:
		fmt.Printf("invalid opcode: %X\n", op)
	}
}

// Increment the PC (program-counter) register by 2.
func (c *Chip8) nextInstruction() {
	c.PC += 2
}

// Load a binary Chip8 program, check its size, then copy it from 0x200 memory address.
func (c *Chip8) LoadProgram(fileName string) error {
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	if len(fileContent) > len(c.Memory)-0x200 {
		return fmt.Errorf("not enought memory")
	}

	copy(c.Memory[0x200:], fileContent)

	return nil
}
