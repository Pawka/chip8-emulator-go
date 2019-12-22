package chip8

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type displayMock struct {
	clear bool
}

func (d *displayMock) Clear() {
	d.clear = true
}

func TestExec(t *testing.T) {
	testCases := map[string]struct {
		opcode []byte
		setup  func(ch *chip8)
		assert func(t *testing.T, ch *chip8)
	}{
		// 00E0
		"clear_display": {
			opcode: []byte{0x00, 0xE0},
			setup: func(ch *chip8) {
				ch.display = &displayMock{}
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.True(t, ch.display.(*displayMock).clear)
				assert.Equal(t, uint16(0x202), ch.pc)
			},
		},
		// 1NNN
		"jmp": {
			opcode: []byte{0x12, 0xEE},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, uint16(0x2EE), ch.pc)
			},
		},
		// 2NNN
		"call_subroutine": {
			opcode: []byte{0x2A, 0xEE},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, 1, len(ch.stack))
				assert.Equal(t, uint16(0x202), ch.stack[0])
				assert.Equal(t, uint16(0xAEE), ch.pc)
			},
		},
	}
	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			chip8 := NewChip8().(*chip8)
			chip8.ram.Memory[0x200] = test.opcode[0]
			chip8.ram.Memory[0x200+1] = test.opcode[1]
			if test.setup != nil {
				test.setup(chip8)
			}
			chip8.exec(programStartPos)
			test.assert(t, chip8)
		})
	}
}
