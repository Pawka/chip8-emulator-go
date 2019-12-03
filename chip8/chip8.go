package chip8

import (
	"encoding/binary"
	"fmt"
)

// Chip8 is and interface of CHIP-8 emulator.
type Chip8 interface {
	// Run executes provided rom.
	Run(string)
}

type chip8 struct {
	ram *ram
	// v is vector of registers. CHIP-8 has 16 8-bit data registers named V0 to
	// VF.
	v []byte

	// 16 level stack
	stack []int

	delayTimer int
	soundTimer int

	// program counter
	pc int
	// 16bit register (For memory address)
	i int
}

const registersCount = 16
const stackSize = 16
const timerInitialValue = 60

// NewChip8 creates a new instance of emulator.
func NewChip8() Chip8 {
	c := &chip8{
		ram:        newRAM(),
		v:          make([]byte, registersCount),
		stack:      make([]int, stackSize),
		delayTimer: timerInitialValue,
		soundTimer: timerInitialValue,
		pc:         0x200,
	}

	return c
}

// Run implements the interface
func (c *chip8) Run(path string) {
	c.ram.Load(path)

	for pc := 0x200; pc < memorySize; pc = pc + 2 {
		c.disassemble(pc)
	}
}

func (c *chip8) disassemble(pc int) {
	code := binary.BigEndian.Uint16(c.ram.Memory[pc : pc+2])

	first := code & 0xF000

	var expl string
	switch first {
	case 0x00:
		switch code & 0x00FF {
		case 0xEE:
			expl = "return;"
		case 0xE0:
			expl = "disp_clear();"
		}
	case 0x1000:
		addr := code & 0x0FFF
		expl = fmt.Sprintf("JMP #%x", addr)
	default:
		expl = fmt.Sprintf("%X", first)
	}

	fmt.Printf("%04x\t%04X\t%s\n", pc, code, expl)
}
