package main

import (
	"flag"
	"fmt"

	"github.com/SharkEzz/chip8-go/pkg/disassembler"
)

func main() {
	file := flag.String("file", "", "The file to disassemble")
	outputFile := flag.String("outputFile", "out.asm", "The file to which output the disassembled content")

	flag.Parse()

	if *file == "" {
		panic(fmt.Errorf("file cannot be empty"))
	}

	d, err := disassembler.NewDisassembler(*file)
	if err != nil {
		panic(err)
	}

	d.WriteToFile(*outputFile)
}
