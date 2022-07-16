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

func (d *Disassembler) Disassemble() []Line {
	lines := []Line{}

	for d.pc != uint16(len(d.fileContent)-1) {
		op := uint16(d.fileContent[d.pc])<<8 | uint16(d.fileContent[d.pc+1])

		x := (op & 0x0F00) >> 8
		y := (op & 0x00F0) >> 4
		nnn := op & 0x0FFF
		kk := op & 0x00FF

		switch op & 0xF000 {
		case 0x0000:
			switch op & 0x000F {
			case 0x0000:
				lines = append(lines, Line{formatOP(op), "CLS"})
			case 0x000E:
				lines = append(lines, Line{formatOP(op), "RET"})
			default:
				break
			}
		case 0x1000:
			lines = append(lines, Line{formatOP(op), fmt.Sprintf("JP 0x%04X", nnn)})
		case 0x2000:
			lines = append(lines, Line{formatOP(op), fmt.Sprintf("CALL 0x%04X", nnn)})
		case 0x3000:
			lines = append(lines, Line{formatOP(op), fmt.Sprintf("SE V%X, 0x%04X", x, kk)})
		case 0x4000:
			lines = append(lines, Line{formatOP(op), fmt.Sprintf("SNE V%X, 0x%04X", x, kk)})
		case 0x5000:
			lines = append(lines, Line{formatOP(op), fmt.Sprintf("SE V%X, V%X", x, y)})
		case 0x6000:
			lines = append(lines, Line{formatOP(op), fmt.Sprintf("LD V%X, 0x%04X", x, kk)})
		case 0x7000:
			lines = append(lines, Line{formatOP(op), fmt.Sprintf("ADD V%X, 0x%04X", x, kk)})
		case 0x8000:
			switch op & 0x000F {
			case 0x0000:
				lines = append(lines, Line{formatOP(op), fmt.Sprintf("LD V%X, V%X", x, y)})
			case 0x0001:
				lines = append(lines, Line{formatOP(op), fmt.Sprintf("OR V%X, V%X", x, y)})
			case 0x0002:
				lines = append(lines, Line{formatOP(op), fmt.Sprintf("AND V%X, V%X", x, y)})
			case 0x0003:
				lines = append(lines, Line{formatOP(op), fmt.Sprintf("XOR V%X, V%X", x, y)})
			case 0x0004:
				lines = append(lines, Line{formatOP(op), fmt.Sprintf("ADD V%X, V%X", x, y)})
			case 0x0005:
				lines = append(lines, Line{formatOP(op), fmt.Sprintf("SUB V%X, V%X", x, y)})
			case 0x0006:
				lines = append(lines, Line{formatOP(op), fmt.Sprintf("SHR V%X, V%X", x, y)})
			case 0x0007:
				lines = append(lines, Line{formatOP(op), fmt.Sprintf("SUBN V%X, V%X", x, y)})
			case 0x000E:
				lines = append(lines, Line{formatOP(op), fmt.Sprintf("SHL V%X, V%X", x, y)})
			default:
				fmt.Printf("invalid opcode: %X\n", op)
			}
		case 0x9000:
			lines = append(lines, Line{formatOP(op), fmt.Sprintf("SNE V%X, V%X", x, y)})
		case 0xA000:
			lines = append(lines, Line{formatOP(op), fmt.Sprintf("LD I, 0x%04X", nnn)})
		case 0xB000:
			lines = append(lines, Line{formatOP(op), fmt.Sprintf("JP V0, 0x%04X", nnn)})
		case 0xC000:
			lines = append(lines, Line{formatOP(op), fmt.Sprintf("RND V%X, 0x%04X", x, kk)})
		case 0xD000:
			lines = append(lines, Line{formatOP(op), fmt.Sprintf("DRW V%X, V%X, 0x%04X", x, y, op&0x000F)})
		case 0xE000:
			switch op & 0x00FF {
			case 0x009E:
				lines = append(lines, Line{formatOP(op), fmt.Sprintf("SKP V%X", x)})
			case 0x00A1:
				lines = append(lines, Line{formatOP(op), fmt.Sprintf("SKNP V%X", x)})
			default:
				fmt.Printf("invalid opcode: %X\n", op)
			}
		case 0xF000:
			switch op & 0x00FF {
			case 0x0007:
				lines = append(lines, Line{formatOP(op), fmt.Sprintf("LD V%X, DT", x)})
			case 0x000A:
				lines = append(lines, Line{formatOP(op), fmt.Sprintf("LD V%X, K", x)})
			case 0x0015:
				lines = append(lines, Line{formatOP(op), fmt.Sprintf("LD DT, V%X", x)})
			case 0x0018:
				lines = append(lines, Line{formatOP(op), fmt.Sprintf("LD DT, V%X", x)})
			case 0x001E:
				lines = append(lines, Line{formatOP(op), fmt.Sprintf("ADD I, V%X", x)})
			case 0x0029:
				lines = append(lines, Line{formatOP(op), fmt.Sprintf("LD F, V%X", x)})
			case 0x0033:
				lines = append(lines, Line{formatOP(op), fmt.Sprintf("LD B, V%X", x)})
			case 0x0055:
				lines = append(lines, Line{formatOP(op), fmt.Sprintf("LD [I], V%X", x)})
			case 0x0065:
				lines = append(lines, Line{formatOP(op), fmt.Sprintf("LD V%X, [I]", x)})
			}
		default:
			fmt.Printf("invalid opcode: %X\n", op)
		}

		d.pc += 2
	}

	return lines
}

// Format an opcode to its hexadecimal representation (0x0000)
func formatOP(op uint16) string {
	return fmt.Sprintf("0x%04X", op)
}
