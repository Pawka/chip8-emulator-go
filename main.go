package main

import (
	"fmt"
	"os"

	"github.com/Pawka/chip8-emulator/chip8"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <path-to-rom>\n", os.Args[0])
		os.Exit(2)
	}
	cpu := chip8.NewChip8()
	cpu.Load(os.Args[1])
}
