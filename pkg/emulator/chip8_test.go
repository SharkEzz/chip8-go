package emulator

import "testing"

func TestReturnOpcode00EE(t *testing.T) {
	c := Init()
	c.SP = 1
	c.Stack[0] = 0x456
	c.Memory[c.PC] = 0x00
	c.Memory[c.PC+1] = 0xEE

	c.Cycle()

	if c.PC != 0x456 {
		t.Fatalf("PC = %#x, want %#x", c.PC, 0x456)
	}
	if c.SP != 0 {
		t.Fatalf("SP = %d, want 0", c.SP)
	}
}

func TestCallOpcodeStoresCurrentPC(t *testing.T) {
	c := Init()
	c.Memory[c.PC] = 0x22
	c.Memory[c.PC+1] = 0x34

	c.Cycle()

	if c.PC != 0x234 {
		t.Fatalf("PC = %#x, want %#x", c.PC, 0x234)
	}
	if c.SP != 1 {
		t.Fatalf("SP = %d, want 1", c.SP)
	}
	if c.Stack[0] != 0x202 {
		t.Fatalf("Stack[0] = %#x, want %#x", c.Stack[0], 0x202)
	}
}

func TestFx55AndFx65DoNotMutateI(t *testing.T) {
	c := Init()
	c.I = 0x300
	c.V[0] = 0xAA
	c.V[1] = 0xBB
	c.V[2] = 0xCC

	c.processOP(0xF255)
	if c.I != 0x300 {
		t.Fatalf("I after Fx55 = %#x, want %#x", c.I, 0x300)
	}
	if c.Memory[0x300] != 0xAA || c.Memory[0x301] != 0xBB || c.Memory[0x302] != 0xCC {
		t.Fatalf("memory not written as expected")
	}

	c.V[0], c.V[1], c.V[2] = 0, 0, 0
	c.processOP(0xF265)
	if c.I != 0x300 {
		t.Fatalf("I after Fx65 = %#x, want %#x", c.I, 0x300)
	}
	if c.V[0] != 0xAA || c.V[1] != 0xBB || c.V[2] != 0xCC {
		t.Fatalf("registers not restored as expected")
	}
}
