package chip8

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCtxFromArgs(t *testing.T) {
	testCases := map[string]struct {
		args    []string
		want    Ctx
		wantErr string
	}{
		"happy_path": {
			args: []string{"program", "file"},
			want: Ctx{
				path: "file",
			},
		},
		"disassembler_flag_provided": {
			args: []string{"program", "-d", "file"},
			want: Ctx{
				disassemble: true,
				path:        "file",
			},
		},
		"no_program_path_provided": {
			args: []string{"program", "-d"},
			want: Ctx{
				disassemble: true,
			},
			wantErr: "provide path to program",
		},
	}
	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			ctx, err := NewCtxFromArgs(test.args)
			assert.Equal(t, test.want, ctx)
			if test.wantErr != "" {
				assert.EqualError(t, err, test.wantErr)
				return
			}

			assert.NoError(t, err)
		})
	}
}
