package display

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell"
)

const (
	width  = 64
	height = 32
)

// Display defines interface of CHIP8 display.
type Display interface {
	// Clear the screen.
	Clear()
	// Show the screen.
	Show()
}

type display struct {
	s      tcell.Screen
	screen [][]int
}

// New initializes a new display
func New() (Display, error) {
	s, err := tcell.NewScreen()
	if err != nil {
		return nil, fmt.Errorf("creating screen: %v", err)
	}
	err = s.Init()
	if err != nil {
		return nil, fmt.Errorf("initializing screen: %v", err)
	}
	d := &display{
		s: s,
	}
	return d, nil
}

func (d *display) Show() {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	d.s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack))
	d.Clear()

	quit := make(chan struct{})
	go func() {
		for {
			ev := d.s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyCtrlC || ev.Key() == tcell.KeyEscape {
					close(quit)
					return
				}
			case *tcell.EventResize:
				d.Clear()
				d.s.Sync()
			}
		}
	}()

loop:
	for {
		select {
		case <-quit:
			break loop
		case <-time.After(time.Millisecond * 50):
		}

		d.DrawScreen(width, height)
		d.s.Show()
	}
	d.s.Fini()
}

func (d *display) SetContent(x int, y int, mainc rune, combc []rune, style tcell.Style) {
	dw, dh := d.s.Size()
	_y := dh/2 - height/2
	_x := dw/2 - width/2
	d.s.SetContent(_x+x, _y+y, mainc, combc, style)
}

func (d *display) DrawScreen(w, h int) {
	st := tcell.StyleDefault.
		Background(tcell.ColorWhite)
	d.SetContent(0, 0, ' ', nil, st)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			d.SetContent(x, y, ' ', nil, st)
		}
	}
}

func (d *display) Clear() {
	d.s.Clear()
}
