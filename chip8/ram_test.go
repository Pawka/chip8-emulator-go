package chip8

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const binaryPath = "testdata/binary"

func TestLoad(t *testing.T) {
	ram := newRAM()
	ram.Load(binaryPath)
	want := []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x0a}
	assert.Equal(t, want, ram.Memory[programStartPos:programStartPos+len(want)])
}
