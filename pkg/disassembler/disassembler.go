package disassembler

import (
	"fmt"
	"os"
)

type Line struct {
	OPCode      string // The hexadecimal representation of the value
	Instruction string // The assembly representation of the OPCode
}

type Disassembler struct {
	fileContent []byte
	pc          uint16
	Lines       []Line
}

func NewDisassembler(fileName string) (*Disassembler, error) {
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return &Disassembler{
		fileContent: fileContent,
		pc:          0,
	}, nil
}

func (d *Disassembler) Disassemble() []*Line {
	lines := []*Line{}

	for d.pc != uint16(len(d.fileContent)-1) && int(d.pc+1) <= len(d.fileContent) {
		op := uint16(d.fileContent[d.pc])<<8 | uint16(d.fileContent[d.pc+1])

		line := d.DisassembleOPCode(op)
		if line != nil {
			lines = append(lines, line)
		}

		d.pc += 2
	}

	return lines
}

func (d *Disassembler) DisassembleOPCode(op uint16) *Line {
	x := (op & 0x0F00) >> 8
	y := (op & 0x00F0) >> 4
	nnn := op & 0x0FFF
	kk := op & 0x00FF

	switch op & 0xF000 {
	case 0x0000:
		switch op & 0x000F {
		case 0x0000:
			return &Line{formatOP(op), "CLS"}
		case 0x000E:
			return &Line{formatOP(op), "RET"}
		default:
			return nil
		}
	case 0x1000:
		return &Line{formatOP(op), fmt.Sprintf("JP 0x%04X", nnn)}
	case 0x2000:
		return &Line{formatOP(op), fmt.Sprintf("CALL 0x%04X", nnn)}
	case 0x3000:
		return &Line{formatOP(op), fmt.Sprintf("SE V%X, 0x%04X", x, kk)}
	case 0x4000:
		return &Line{formatOP(op), fmt.Sprintf("SNE V%X, 0x%04X", x, kk)}
	case 0x5000:
		return &Line{formatOP(op), fmt.Sprintf("SE V%X, V%X", x, y)}
	case 0x6000:
		return &Line{formatOP(op), fmt.Sprintf("LD V%X, 0x%04X", x, kk)}
	case 0x7000:
		return &Line{formatOP(op), fmt.Sprintf("ADD V%X, 0x%04X", x, kk)}
	case 0x8000:
		switch op & 0x000F {
		case 0x0000:
			return &Line{formatOP(op), fmt.Sprintf("LD V%X, V%X", x, y)}
		case 0x0001:
			return &Line{formatOP(op), fmt.Sprintf("OR V%X, V%X", x, y)}
		case 0x0002:
			return &Line{formatOP(op), fmt.Sprintf("AND V%X, V%X", x, y)}
		case 0x0003:
			return &Line{formatOP(op), fmt.Sprintf("XOR V%X, V%X", x, y)}
		case 0x0004:
			return &Line{formatOP(op), fmt.Sprintf("ADD V%X, V%X", x, y)}
		case 0x0005:
			return &Line{formatOP(op), fmt.Sprintf("SUB V%X, V%X", x, y)}
		case 0x0006:
			return &Line{formatOP(op), fmt.Sprintf("SHR V%X, V%X", x, y)}
		case 0x0007:
			return &Line{formatOP(op), fmt.Sprintf("SUBN V%X, V%X", x, y)}
		case 0x000E:
			return &Line{formatOP(op), fmt.Sprintf("SHL V%X, V%X", x, y)}
		default:
			return nil
		}
	case 0x9000:
		return &Line{formatOP(op), fmt.Sprintf("SNE V%X, V%X", x, y)}
	case 0xA000:
		return &Line{formatOP(op), fmt.Sprintf("LD I, 0x%04X", nnn)}
	case 0xB000:
		return &Line{formatOP(op), fmt.Sprintf("JP V0, 0x%04X", nnn)}
	case 0xC000:
		return &Line{formatOP(op), fmt.Sprintf("RND V%d, 0x%04X", x, kk)}
	case 0xD000:
		return &Line{formatOP(op), fmt.Sprintf("DRW V%X, V%X, 0x%04X", x, y, op&0x000F)}
	case 0xE000:
		switch op & 0x00FF {
		case 0x009E:
			return &Line{formatOP(op), fmt.Sprintf("SKP V%X", x)}
		case 0x00A1:
			return &Line{formatOP(op), fmt.Sprintf("SKNP V%X", x)}
		default:
			return nil
		}
	case 0xF000:
		switch op & 0x00FF {
		case 0x0007:
			return &Line{formatOP(op), fmt.Sprintf("LD V%X, DT", x)}
		case 0x000A:
			return &Line{formatOP(op), fmt.Sprintf("LD V%X, K", x)}
		case 0x0015:
			return &Line{formatOP(op), fmt.Sprintf("LD DT, V%X", x)}
		case 0x0018:
			return &Line{formatOP(op), fmt.Sprintf("LD DT, V%X", x)}
		case 0x001E:
			return &Line{formatOP(op), fmt.Sprintf("ADD I, V%X", x)}
		case 0x0029:
			return &Line{formatOP(op), fmt.Sprintf("LD F, V%X", x)}
		case 0x0033:
			return &Line{formatOP(op), fmt.Sprintf("LD B, V%X", x)}
		case 0x0055:
			return &Line{formatOP(op), fmt.Sprintf("LD [I], V%X", x)}
		case 0x0065:
			return &Line{formatOP(op), fmt.Sprintf("LD V%X, [I]", x)}
		}
	default:
		return nil
	}

	return nil
}

func (d *Disassembler) WriteToFile(fileName string) error {
	bf := []byte{}

	for _, line := range d.Lines {
		bf = append(bf, []byte(fmt.Sprintln(line.Instruction))...)
	}

	return os.WriteFile(fileName, bf, 0666)
}

// Format an opcode to its hexadecimal representation (0x0000)
func formatOP(op uint16) string {
	return fmt.Sprintf("0x%04X", op)
}
