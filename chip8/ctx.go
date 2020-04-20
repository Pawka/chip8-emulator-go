package chip8

import (
	"errors"
	"flag"
)

// Ctx is context of the program which holds command line arguments.
type Ctx struct {
	disassemble bool
	path        string
}

// IsDisplay returns true if display is supposed to be created.
func (c Ctx) IsDisplay() bool {
	return !c.disassemble
}

const (
	_ = iota

	// PathArgPosition holds position of path argument
	PathArgPosition
)

// NewCtxFromArgs initializes a new context.
func NewCtxFromArgs(args []string) (Ctx, error) {
	ctx := Ctx{}
	set := flag.NewFlagSet(args[0], flag.ExitOnError)
	set.BoolVar(&ctx.disassemble, "d", false, "Run disassembler for given program")
	set.Parse(args[1:])

	if set.NArg() < PathArgPosition {
		return ctx, errors.New("provide path to program")
	}

	ctx.path = set.Arg(PathArgPosition - 1)

	return ctx, nil
}
