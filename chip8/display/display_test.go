package display

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	d := New()
	assert.NotNil(t, d.(*display).screen)
}

func TestClear(t *testing.T) {
	d := New()
	dd := d.(*display)
	x := rand.Intn(width)
	y := rand.Intn(height)
	dd.screen[y][x] = 1
	dd.Clear()
	assert.Equal(t, 0, dd.screen[y][x])
}
