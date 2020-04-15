package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Pawka/chip8-emulator/chip8/display"
)

func main() {
	// if len(os.Args) < 2 {
	// 	fmt.Fprintf(os.Stderr, "Usage: %s <path-to-rom>\n", os.Args[0])
	// 	os.Exit(2)
	// }
	// ctx, err := chip8.NewCtxFromArgs(os.Args)
	// if err != nil {
	// 	panic(err)
	// }
	// cpu := chip8.NewChip8()
	// cpu.Run(ctx)
	screen, err := display.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	quit := make(chan struct{})
	go func() {
		screen.Show()
		close(quit)
	}()

	go func() {
		// INVADER
		//
		// X.XXX.X.     0b10111010  $BA
		// .XXXXX..     0b01111100  $7C
		// XX.X.XX.     0b11010110  $D6
		// XXXXXXX.     0b11111110  $FE
		// .X.X.X..     0b01010100  $54
		// X.X.X.X.     0b10101010  $AA
		invader := []byte{0xBA, 0x7C, 0xD6, 0xFE, 0x54, 0xAA}
		for {
			for i := -7; i < 40; i += 3 {
				screen.Sprite(i*2, i, invader)
				time.Sleep(time.Millisecond * 100)
				screen.Clear()
			}
		}
	}()

	<-quit
}
