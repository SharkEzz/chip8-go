package main

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
	display [32][64]uint8 // display

	memory [4096]uint8 // memory

	vx [16]uint8 // registers -> V0 -> VF
	vi uint16    // index register

	dt uint8 // delay timer
	st uint8 // sound timer

	pc uint16 // program counter, points to the next instruction to be executed
	sp uint16 // stack pointer

	oc uint16 // current opcode, holds the current instruction

	key [16]uint8 // keypad -> true if pressed

	stack [16]uint16 // stack

	beeper func() // beeper function
}

func Init() *Chip8 {
	c := &Chip8{
		pc:     0x200,
		beeper: func() {},
	}

	// Copy fontset to memory
	copy(c.memory[0:], fontSet)

	return c
}

func (c *Chip8) Buffer() [32][64]uint8 {
	return c.display
}

func (c *Chip8) Beeper(f func()) {
	c.beeper = f
}

func (c *Chip8) Key(num uint8, down bool) {
	if down {
		c.key[num] = 1
	} else {
		c.key[num] = 0
	}
}

func (c *Chip8) Cycle() {
	c.oc = uint16(c.memory[c.pc])<<8 | uint16(c.memory[c.pc+1]) // Combine the 2 bytes of the opcode

	fmt.Printf("op=0x%04X \n", c.oc)

	switch c.oc & 0xF000 { // 0xF000 because we need to get only the first nibble
	case 0x0000:
		switch c.oc & 0x000F { // 0x000F because we need to get only the last nibble
		case 0x0000:
			for i := 0; i < 32; i++ {
				for j := 0; j < 64; j++ {
					c.display[i][j] = 0
				}
			}
			c.pc += 2
		case 0x000E: // Return from subroutine
			c.pc = c.stack[c.sp] // Set program counter to the address at the top of the stack
			c.sp--               // Subtract 1 from stack pointer
			c.pc += 2
		default:
			fmt.Printf("invalid opcode: %X\n", c.oc)
		}
	case 0x1000: // Jump to address NNN
		c.pc = uint16(c.oc & 0x0FFF)
	case 0x2000: // Call subroutine at NNN
		c.sp++
		c.stack[c.sp] = c.pc // Save current program counter to stack
		c.pc = uint16(c.oc & 0x0FFF)
	case 0x3000: // Skip next instruction if VX equals NN
		if uint16(c.vx[(c.oc&0x0F00)>>8]) == c.oc&0x00FF {
			c.pc += 4
		} else {
			c.pc += 2
		}
	case 0x4000: // Skip next instruction if VX doesn't equal NN
		if uint16(c.vx[(c.oc&0x0F00)>>8]) != c.oc&0x00FF {
			c.pc += 4
		} else {
			c.pc += 2
		}
	case 0x5000: // Skip next instruction if Vx = Vy.
		if c.vx[(c.oc&0x0F00)>>8] == c.vx[(c.oc&0x00F0)>>4] {
			c.pc += 4
		} else {
			c.pc += 2
		}
	case 0x6000: // Set Vx = kk.
		c.vx[c.oc&0x0F00>>8] = uint8(c.oc & 0x00FF)
		c.pc += 2
	case 0x7000: // Set Vx = Vx + kk.
		c.vx[(c.oc&0x0F00)>>8] = c.vx[(c.oc&0x0F00)>>8] + uint8(c.oc&0x00FF)
		c.pc += 2
	case 0x8000:
		switch c.oc & 0x000F {
		case 0x0000: // Set Vx = Vy.
			c.vx[(c.oc&0x0F00)>>8] = c.vx[(c.oc&0x00F0)>>4]
			c.pc += 2
		case 0x0001: // Set Vx = Vx OR Vy.
			c.vx[(c.oc&0x0F00)>>8] = c.vx[(c.oc&0x0F00)>>8] | c.vx[(c.oc&0x00F0)>>4]
			c.pc += 2
		case 0x0002: // Set Vx = Vx AND Vy.
			c.vx[(c.oc&0x0F00)>>8] = c.vx[(c.oc&0x0F00)>>8] & c.vx[(c.oc&0x00F0)>>4]
			c.pc += 2
		case 0x0003: // Set Vx = Vx XOR Vy.
			c.vx[(c.oc&0x0F00)>>8] = c.vx[(c.oc&0x0F00)>>8] ^ c.vx[(c.oc&0x00F0)>>4]
			c.pc += 2
		case 0x0004: // Set Vx = Vx + Vy, set VF = carry.
			if c.vx[(c.oc&0x00F0)>>4] > 0xFF-c.vx[(c.oc&0x0F00)>>8] {
				c.vx[0xF] = 1
			} else {
				c.vx[0xF] = 0
			}
			c.vx[(c.oc&0x0F00)>>8] = c.vx[(c.oc&0x0F00)>>8] + c.vx[(c.oc&0x00F0)>>4]
			c.pc = c.pc + 2
		case 0x0005: // Set Vx = Vx - Vy, set VF = NOT borrow.
			if c.vx[(c.oc&0x00F0)>>4] > c.vx[(c.oc&0x0F00)>>8] {
				c.vx[0xF] = 0
			} else {
				c.vx[0xF] = 1
			}
			c.vx[(c.oc&0x0F00)>>8] = c.vx[(c.oc&0x0F00)>>8] - c.vx[(c.oc&0x00F0)>>4]
			c.pc = c.pc + 2
		case 0x0006: // Set Vx = Vx SHR 1.
			c.vx[0xF] = c.vx[(c.oc&0x0F00)>>8] >> 7
			c.vx[(c.oc&0x0F00)>>8] = c.vx[(c.oc&0x0F00)>>8] << 1
			c.pc = c.pc + 2
		case 0x0007: // Set Vx = Vy - Vx, set VF = NOT borrow.
			if c.vx[(c.oc&0x0F00)>>8] > c.vx[(c.oc&0x00F0)>>4] {
				c.vx[0xF] = 0
			} else {
				c.vx[0xF] = 1
			}
			c.vx[(c.oc&0x0F00)>>8] = c.vx[(c.oc&0x00F0)>>4] - c.vx[(c.oc&0x0F00)>>8]
			c.pc = c.pc + 2
		case 0x000E: // Set Vx = Vx SHL 1.
			c.vx[0xF] = c.vx[(c.oc&0x0F00)>>8] >> 7
			c.vx[(c.oc&0x0F00)>>8] = c.vx[(c.oc&0x0F00)>>8] << 1
			c.pc = c.pc + 2
		default:
			fmt.Printf("invalid opcode: %X\n", c.oc)
		}
	case 0x9000: // Skip next instruction if Vx != Vy.
		if c.vx[c.oc&0x0F00>>8] != c.vx[(c.oc&0x00F0)>>4] {
			c.pc += 4
		} else {
			c.pc += 2
		}
	case 0xA000: // Set I = nnn.
		c.vi = c.oc & 0x0FFF
		c.pc += 2
	case 0xB000: // Jump to location nnn + V0.
		c.pc = (c.oc & 0x0FFF) + uint16(c.vx[0x0])
	case 0xC000: // Set Vx = random byte AND kk.
		c.vx[c.oc&0x0F00>>8] = uint8(rand.Intn(256)) & uint8(c.oc&0x00FF)
		c.pc += 2
	case 0xD000: // Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision.
		x := c.vx[(c.oc&0x0F00)>>8]
		y := c.vx[(c.oc&0x00F0)>>4]
		n := c.oc & 0x000F
		c.vx[0xF] = 0
		var j uint16 = 0
		var i uint16 = 0
		for j = 0; j < n; j++ {
			pixel := c.memory[c.vi+j]
			for i = 0; i < 8; i++ {
				if (pixel & (0x80 >> i)) != 0 {
					if c.display[(y + uint8(j))][x+uint8(i)] == 1 {
						c.vx[0xF] = 1
					}
					c.display[(y + uint8(j))][x+uint8(i)] ^= 1
				}
			}
		}
		c.pc = c.pc + 2
	case 0xE000:
		switch c.oc & 0x00FF {
		case 0x009E: // Skip next instruction if key with the value of Vx is pressed.
			if c.key[c.vx[(c.oc&0x0F00)>>8]] == 1 {
				c.pc += 4
			} else {
				c.pc += 2
			}
		case 0x00A1: // Skip next instruction if key with the value of Vx is not pressed.
			if c.key[c.vx[(c.oc&0x0F00)>>8]] == 0 {
				c.pc += 4
			} else {
				c.pc += 2
			}
		default:
			fmt.Printf("invalid opcode: %X\n", c.oc)
		}
	case 0xF000:
		switch c.oc & 0x00FF {
		case 0x0007: // Set Vx = delay timer value.
			c.vx[(c.oc&0x0F00)>>8] = c.dt
			c.pc += 2
		case 0x000A: // Wait for a key press, store the value of the key in Vx.
			pressed := false
			for i := 0; i < len(c.key); i++ {
				if c.key[i] != 0 {
					c.vx[(c.oc&0x0F00)>>8] = uint8(i)
					pressed = true
				}
			}
			if !pressed {
				return
			}
			c.pc += 2
		case 0x0015: // Set delay timer = Vx.
			c.dt = c.vx[(c.oc&0xF00)>>8]
			c.pc += 2
		case 0x0018: // Set sound timer = Vx.
			c.st = c.vx[(c.oc&0xF00)>>8]
			c.pc += 2
		case 0x001E: // Set I = I + Vx.
			c.vi = c.vi + uint16(c.vx[(c.oc&0x0F00)>>8])
			c.pc += 2
		case 0x0029: // Set I = location of sprite for digit Vx.
			c.vi = uint16(c.vx[(c.oc&0x0F00)>>8]) * 0x5
			c.pc += 2
		case 0x0033: // 0xFX33 Stores the binary-coded decimal representation of VX, with the most significant of three digits at the address in I, the middle digit at I plus 1, and the least significant digit at I plus 2
			c.memory[c.vi] = c.vx[(c.oc&0x0F00)>>8] / 100
			c.memory[c.vi+1] = (c.vx[(c.oc&0x0F00)>>8] / 10) % 10
			c.memory[c.vi+2] = (c.vx[(c.oc&0x0F00)>>8] % 100) / 10
			c.pc = c.pc + 2
		case 0x0055: // 0xFX55 Stores V0 to VX (including VX) in memory starting at address I. I is increased by 1 for each value written
			for i := 0; i < int((c.oc&0x0F00)>>8)+1; i++ {
				c.memory[uint16(i)+c.vi] = c.vx[i]
			}
			c.vi = ((c.oc & 0x0F00) >> 8) + 1
			c.pc = c.pc + 2
		case 0x0065: // 0xFX65 Fills V0 to VX (including VX) with values from memory starting at address I. I is increased by 1 for each value written
			for i := 0; i < int((c.oc&0x0F00)>>8)+1; i++ {
				c.vx[i] = c.memory[c.vi+uint16(i)]
			}
			c.vi = ((c.oc & 0x0F00) >> 8) + 1
			c.pc = c.pc + 2
		}
	default:
		fmt.Printf("invalid opcode: %X\n", c.oc)
	}

	if c.dt > 0 {
		c.dt -= 1
	}
	if c.st > 0 {
		c.beeper()
		c.st -= 1
	}
}

func (c *Chip8) LoadProgram(fileName string) error {
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	if len(fileContent) > len(c.memory)-0x200 {
		return fmt.Errorf("not enought memory")
	}

	copy(c.memory[0x200:], fileContent)

	return nil
}
