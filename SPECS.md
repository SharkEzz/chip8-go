# CPU

## 16 registers of 1 byte (general purpose):
V0, V1, V2, V3, V4, V5, V6, V7, V8, V9, VA, VB, VC, VD, VE, VF

I -> store memory addresses (12 right bits are used)

## 2 special registers for delay & sound

When those registers are non-zero, they are decremented at the rate of 60Hz

- DT: delay timer -> 8 bit
- ST: sound timer -> 8 bit

## Pseudo-registers (not accessibles by programs)
- PC: program counter -> 16 bit -> store the currently executing address
- SP: stack pointer -> 8 bit -> point to the top of the stack

## Stack

Array of 16 16 bit values -> store the address to return when a subroutine is finished -> 16 because Chip-8 allow up to 16 nester subroutines.

# Memory
0x000 -> 0xFFF 			: Total memory (4096 bytes)

0x000 -> 0x1FF 			: Interpreter

(0x200 | 0x600) -> 0xFFF 	: Program

# Keyboard

Keys:
- 1
- 2
- 3
- 4
- 5
- 6
- 7
- 8
- 9
- 0
- A
- B
- C
- D
- E
- F

# Display

64x32 pixels monochrome, horizontal

## Sprite

A group of bytes who are the representation of the picture (0 or 1)

### Built-in

5 bytes groups, should be included in the interpreter area of the memory

- 0
- 1
- 2
- 3
- 4
- 5
- 6
- 7
- 8
- 9
- A
- B
- C
- D
- E
- F

# Sound & Timers

- When ST > 0 -> buzzer enabled
- The buzzer has only one tone (decided by the author of the interpreter)

# Instructions

36 instructions in the original implementation.

All instructions are 2 bytes long, stored with the most significant byte first. The first byte of each instruction must be located at an even address.

If sprite data is included, it must be filled so the instructions following it will be properly placed in RAM.

- _nnn_ or _addr_ -> 12 bit value, lowest 12 bits of the instruction
- _n_ or _nibble_ -> 4 bit value, lowest 4 bits of the instruction
- _x_ -> 4 bit value -> lower 4 bits of the high byte of the instruction
- _y_ -> 4 bit value -> upper 4 bits of the low byte of the instruction
- _kk_ or _byte_ -> 8 bit value, lowest 8 bits of the instruction

[Reference](http://devernay.free.fr/hacks/chip8/C8TECH10.HTM)