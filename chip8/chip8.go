package chip8

// Chip8 is and interface of CHIP-8 emulator.
type Chip8 interface {
	// Load the rom into memory.
	Load(string)
}

const memorySize = 4098

// NewChip8 creates a new instance of emulator.
func NewChip8() Chip8 {
	c := &chip8{}
	c.initializeRAM()

	return c
}

type chip8 struct {
	ram []byte
}

func (c *chip8) initializeRAM() {
	c.ram = make([]byte, memorySize)
}

// Load implements the interface
func (c *chip8) Load(path string) {

}
