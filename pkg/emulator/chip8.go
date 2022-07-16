package emulator

import (
	"fmt"
	"math/rand"
	"os"
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
}

func Init() *Chip8 {
	c := &Chip8{
		PC:     0x200, // Program start at 0x200
		Beeper: func() { fmt.Print("\a") },
	}

	// Copy fontset to memory, starting at 0x000
	copy(c.Memory[0:], fontSet)

	return c
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
func (c *Chip8) Cycle() {
	op := uint16(c.Memory[c.PC])<<8 | uint16(c.Memory[c.PC+1]) // Combine the 2 bytes of the opcode

	fmt.Printf("op=0x%04X\n", op)

	c.processOP(op)

	if c.DT > 0 {
		c.DT -= 1
	}
	if c.ST > 0 {
		c.Beeper()
		c.ST -= 1
	}
}

// processOP take a opcode and process it
func (c *Chip8) processOP(op uint16) {
	switch op & 0xF000 { // 0xF000 because we need to get only the first nibble
	case 0x0000:
		switch op & 0x000F { // 0x000F because we need to get only the last nibble
		case 0x0000:
			for i := 0; i < 32; i++ {
				for j := 0; j < 64; j++ {
					c.Display[i][j] = 0
				}
			}
			c.incrementPC()
			c.ShouldDraw = true
		case 0x000E: // Return from subroutine
			c.PC = c.Stack[c.SP] // Set program counter to the address at the top of the stack
			c.SP--               // Subtract 1 from stack pointer
			c.incrementPC()
		default:
			fmt.Printf("invalid opcode: %X\n", op)
		}
	case 0x1000: // Jump to address NNN
		c.PC = uint16(op & 0x0FFF)
	case 0x2000: // Call subroutine at NNN
		c.SP++
		c.Stack[c.SP] = c.PC // Save current program counter to stack
		c.PC = uint16(op & 0x0FFF)
	case 0x3000: // Skip next instruction if VX equals NN
		c.incrementPC()
		if uint16(c.V[(op&0x0F00)>>8]) == op&0x00FF {
			c.incrementPC()
		}
	case 0x4000: // Skip next instruction if VX doesn't equal NN
		c.incrementPC()
		if uint16(c.V[(op&0x0F00)>>8]) != op&0x00FF {
			c.incrementPC()
		}
	case 0x5000: // Skip next instruction if Vx = Vy.
		c.incrementPC()
		if c.V[(op&0x0F00)>>8] == c.V[(op&0x00F0)>>4] {
			c.incrementPC()
		}
	case 0x6000: // Set Vx = kk.
		c.V[op&0x0F00>>8] = uint8(op & 0x00FF)
		c.incrementPC()
	case 0x7000: // Set Vx = Vx + kk.
		c.V[(op&0x0F00)>>8] = c.V[(op&0x0F00)>>8] + uint8(op&0x00FF)
		c.incrementPC()
	case 0x8000:
		x := (op & 0x0F00) >> 8
		y := (op & 0x00F0) >> 4
		switch op & 0x000F {
		case 0x0000: // Set Vx = Vy.
			c.V[x] = c.V[y]
			c.incrementPC()
		case 0x0001: // Set Vx = Vx OR Vy.
			c.V[x] = c.V[x] | c.V[y]
			c.incrementPC()
		case 0x0002: // Set Vx = Vx AND Vy.
			c.V[x] = c.V[x] & c.V[y]
			c.incrementPC()
		case 0x0003: // Set Vx = Vx XOR Vy.
			c.V[x] = c.V[x] ^ c.V[y]
			c.incrementPC()
		case 0x0004: // Set Vx = Vx + Vy, set VF = carry.
			r := uint16(c.V[x]) + uint16(c.V[y])
			var cf byte
			if r > 0xFF {
				cf = 1
			}
			c.V[0xF] = cf
			c.V[x] = byte(r)
			c.incrementPC()
		case 0x0005: // Set Vx = Vx - Vy, set VF = NOT borrow.
			var cf byte
			if c.V[x] > c.V[y] {
				cf = 1
			}
			c.V[0xF] = cf
			c.V[x] = c.V[x] - c.V[y]
			c.incrementPC()
		case 0x0006: // Set Vx = Vx SHR 1.
			var cf byte
			if (c.V[x] & 0x01) == 0x01 {
				cf = 1
			}
			c.V[0xF] = cf
			c.V[x] = c.V[x] / 2
			c.incrementPC()
		case 0x0007: // Set Vx = Vy - Vx, set VF = NOT borrow.
			var cf byte
			if c.V[y] > c.V[x] {
				cf = 1
			}
			c.V[0xF] = cf
			c.V[x] = c.V[y] - c.V[x]
			c.incrementPC()
		case 0x000E: // Set Vx = Vx SHL 1.
			var cf byte
			if (c.V[x] & 0x80) == 0x80 {
				cf = 1
			}
			c.V[0xF] = cf
			c.V[x] = c.V[x] * 2
			c.incrementPC()
		default:
			fmt.Printf("invalid opcode: %X\n", op)
		}
	case 0x9000: // Skip next instruction if Vx != Vy.
		c.incrementPC()
		if c.V[op&0x0F00>>8] != c.V[(op&0x00F0)>>4] {
			c.incrementPC()
		}
	case 0xA000: // Set I = nnn.
		c.I = op & 0x0FFF
		c.incrementPC()
	case 0xB000: // Jump to location nnn + V0.
		c.PC = (op & 0x0FFF) + uint16(c.V[0x0])
	case 0xC000: // Set Vx = random byte AND kk.
		c.V[op&0x0F00>>8] = uint8(rand.Intn(256)) & uint8(op&0x00FF)
		c.incrementPC()
	case 0xD000: // Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision.
		x := c.V[(op&0x0F00)>>8]
		y := c.V[(op&0x00F0)>>4]
		n := op & 0x000F
		c.V[0xF] = 0
		var j uint16
		var i uint16
		for j = 0; j < n; j++ {
			pixel := c.Memory[c.I+j]
			for i = 0; i < 8; i++ {
				if (pixel & (0x80 >> i)) != 0 {
					posY := y + uint8(j)
					posX := x + uint8(i)

					if posY >= 32 {
						posY = 31
					}
					if posX >= 64 {
						posX = 63
					}

					if c.Display[posY][posX] == 1 {
						c.V[0xF] = 1
					}
					c.Display[posY][posX] ^= 1
				}
			}
		}
		c.incrementPC()
		c.ShouldDraw = true
	case 0xE000:
		switch op & 0x00FF {
		case 0x009E: // Skip next instruction if key with the value of Vx is pressed.
			c.incrementPC()
			if c.Key[c.V[(op&0x0F00)>>8]] == 1 {
				c.incrementPC()
			}
		case 0x00A1: // Skip next instruction if key with the value of Vx is not pressed.
			c.incrementPC()
			if c.Key[c.V[(op&0x0F00)>>8]] == 0 {
				c.incrementPC()
			}
		default:
			fmt.Printf("invalid opcode: %X\n", op)
		}
	case 0xF000:
		switch op & 0x00FF {
		case 0x0007: // Set Vx = delay timer value.
			c.V[(op&0x0F00)>>8] = c.DT
			c.incrementPC()
		case 0x000A: // Wait for a key press, store the value of the key in Vx.
			pressed := false
			for i := 0; i < len(c.Key); i++ {
				if c.Key[i] != 0 {
					c.V[(op&0x0F00)>>8] = uint8(i)
					pressed = true
				}
			}
			if !pressed {
				return
			}
			c.incrementPC()
		case 0x0015: // Set delay timer = Vx.
			c.DT = c.V[(op&0xF00)>>8]
			c.incrementPC()
		case 0x0018: // Set sound timer = Vx.
			c.ST = c.V[(op&0xF00)>>8]
			c.incrementPC()
		case 0x001E: // Set I = I + Vx.
			c.I = c.I + uint16(c.V[(op&0x0F00)>>8])
			c.incrementPC()
		case 0x0029: // Set I = location of sprite for digit Vx.
			c.I = uint16(c.V[(op&0x0F00)>>8]) * 0x5
			c.incrementPC()
		case 0x0033: // 0xFX33 Stores the binary-coded decimal representation of VX, with the most significant of three digits at the address in I, the middle digit at I plus 1, and the least significant digit at I plus 2
			c.Memory[c.I] = c.V[(op&0x0F00)>>8] / 100
			c.Memory[c.I+1] = (c.V[(op&0x0F00)>>8] / 10) % 10
			c.Memory[c.I+2] = (c.V[(op&0x0F00)>>8] % 100) % 10
			c.incrementPC()
		case 0x0055: // 0xFX55 Stores V0 to VX (including VX) in memory starting at address I. I is increased by 1 for each value written
			for i := 0; i < int((op&0x0F00)>>8)+1; i++ {
				c.Memory[uint16(i)+c.I] = c.V[i]
			}
			c.I = ((op & 0x0F00) >> 8) + 1
			c.incrementPC()
		case 0x0065: // 0xFX65 Fills V0 to VX (including VX) with values from memory starting at address I. I is increased by 1 for each value written
			for i := 0; i < int((op&0x0F00)>>8)+1; i++ {
				c.V[i] = c.Memory[c.I+uint16(i)]
			}
			c.I = ((op & 0x0F00) >> 8) + 1
			c.incrementPC()
		}
	default:
		fmt.Printf("invalid opcode: %X\n", op)
	}
}

// Increment the PC register by 2.
func (c *Chip8) incrementPC() {
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
