package chip8

import (
	"encoding/binary"
	"flag"
	"fmt"
	"math/rand"
	"time"

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
func NewChip8(ctx Ctx) Chip8 {
	var d display.Display

	// Do not initialize a new display during test run.
	// Creating it breaks test output.
	if ctx.IsDisplay() == true && flag.Lookup("test.v") == nil {
		var err error
		d, err = display.New()
		if err != nil {
			// TODO: Return as error
			panic(err)
		}
	}

	c := &chip8{
		display:    d,
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

	quit := make(chan struct{})
	if ctx.IsDisplay() {
		go func() {
			c.display.Show()
			close(quit)
		}()
	} else {
		close(quit)
	}

	go func() {
		// INVADER
		//
		// X.XXX.X.     0b10111010  $BA
		// .XXXXX..     0b01111100  $7C
		// XX.X.XX.     0b11010110  $D6
		// XXXXXXX.     0b11111110  $FE
		// .X.X.X..     0b01010100  $54
		// X.X.X.X.     0b10101010  $AA
		invader := []byte{0xBA, 0x7C, 0xD6, 0xFE, 0x54, 0xAA}
		for {
			for i := -7; i < 40; i += 3 {
				c.display.Sprite(i*2, i, invader)
				time.Sleep(time.Millisecond * 100)
				c.display.Clear()
			}
		}
	}()

	c.pc = 0x200
	//for {
	//	c.exec(c.pc)
	//}

	<-quit
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
	case 0x3:
		vx := code & 0x0F00 >> 8
		nn := code & 0x00FF
		c.pc += 2
		if uint8(nn) == c.v[vx] {
			c.pc += 2
		}
	case 0x4:
		vx := code & 0x0F00 >> 8
		nn := code & 0x00FF
		c.pc += 2
		if uint8(nn) != c.v[vx] {
			c.pc += 2
		}
	case 0x5:
		vx := code & 0x0F00 >> 8
		vy := code & 0x00F0 >> 4
		c.pc += 2
		if c.v[vy] == c.v[vx] {
			c.pc += 2
		}
	case 0x6:
		vx := code & 0x0F00 >> 8
		nn := code & 0x00FF
		c.v[vx] = uint8(nn)
		c.pc += 2
	case 0x7:
		vx := code & 0x0F00 >> 8
		nn := code & 0x00FF
		c.v[vx] += uint8(nn)
		c.pc += 2
	case 0x8:
		vx := code & 0x0F00 >> 8
		vy := code & 0x00F0 >> 4
		last := code & 0x000F
		switch last {
		case 0x0:
			c.v[vx] = c.v[vy]
		case 0x1:
			c.v[vx] = c.v[vx] | c.v[vy]
		case 0x2:
			c.v[vx] = c.v[vx] & c.v[vy]
		case 0x3:
			c.v[vx] = c.v[vx] ^ c.v[vy]
		case 0x4:
			res := uint16(c.v[vx]) + uint16(c.v[vy])
			c.v[0xF] = 0
			if res > 0xFF {
				c.v[0xF] = 1
			}
			c.v[vx] = c.v[vx] + c.v[vy]
		case 0x5:
			c.v[0xF] = 0
			// NOTE: Not sure if 0 should enable underoverflow flag.
			if c.v[vx] < c.v[vy] {
				c.v[0xF] = 1
			}
			c.v[vx] = c.v[vx] - c.v[vy]
		case 0x6:
			c.v[0xF] = c.v[vx] & 0x1
			c.v[vx] = c.v[vy] >> 1
		case 0x7:
			c.v[0xF] = 0
			if c.v[vy] > c.v[vx] {
				c.v[0xF] = 1
			}
			c.v[vx] = c.v[vy] - c.v[vx]
		case 0xE:
			c.v[0xF] = c.v[vx] >> 7
			c.v[vx] = c.v[vy] << 1
		default:
			panic("Not implemented")
		}
		c.pc += 2
	case 0x9:
		vx := code & 0x0F00 >> 8
		vy := code & 0x00F0 >> 4
		if vx != vy {
			c.pc += 2
		}
		c.pc += 2
	case 0xA:
		addr := code & 0x0FFF
		c.i = int(addr)
		c.pc += 2
	case 0xB:
		addr := code & 0x0FFF
		c.pc = addr + uint16(c.v[0])
	case 0xC:
		vx := code & 0x0F00 >> 8
		last := code & 0x00FF
		s := rand.NewSource(time.Now().UnixNano())
		r := rand.New(s)
		c.v[vx] = byte(r.Intn(256)) & byte(last)
		c.pc += 2
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
