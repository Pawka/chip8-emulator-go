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
	// Point draws a point at x y
	Point(x, y int)
	// Sprite at x, y coordinates draws a sprite which is 8 symbols width and n
	// lines height.
	Sprite(x, y int, payload []byte)
}

type display struct {
	s                tcell.Screen
	screen           [][]int
	sprites          chan sprite
	bgStyle, fgStyle tcell.Style
}

type sprite struct {
	x, y    int
	payload []byte
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
	fg := tcell.StyleDefault.Background(tcell.ColorBlack)
	bg := tcell.StyleDefault.Background(tcell.ColorWhite)
	d := &display{
		s:       s,
		sprites: make(chan sprite),
		bgStyle: bg,
		fgStyle: fg,
	}
	return d, nil
}

func (d *display) Show() {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	d.s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack))

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
				d.s.Sync()
			}
		}
	}()

	d.Clear()
	d.s.Show()
loop:
	for {
		select {
		case <-quit:
			break loop
		case sp := <-d.sprites:
			d.drawSprite(sp)
		case <-time.After(time.Millisecond * 50):
		}
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
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			d.SetContent(x, y, ' ', nil, d.bgStyle)
		}
	}
}

func (d *display) drawSprite(s sprite) {
	d.SetContent(s.x, s.y, ' ', nil, d.fgStyle)
}

func (d *display) Point(x, y int) {
	var point byte
	point = 1 << 7
	sp := sprite{x, y, []byte{point}}
	d.sprites <- sp
}

func (d *display) Sprite(x, y int, payload []byte) {
	sp := sprite{x, y, payload}
	d.sprites <- sp
}

func (d *display) Clear() {
	d.s.Clear()
	d.DrawScreen(width, height)
}
