package display

const (
	width  = 64
	height = 32
)

// Display defines interface of CHIP8 display.
type Display interface {
	// Clear the screen.
	Clear()
}

type display struct {
	screen [][]int
}

// New initializes a new display
func New() Display {
	d := &display{}
	d.init()
	return d
}

func (d *display) init() {
	d.screen = make([][]int, height)
	for i := range d.screen {
		d.screen[i] = make([]int, width)
	}
}

func (d *display) Clear() {
	d.init()
}
