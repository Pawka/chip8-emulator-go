package chip8

import (
	"encoding/binary"
	"fmt"

	"github.com/Pawka/chip8-emulator/chip8/display"
)

// Chip8 is and interface of CHIP-8 emulator.
type Chip8 interface {
	// Run executes provided rom.
	Run(ctx Ctx)
}

type chip8 struct {
	display display.Display

	ram *ram
	// v is vector of registers. CHIP-8 has 16 8-bit data registers named V0 to
	// VF.
	v []byte

	// 16 level stack
	stack []uint16

	delayTimer int
	soundTimer int

	// program counter
	pc uint16
	// 16bit register (For memory address)
	i int
}

const registersCount = 16
const stackSize = 16
const timerInitialValue = 60

// NewChip8 creates a new instance of emulator.
func NewChip8() Chip8 {
	c := &chip8{
		display:    display.New(),
		ram:        newRAM(),
		v:          make([]byte, registersCount),
		stack:      make([]uint16, 0, stackSize),
		delayTimer: timerInitialValue,
		soundTimer: timerInitialValue,
		pc:         0x200,
	}

	return c
}

// Run implements the interface
func (c *chip8) Run(ctx Ctx) {
	c.ram.Load(ctx.path)

	if ctx.disassemble {
		for pc := 0x200; pc < memorySize; pc = pc + 2 {
			c.disassemble(pc)
		}
		return
	}

	c.pc = 0x200
	for {
		c.exec(c.pc)
	}
}

func (c *chip8) exec(pc uint16) {
	code := binary.BigEndian.Uint16(c.ram.Memory[pc : pc+2])
	first := code & 0xF000 >> 12

	switch first {
	case 0x0:
		switch code & 0x00FF {
		case 0xE0:
			c.pc += 2
			c.display.Clear()
		case 0xEE:
			c.pc = c.stack[len(c.stack)-1]
			c.stack = c.stack[:len(c.stack)-1]
		default:
			panic("Not implemented")
		}
	case 0x1:
		addr := code & 0x0FFF
		c.pc = addr
	case 0x2:
		if len(c.stack) == stackSize-1 {
			panic("Stack overflow")
		}
		addr := code & 0x0FFF
		c.stack = append(c.stack, c.pc+2)
		c.pc = addr
	default:
		panic("Not implemented")
	}
}

func (c *chip8) disassemble(pc int) {
	code := binary.BigEndian.Uint16(c.ram.Memory[pc : pc+2])
	first := code & 0xF000 >> 12

	var expl string
	switch first {
	case 0x0:
		switch code & 0x00FF {
		case 0xEE:
			expl = "return;"
		case 0xE0:
			expl = "disp_clear();"
		}
	case 0x1:
		addr := code & 0x0FFF
		expl = fmt.Sprintf("JMP #%x", addr)
	case 0x2:
		addr := code & 0x0FFF
		expl = fmt.Sprintf("CALL #%x", addr)
	case 0x3:
		vx := code & 0x0F00 >> 8
		nn := code & 0x00FF
		expl = fmt.Sprintf("SE V%X, %X", vx, nn)
	case 0x4:
		vx := code & 0x0F00 >> 8
		nn := code & 0x00FF
		expl = fmt.Sprintf("SNE V%X, %X", vx, nn)
	case 0x5:
		vx := code & 0x0F00 >> 8
		vy := code & 0x00F0 >> 4
		expl = fmt.Sprintf("SE V%X, V%X", vx, vy)
	case 0x6:
		vx := code & 0x0F00 >> 8
		nn := code & 0x00FF
		expl = fmt.Sprintf("LD V%X, %X", vx, nn)
	case 0x7:
		vx := code & 0x0F00 >> 8
		nn := code & 0x00FF
		expl = fmt.Sprintf("ADD V%X, %X", vx, nn)
	case 0x8:
		vx := code & 0x0F00 >> 8
		vy := code & 0x00F0 >> 4
		last := code & 0x000F
		switch last {
		case 0x0:
			expl = fmt.Sprintf("SUB V%X, V%X", vx, vy)
		case 0x6:
			expl = fmt.Sprintf("SHR V%X, V%X", vx, vy)
		case 0x7:
			expl = fmt.Sprintf("SUBN V%X, V%X", vx, vy)
		case 0xE:
			expl = fmt.Sprintf("SHL V%X, V%X", vx, vy)
		default:
			expl = fmt.Sprintf("> %X", first)
		}
	case 0x9:
		vx := code & 0x0F00 >> 8
		vy := code & 0x00F0 >> 4
		expl = fmt.Sprintf("SNE V%X, V%X", vx, vy)
	case 0xA:
		addr := code & 0x0FFF
		expl = fmt.Sprintf("LD I, #%x", addr)
	case 0xB:
		addr := code & 0x0FFF
		expl = fmt.Sprintf("JMP V0, #%x", addr)
	case 0xC:
		vx := code & 0x0F00 >> 8
		nn := code & 0x00FF
		expl = fmt.Sprintf("RND V%X, %X", vx, nn)
	case 0xD:
		vx := code & 0x0F00 >> 8
		vy := code & 0x00F0 >> 4
		n := code & 0x000F
		expl = fmt.Sprintf("DRV V%X, V%X, %X", vx, vy, n)
	case 0xE:
		vx := code & 0x0F00 >> 8
		last := code & 0x00FF
		switch last {
		case 0x9E:
			expl = fmt.Sprintf("SKP V%X", vx)
		case 0xA1:
			expl = fmt.Sprintf("SKPN V%X", vx)
		}
	case 0xF:
		vx := code & 0x0F00 >> 8
		last := code & 0x00FF
		switch last {
		case 0x07:
			expl = fmt.Sprintf("LD V%X, DT", vx)
		case 0x0A:
			expl = fmt.Sprintf("LD V%X, KEY", vx)
		case 0x15:
			expl = fmt.Sprintf("LD DT, V%X", vx)
		case 0x18:
			expl = fmt.Sprintf("LD ST V%X", vx)
		case 0x1E:
			expl = fmt.Sprintf("ADD I, V%X", vx)
		case 0x29:
			expl = fmt.Sprintf("LD I, FONT(V%X)", vx)
		case 0x33:
			expl = fmt.Sprintf("BCD V%X", vx)
		case 0x55:
			expl = fmt.Sprintf("LD [I], V%X", vx)
		case 0x65:
			expl = fmt.Sprintf("LD V%X, [I]", vx)
		}

	default:
		expl = fmt.Sprintf("> %X", first)
	}

	fmt.Printf("%04x\t%04X\t%s\n", pc, code, expl)
}
