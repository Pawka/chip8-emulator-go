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
		for i := 10; i < 20; i++ {
			screen.Point(i, i)
			time.Sleep(time.Millisecond * 500)
		}
	}()

	<-quit
}
