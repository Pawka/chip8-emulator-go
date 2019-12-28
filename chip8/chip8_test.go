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
		opcode uint16
		setup  func(ch *chip8)
		assert func(t *testing.T, ch *chip8)
	}{
		// 00E0
		"clear_display": {
			opcode: 0x00E0,
			setup: func(ch *chip8) {
				ch.display = &displayMock{}
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.True(t, ch.display.(*displayMock).clear)
				assert.Equal(t, uint16(0x202), ch.pc)
			},
		},
		// 00EE
		"return_from_subroutine": {
			opcode: 0x00EE,
			setup: func(ch *chip8) {
				ch.stack = append(ch.stack, 0x0260)
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.Len(t, ch.stack, 0)
				assert.Equal(t, uint16(0x260), ch.pc)
			},
		},
		// 1NNN
		"jmp": {
			opcode: 0x12EE,
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, uint16(0x2EE), ch.pc)
			},
		},
		// 2NNN
		"call_subroutine": {
			opcode: 0x2AEE,
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, 1, len(ch.stack))
				assert.Equal(t, uint16(0x202), ch.stack[0])
				assert.Equal(t, uint16(0xAEE), ch.pc)
			},
		},
		// 3XNN
		"skip_if_equal_condition_must_be_skipped_because_values_are_equal": {
			opcode: 0x3510,
			setup: func(ch *chip8) {
				ch.v[5] = 0x10
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, uint16(0x204), ch.pc)
			},
		},
		"skip_if_equal_condition_must_not_be_skipped_because_values_are_not_equal": {
			opcode: 0x3510,
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, uint16(0x202), ch.pc)
			},
		},
		// 4XNN
		"skip_if_not_equal_condition_must_be_skipped_because_values_are_not_equal": {
			opcode: 0x4510,
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, uint16(0x204), ch.pc)
			},
		},
		"skip_if_not_equal_condition_must_not_be_skipped_because_values_are_equal": {
			opcode: 0x4510,
			setup: func(ch *chip8) {
				ch.v[5] = 0x10
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, uint16(0x202), ch.pc)
			},
		},
		// 5XY0
		"skip_if_registers_equal_must_be_skipped_because_values_are_equal": {
			opcode: 0x5540,
			setup: func(ch *chip8) {
				ch.v[4] = 0x10
				ch.v[5] = 0x10
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, uint16(0x204), ch.pc)
			},
		},
		"skip_if_registers_equal_must_not_be_skipped_because_values_are_not_equal": {
			opcode: 0x5540,
			setup: func(ch *chip8) {
				ch.v[5] = 0x10
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, uint16(0x202), ch.pc)
			},
		},
	}
	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			chip8 := NewChip8().(*chip8)
			a := (test.opcode & 0xFF00) >> 8
			b := test.opcode & 0x00FF
			chip8.ram.Memory[0x200] = uint8(a)
			chip8.ram.Memory[0x200+1] = uint8(b)
			if test.setup != nil {
				test.setup(chip8)
			}
			chip8.exec(programStartPos)
			test.assert(t, chip8)
		})
	}
}
