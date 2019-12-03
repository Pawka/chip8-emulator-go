package chip8

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
	}

	return c
}

// Run implements the interface
func (c *chip8) Run(path string) {
	c.ram.Load(path)
}
