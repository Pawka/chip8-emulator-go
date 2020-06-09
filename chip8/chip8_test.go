package chip8

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type displayMock struct {
	clear, show, point, sprite, pollkey bool
	x, y                                int
	payload                             []byte
	// Keys to be returned by PollKey method
	keys []*rune
}

func (d *displayMock) Show() {
	d.show = true
}

func (d *displayMock) Point(x int, y int) bool {
	d.point = true
	d.x = x
	d.y = y
	return false
}

func (d *displayMock) Sprite(x int, y int, payload []byte) bool {
	d.sprite = true
	d.x = x
	d.y = y
	d.payload = payload
	return true
}

func (d *displayMock) Clear() {
	d.clear = true
}

func (d *displayMock) Debug(line string) {

}

func (d *displayMock) PollKey() *rune {
	if len(d.keys) == 0 {
		return nil
	}
	val := d.keys[0]
	d.keys = d.keys[1:len(d.keys)]
	return val
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
		// 6XNN
		"set_nn_value_to_x": {
			opcode: 0x6540,
			setup: func(ch *chip8) {
				ch.v[5] = 0x10
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, uint8(0x40), ch.v[5])
				assert.Equal(t, uint16(0x202), ch.pc)
			},
		},
		// 7XNN
		"add_nn_value_to_x": {
			opcode: 0x7501,
			setup: func(ch *chip8) {
				ch.v[5] = 0x10
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, uint8(0x11), ch.v[5])
				assert.Equal(t, uint16(0x202), ch.pc)
			},
		},
		"add_nn_value_to_x_with_overflow": {
			opcode: 0x7501,
			setup: func(ch *chip8) {
				ch.v[5] = 0xFF
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, uint8(0x0), ch.v[5])
				assert.Equal(t, uint8(0x0), ch.v[0xF])
				assert.Equal(t, uint16(0x202), ch.pc)
			},
		},
		// 8XY0
		"copy_y_value_to_x": {
			opcode: 0x8230,
			setup: func(ch *chip8) {
				ch.v[2] = 0x10
				ch.v[3] = 0x20
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, uint8(0x20), ch.v[2])
				assert.Equal(t, ch.v[2], ch.v[3])
				assert.Equal(t, uint16(0x202), ch.pc)
			},
		},
		"set_vx_or_vy_to_vx": {
			opcode: 0x8231,
			setup: func(ch *chip8) {
				ch.v[2] = 0x12
				ch.v[3] = 0x3
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, uint8(0x13), ch.v[2])
				assert.Equal(t, uint16(0x202), ch.pc)
			},
		},
		"set_vx_and_vy_to_vx": {
			opcode: 0x8232,
			setup: func(ch *chip8) {
				ch.v[2] = 0x5
				ch.v[3] = 0x6
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, uint8(0x4), ch.v[2])
				assert.Equal(t, uint16(0x202), ch.pc)
			},
		},
		"set_vx_xor_vy_to_vx": {
			opcode: 0x8233,
			setup: func(ch *chip8) {
				ch.v[2] = 0x12
				ch.v[3] = 0x3
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, uint8(0x11), ch.v[2])
				assert.Equal(t, uint16(0x202), ch.pc)
			},
		},
		"add_nn_value_to_x_with_overflow_flag_set_to_0": {
			opcode: 0x8234,
			setup: func(ch *chip8) {
				ch.v[2] = 0x1
				ch.v[3] = 0x2
				ch.v[0xF] = 0xFF
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, uint8(0x3), ch.v[2])
				assert.Equal(t, uint8(0x0), ch.v[0xF])
				assert.Equal(t, uint16(0x202), ch.pc)
			},
		},
		"add_nn_value_to_x_with_overflow_flag_set_to_1": {
			opcode: 0x8234,
			setup: func(ch *chip8) {
				ch.v[2] = 0xFF
				ch.v[3] = 0x02
				ch.v[0xF] = 0xFF
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, uint8(0x1), ch.v[2])
				assert.Equal(t, uint8(0x1), ch.v[0xF])
				assert.Equal(t, uint16(0x202), ch.pc)
			},
		},
		"substract_nn_value_from_x_with_overflow_flag_set_to_0": {
			opcode: 0x8235,
			setup: func(ch *chip8) {
				ch.v[2] = 0x2
				ch.v[3] = 0x1
				ch.v[0xF] = 0xFF
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, uint8(0x1), ch.v[2])
				assert.Equal(t, uint8(0x0), ch.v[0xF])
				assert.Equal(t, uint16(0x202), ch.pc)
			},
		},
		"substract_nn_value_from_x_with_overflow_flag_set_to_1": {
			opcode: 0x8235,
			setup: func(ch *chip8) {
				ch.v[2] = 0x1
				ch.v[3] = 0x2
				ch.v[0xF] = 0xFF
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, uint8(0xFF), ch.v[2])
				assert.Equal(t, uint8(0x1), ch.v[0xF])
				assert.Equal(t, uint16(0x202), ch.pc)
			},
		},
		"right_shift_with_vf_set_to_0": {
			opcode: 0x8236,
			setup: func(ch *chip8) {
				ch.v[2] = 0x4
				ch.v[3] = 0x5
				ch.v[0xF] = 0xFF
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, uint8(0x2), ch.v[2])
				assert.Equal(t, uint8(0x0), ch.v[0xF])
				assert.Equal(t, uint16(0x202), ch.pc)
			},
		},
		"right_shift_with_vf_set_to_1": {
			opcode: 0x8236,
			setup: func(ch *chip8) {
				ch.v[2] = 0x5
				ch.v[3] = 0x5
				ch.v[0xF] = 0xFF
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, uint8(0x2), ch.v[2])
				assert.Equal(t, uint8(0x1), ch.v[0xF])
				assert.Equal(t, uint16(0x202), ch.pc)
			},
		},
		"set_vx_equal_vy_minus_xy_when_vx_is_less_than_vy": {
			opcode: 0x8237,
			setup: func(ch *chip8) {
				ch.v[2] = 0x3
				ch.v[3] = 0x4
				ch.v[0xF] = 0xFF
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, uint8(0x1), ch.v[2])
				assert.Equal(t, uint8(0x1), ch.v[0xF])
				assert.Equal(t, uint16(0x202), ch.pc)
			},
		},
		"set_vx_equal_vy_minus_xy_when_vx_is_not_less_than_vy": {
			opcode: 0x8237,
			setup: func(ch *chip8) {
				ch.v[2] = 0x4
				ch.v[3] = 0x4
				ch.v[0xF] = 0xFF
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, uint8(0x0), ch.v[2])
				assert.Equal(t, uint8(0x0), ch.v[0xF])
				assert.Equal(t, uint16(0x202), ch.pc)
			},
		},
		"left_shift_with_vf_set_to_0": {
			opcode: 0x823E,
			setup: func(ch *chip8) {
				ch.v[2] = 0x0
				ch.v[3] = 0x4
				ch.v[0xF] = 0xFF
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, uint8(0x8), ch.v[2])
				assert.Equal(t, uint8(0x0), ch.v[0xF])
				assert.Equal(t, uint16(0x202), ch.pc)
			},
		},
		"left_shift_with_vf_set_to_1": {
			opcode: 0x823E,
			setup: func(ch *chip8) {
				ch.v[2] = 0xFF
				ch.v[3] = 0x8
				ch.v[0xF] = 0xFF
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, uint8(0x10), ch.v[2])
				assert.Equal(t, uint8(0x1), ch.v[0xF])
				assert.Equal(t, uint16(0x202), ch.pc)
			},
		},
		// 9XY0
		"skip_next_function_when_x_not_equal_y": {
			opcode: 0x9120,
			setup: func(ch *chip8) {
				ch.v[1] = 0x1
				ch.v[2] = 0x2
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, uint16(0x204), ch.pc)
			},
		},
		"not_skip_next_function_when_x_equal_y": {
			opcode: 0x9110,
			setup: func(ch *chip8) {
				ch.v[1] = 0x1
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, uint16(0x202), ch.pc)
			},
		},
		// ANNN
		"set_i_to_nnn": {
			opcode: 0xA123,
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, 0x123, ch.i)
			},
		},
		// BNNN
		"jump_to_nnn+v0": {
			opcode: 0xB123,
			setup: func(ch *chip8) {
				ch.v[0] = 0x2
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, uint16(0x125), ch.pc)
			},
		},
		// CXNN
		"set_vx_to_number_and_bitwise_0x0F": {
			opcode: 0xC30F,
			assert: func(t *testing.T, ch *chip8) {
				assert.GreaterOrEqual(t, ch.v[3], byte(0))
				assert.LessOrEqual(t, ch.v[3], byte(0xF))
				assert.Equal(t, uint16(0x202), ch.pc)
			},
		},
		// DXYN
		"draw_a_sprite": {
			opcode: 0xD123,
			setup: func(ch *chip8) {
				ch.display = &displayMock{}
				ch.v[1] = 10
				ch.v[2] = 20
				ch.i = 0x300
				ch.ram.Memory[0x300] = 0x2
				ch.ram.Memory[0x301] = 0x3
				ch.ram.Memory[0x302] = 0x4
				ch.ram.Memory[0x303] = 0x5
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.True(t, ch.display.(*displayMock).sprite)
				assert.Equal(t, byte(1), ch.v[0xF])
				assert.Equal(t, 10, ch.display.(*displayMock).x)
				assert.Equal(t, 20, ch.display.(*displayMock).y)
				assert.Equal(t, []byte{0x2, 0x3, 0x4}, ch.display.(*displayMock).payload)
				assert.Equal(t, uint16(0x202), ch.pc)
			},
		},
		// EX9E
		"skip_instruction_when_key_is_pressed": {
			opcode: 0xEA9E,
			setup: func(ch *chip8) {
				key := 'd' // 'd' is mapped to '9' on keyboard
				ch.display = &displayMock{
					keys: []*rune{&key},
				}
				ch.v[0xA] = 9
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, uint16(0x204), ch.pc)
			},
		},
		"do_not_skip_instruction_when_other_key_is_pressed": {
			opcode: 0xEA9E,
			setup: func(ch *chip8) {
				key1 := 'd' // 'd' is mapped to '9' on keyboard
				key2 := '1'
				ch.display = &displayMock{
					keys: []*rune{&key1, &key2},
				}
				ch.v[0xA] = 2
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, uint16(0x202), ch.pc)
			},
		},
		// EXA1
		"skip_instruction_when_key_not_is_pressed": {
			opcode: 0xEAA1,
			setup: func(ch *chip8) {
				key := '9'
				ch.display = &displayMock{
					keys: []*rune{&key},
				}
				ch.v[0xA] = 8
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, uint16(0x204), ch.pc)
			},
		},
		"do_not_skip_instruction_when_key_is_pressed": {
			opcode: 0xEAA1,
			setup: func(ch *chip8) {
				key1 := 'd' // 'd' is mapped to '9' on keyboard
				key2 := '1'
				ch.display = &displayMock{
					keys: []*rune{&key1, &key2},
				}
				ch.v[0xA] = 9
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, uint16(0x202), ch.pc)
			},
		},
		// FX07
		"set_delay_timer_value_to_vx": {
			opcode: 0xF407,
			setup: func(ch *chip8) {
				ch.delayTimer = 0x10
				ch.v[0x4] = 123
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, uint8(0x10), ch.v[0x4])
				assert.Equal(t, uint16(0x202), ch.pc)
			},
		},
		// FX0A
		"set_pressed_key_to_vx": {
			opcode: 0xF40A,
			setup: func(ch *chip8) {
				key1 := 'd' // 'd' is mapped to '9' on keyboard
				ch.display = &displayMock{
					keys: []*rune{&key1},
				}
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, uint8(0x9), ch.v[0x4])
				assert.Equal(t, uint16(0x202), ch.pc)
			},
		},
		// FX15
		"set_delay_timer_to_vx": {
			opcode: 0xF415,
			setup: func(ch *chip8) {
				ch.v[0x4] = 0x10
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, uint8(0x10), ch.delayTimer)
				assert.Equal(t, uint16(0x202), ch.pc)
			},
		},
		// FX18
		"set_sound_timer_to_vx": {
			opcode: 0xF418,
			setup: func(ch *chip8) {
				ch.v[0x4] = 0x10
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, uint8(0x10), ch.soundTimer)
				assert.Equal(t, uint16(0x202), ch.pc)
			},
		},
		// FX1E
		"add_vx_to_i": {
			opcode: 0xF41E,
			setup: func(ch *chip8) {
				ch.i = 10
				ch.v[0x4] = 0x1
				ch.v[0xF] = 0x10
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, 11, ch.i)
				assert.Equal(t, uint8(0), ch.v[0xF])
				assert.Equal(t, uint16(0x202), ch.pc)
			},
		},
		// FX1E
		"add_vx_to_i_with_overflow": {
			opcode: 0xF41E,
			setup: func(ch *chip8) {
				ch.i = memorySize - 2
				ch.v[0x4] = 0x9
				ch.v[0xF] = 0x10
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, memorySize+7, ch.i)
				assert.Equal(t, uint8(1), ch.v[0xF])
				assert.Equal(t, uint16(0x202), ch.pc)
			},
		},
		// FX29
		"point_i_to_font_position": {
			opcode: 0xF429,
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, uint16(0x202), ch.pc)
				assert.Equal(t, 40, ch.i)
			},
		},
		// FX33
		"store_decimal_at_vx_to_i": {
			opcode: 0xF433,
			setup: func(ch *chip8) {
				ch.i = 2
				ch.v[4] = 123
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, uint16(0x202), ch.pc)
				assert.Equal(t, uint8(1), ch.ram.Memory[ch.i])
				assert.Equal(t, uint8(2), ch.ram.Memory[ch.i+1])
				assert.Equal(t, uint8(3), ch.ram.Memory[ch.i+2])
			},
		},
		// FX55
		"fill_memory_with_values_from_v0_to_vx": {
			opcode: 0xF455,
			setup: func(ch *chip8) {
				ch.i = 2
				ch.v[0x0] = 0x10
				ch.v[0x1] = 0x11
				ch.v[0x2] = 0x12
				ch.v[0x3] = 0x13
				ch.v[0x4] = 0x14
				ch.v[0x5] = 0x15
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, []byte{0, 0, 0x10, 0x11, 0x12, 0x13, 0x14, 0}, ch.ram.Memory[0:8])
				assert.Equal(t, uint16(0x202), ch.pc)
				assert.Equal(t, 2, ch.i)
			},
		},
		// FX65
		"fill_v0_to_vx_registers_from_memory": {
			opcode: 0xF465,
			setup: func(ch *chip8) {
				ch.i = 2
				data := []byte{0, 0, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15}
				copy(ch.ram.Memory, data)
			},
			assert: func(t *testing.T, ch *chip8) {
				assert.Equal(t, 2, ch.i)
				assert.Equal(t, uint16(0x202), ch.pc)

				for i := 0; i <= 4; i++ {
					assert.Equal(t, ch.v[i], uint8(0x10+i))
				}
				assert.Equal(t, ch.v[0x5], uint8(0))
			},
		},
	}
	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			ctx := Ctx{}
			chip8 := NewChip8(ctx).(*chip8)
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
