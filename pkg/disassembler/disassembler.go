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
	fileContent     []byte
	pc              uint16
	Lines           []*Line
	hasDisassembled bool
}

func NewDisassembler(fileName string) (*Disassembler, error) {
	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return &Disassembler{
		fileContent: fileContent,
	}, nil
}

func (d *Disassembler) Disassemble() []*Line {
	lines := []*Line{}

	for d.pc != uint16(len(d.fileContent)-1) && int(d.pc+1) <= len(d.fileContent) {
		op := uint16(d.fileContent[d.pc])<<8 | uint16(d.fileContent[d.pc+1])

		line := DisassembleOPCode(op)
		if line != nil {
			lines = append(lines, line)
		}

		d.pc += 2
	}

	d.Lines = lines
	d.hasDisassembled = true

	return lines
}

func DisassembleOPCode(op uint16) *Line {
	x := fmt.Sprintf("%X", uint(op&0x0F00)>>8)
	y := fmt.Sprintf("%X", (op&0x00F0)>>4)
	nnn := fmt.Sprintf("0x%04X", op&0x0FFF)
	kk := fmt.Sprintf("0x%04X", op&0x00FF)

	formattedOP := formatOP(op)

	switch op & 0xF000 {
	case 0x0000:
		switch op & 0x00FF {
		case 0x00E0:
			return &Line{formattedOP, "CLS"}
		case 0x000E:
			return &Line{formattedOP, "RET"}
		default:
			return nil
		}
	case 0x1000:
		return &Line{formattedOP, fmt.Sprintf("JP %s", nnn)}
	case 0x2000:
		return &Line{formattedOP, fmt.Sprintf("CALL %s", nnn)}
	case 0x3000:
		return &Line{formattedOP, fmt.Sprintf("SE V%s, %s", x, kk)}
	case 0x4000:
		return &Line{formattedOP, fmt.Sprintf("SNE V%s, %s", x, kk)}
	case 0x5000:
		return &Line{formattedOP, fmt.Sprintf("SE V%s, %s", x, y)}
	case 0x6000:
		return &Line{formattedOP, fmt.Sprintf("LD V%s, %s", x, kk)}
	case 0x7000:
		return &Line{formattedOP, fmt.Sprintf("ADD V%s, %s", x, kk)}
	case 0x8000:
		switch op & 0x000F {
		case 0x0000:
			return &Line{formattedOP, fmt.Sprintf("LD V%s, V%s", x, y)}
		case 0x0001:
			return &Line{formattedOP, fmt.Sprintf("OR V%s, V%s", x, y)}
		case 0x0002:
			return &Line{formattedOP, fmt.Sprintf("AND V%s, V%s", x, y)}
		case 0x0003:
			return &Line{formattedOP, fmt.Sprintf("XOR V%s, V%s", x, y)}
		case 0x0004:
			return &Line{formattedOP, fmt.Sprintf("ADD V%s, V%s", x, y)}
		case 0x0005:
			return &Line{formattedOP, fmt.Sprintf("SUB V%s, V%s", x, y)}
		case 0x0006:
			return &Line{formattedOP, fmt.Sprintf("SHR V%s, V%s", x, y)}
		case 0x0007:
			return &Line{formattedOP, fmt.Sprintf("SUBN V%s, V%s", x, y)}
		case 0x000E:
			return &Line{formattedOP, fmt.Sprintf("SHL V%s, V%s", x, y)}
		default:
			return nil
		}
	case 0x9000:
		return &Line{formattedOP, fmt.Sprintf("SNE V%s, V%s", x, y)}
	case 0xA000:
		return &Line{formattedOP, fmt.Sprintf("LD I, 0x%s", nnn)}
	case 0xB000:
		return &Line{formattedOP, fmt.Sprintf("JP V0, 0x%s", nnn)}
	case 0xC000:
		return &Line{formattedOP, fmt.Sprintf("RND V%s, 0x%s", x, kk)}
	case 0xD000:
		return &Line{formattedOP, fmt.Sprintf("DRW V%s, V%s, 0x%04X", x, y, op&0x000F)}
	case 0xE000:
		switch op & 0x00FF {
		case 0x009E:
			return &Line{formattedOP, fmt.Sprintf("SKP V%s", x)}
		case 0x00A1:
			return &Line{formattedOP, fmt.Sprintf("SKNP V%s", x)}
		default:
			return nil
		}
	case 0xF000:
		switch op & 0x00FF {
		case 0x0007:
			return &Line{formattedOP, fmt.Sprintf("LD V%s, DT", x)}
		case 0x000A:
			return &Line{formattedOP, fmt.Sprintf("LD V%s, K", x)}
		case 0x0015:
			return &Line{formattedOP, fmt.Sprintf("LD DT, V%s", x)}
		case 0x0018:
			return &Line{formattedOP, fmt.Sprintf("LD DT, V%s", x)}
		case 0x001E:
			return &Line{formattedOP, fmt.Sprintf("ADD I, V%s", x)}
		case 0x0029:
			return &Line{formattedOP, fmt.Sprintf("LD F, V%s", x)}
		case 0x0033:
			return &Line{formattedOP, fmt.Sprintf("LD B, V%s", x)}
		case 0x0055:
			return &Line{formattedOP, fmt.Sprintf("LD I, V%s", x)}
		case 0x0065:
			return &Line{formattedOP, fmt.Sprintf("LD V%s, I", x)}
		}
	default:
		return nil
	}

	return nil
}

func (d *Disassembler) WriteToFile(fileName string) error {
	if !d.hasDisassembled {
		d.Disassemble()
	}

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
