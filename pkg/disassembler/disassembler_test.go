package disassembler

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDisassembleOpcodeMappings(t *testing.T) {
	tests := []struct {
		op   uint16
		want string
	}{
		{0x00EE, "RET"},
		{0xA123, "LD I, 0x0123"},
		{0xB321, "JP V0, 0x0321"},
		{0xC4FF, "RND V4, 0x00FF"},
		{0xF218, "LD ST, V2"},
	}

	for _, tt := range tests {
		line := DisassembleOPCode(tt.op)
		if line == nil {
			t.Fatalf("DisassembleOPCode(%#04x) returned nil", tt.op)
		}

		if line.Instruction != tt.want {
			t.Fatalf("DisassembleOPCode(%#04x) = %q, want %q", tt.op, line.Instruction, tt.want)
		}
	}
}

func TestDisassembleOddLengthProgramDoesNotPanic(t *testing.T) {
	dir := t.TempDir()
	program := filepath.Join(dir, "program.ch8")
	if err := os.WriteFile(program, []byte{0x00, 0xE0, 0x12}, 0o644); err != nil {
		t.Fatal(err)
	}

	d, err := NewDisassembler(program)
	if err != nil {
		t.Fatal(err)
	}

	lines := d.Disassemble()
	if len(lines) != 1 {
		t.Fatalf("len(lines) = %d, want 1", len(lines))
	}

	if lines[0].Instruction != "CLS" {
		t.Fatalf("lines[0] = %q, want %q", lines[0].Instruction, "CLS")
	}
}
