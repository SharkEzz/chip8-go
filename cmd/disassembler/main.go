package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/SharkEzz/chip8-go/pkg/disassembler"
)

func main() {
	file := flag.String("file", "", "The file to disassemble")

	flag.Parse()

	d, err := disassembler.NewDisassembler(*file)
	if err != nil {
		panic(err)
	}

	lines := d.Disassemble()
	data, _ := json.MarshalIndent(lines, "", "  ")
	fmt.Println(string(data))
}
