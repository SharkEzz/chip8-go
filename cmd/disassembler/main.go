package main

import (
	"encoding/json"
	"fmt"

	"github.com/SharkEzz/chip8-go/pkg/disassembler"
)

func main() {
	d, err := disassembler.NewDisassembler("./stars.ch8")
	if err != nil {
		panic(err)
	}

	lines := d.Disassemble()
	data, _ := json.MarshalIndent(lines, "", "  ")
	fmt.Println(string(data))
}
