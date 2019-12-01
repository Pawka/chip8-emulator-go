package chip8

// Chip8 is and interface of CHIP-8 emulator.
type Chip8 interface {
	// Run executes provided rom.
	Run(string)
}

type chip8 struct {
	ram *ram
}

// NewChip8 creates a new instance of emulator.
func NewChip8() Chip8 {
	c := &chip8{
		ram: newRAM(),
	}

	return c
}

// Run implements the interface
func (c *chip8) Run(path string) {
	c.ram.Load(path)
}
