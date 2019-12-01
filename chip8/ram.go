package chip8

import (
	"fmt"
	"io/ioutil"
)

const memorySize = 4098

type ram struct {
	Memory []byte
}

func newRAM() *ram {
	return &ram{
		Memory: make([]byte, memorySize),
	}
}

// Load program to memory.
func (r *ram) Load(path string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed load rom at path %q: %s", path, err)
	}
	for k, v := range b {
		r.Memory[0x200+k] = v
		fmt.Printf("%X\n", v)
	}

	return nil
}
