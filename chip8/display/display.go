package display

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell"
)

const (
	width          = 64
	height         = 32
	debuggerHeight = 10
)

// Display defines interface of CHIP8 display.
type Display interface {
	// Clear the screen.
	Clear()
	// Show the screen.
	Show()
	// Point draws a point at x y. Return true if any pixel is flipped from high
	// to low.
	Point(x, y int) bool
	// Sprite at x, y coordinates draws a sprite which is 8 symbols width and n
	// lines height. Return true if any pixel is flipped from high to low.
	Sprite(x, y int, payload []byte) bool

	// PollKey returns a pressed key.
	PollKey() *rune

	// Debug prints information at the top left corner of sceen.
	Debug(line string)
}

type display struct {
	debugLines       []string
	s                tcell.Screen
	keych            chan rune
	quit             chan struct{}
	screen           [][]int
	sprites          chan sprite
	bgStyle, fgStyle tcell.Style
}

type sprite struct {
	x, y        int
	payload     []byte
	collisionch chan bool
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
		debugLines: make([]string, debuggerHeight),
		s:          s,
		sprites:    make(chan sprite),
		keych:      make(chan rune, 10),
		bgStyle:    bg,
		fgStyle:    fg,
	}
	return d, nil
}

func (d *display) Show() {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	d.s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack))

	d.quit = make(chan struct{})
	go func() {
		for {
			ev := d.s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyCtrlC ||
					ev.Key() == tcell.KeyEscape {
					close(d.quit)
					return
				}
				keys := []rune{
					'1', '2', '3', '4',
					'q', 'w', 'e', 'r',
					'a', 's', 'd', 'f',
					'z', 'x', 'c', 'v',
				}
				for _, k := range keys {
					if ev.Rune() == k {
						go func() {
							d.keych <- ev.Rune()
						}()
					}
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
		case <-d.quit:
			break loop
		case sp := <-d.sprites:
			d.drawSprite(sp)
		case <-time.After(time.Millisecond * 50):
		}
		d.s.Show()

	}
	d.s.Fini()
}

func (d *display) setContent(x int, y int, mainc rune, combc []rune, style tcell.Style) {
	dw, dh := d.s.Size()
	_y := dh/2 - height/2
	_x := dw/2 - width/2
	d.s.SetContent(_x+x, _y+y, mainc, combc, style)
}

func (d *display) isSetContent(x int, y int) bool {
	dw, dh := d.s.Size()
	_y := dh/2 - height/2
	_x := dw/2 - width/2
	_, _, style, _ := d.s.GetContent(_x+x, _y+y)
	return style == d.getStyle(1)
}

func (d *display) DrawScreen(w, h int) {
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			d.setContent(x, y, ' ', nil, d.bgStyle)
		}
	}
}

func (d *display) getStyle(b byte) tcell.Style {
	st := map[byte]tcell.Style{
		0: d.bgStyle,
		1: d.fgStyle,
	}

	return st[b]
}

func (d *display) drawSprite(s sprite) {
	const spriteWidth = 8
	collision := false

	for row := 0; row < len(s.payload); row++ {
		for i := spriteWidth - 1; i >= 0; i-- {
			pixel := (s.payload[row] >> i) & 0x1
			x := s.x + spriteWidth - i
			y := s.y + row
			if x < 0 || x >= width || y < 0 || y >= height {
				continue
			}
			if collision == false {
				set := d.isSetContent(x, y)
				if set == true && pixel == 1 {
					collision = true
				}
			}
			d.setContent(x, y, ' ', nil, d.getStyle(pixel))
		}
	}
	s.collisionch <- collision
}

func (d *display) Point(x, y int) bool {
	var point byte
	point = 1 << 7
	ch := make(chan bool)
	sp := sprite{x, y, []byte{point}, ch}
	d.sprites <- sp
	return <-sp.collisionch
}

func (d *display) Sprite(x, y int, payload []byte) bool {
	ch := make(chan bool)
	sp := sprite{x - 1, y, payload, ch}
	d.sprites <- sp
	return <-sp.collisionch
}

func (d *display) Clear() {
	d.s.Clear()
	d.DrawScreen(width, height)
}

func (d *display) PollKey() *rune {
	select {
	case <-d.quit:
		return nil
	case r := <-d.keych:
		return &r
	default:
		return nil
	}
}

func (d *display) Debug(line string) {
	copy(d.debugLines[1:], d.debugLines[:debuggerHeight-1])
	d.debugLines[0] = line

	for i, l := range d.debugLines {
		d.s.SetContent(0, i, ' ', []rune(l), d.fgStyle)
	}
}
