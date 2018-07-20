package main

import (
	"time"
	"strconv"
	"fmt"

	"github.com/jgarff/rpi_ws281x/golang/ws2811"
)

const RED uint32 = 0x00FFFF // 0x00FF00
const BLU uint32 = 0xFF00FF // 0xGGRRBB
const GRE uint32 = 0xFFFF00

const pin = 18
const count = 24
const brightness = 255

func color(red, green, blue, a byte) uint32 {
	return uint32((a << 24) | (red << 16) | (green << 8) | blue)
}

func colorHex(color string) uint32 {
	// check for len
	if len(color) != 6 && (len(color) == 7 && string(color[0]) != "#") {
		return 0x000000
	}
	
	// check for and remove # symbol
	if string(color[0]) == "#" {
		color = string(color[1:])
	}

	rval := string(color[0:2])
	gval := string(color[2:4])
	bval := string(color[4:6])

	code, err := strconv.ParseUint(fmt.Sprintf("0x%s%s%s", gval, rval, bval), 0, 32)
	if err != nil {
		fmt.Println(err.Error())
		return 0x000000
	}

	return uint32(code)
}

func main() {
	defer ws2811.Fini()
	err := ws2811.Init(pin, count, brightness)
	if err != nil {
		panic(err)
	}

	for {
		for i := 0; i < count; i++ {
			ws2811.SetLed(i, colorHex("#4286f4"))
			if err := ws2811.Render(); err != nil {
				ws2811.Clear()
				panic(err)
			}

			time.Sleep(50 * time.Millisecond)
		}

		for i := 0; i < count; i++ {
			ws2811.SetLed(i, GRE)
			if err := ws2811.Render(); err != nil {
				ws2811.Clear()
				panic(err)
			}

			time.Sleep(50 * time.Millisecond)
		}

		for i := 0; i < count; i++ {
			ws2811.SetLed(i, BLU)
			if err := ws2811.Render(); err != nil {
				ws2811.Clear()
				panic(err)
			}

			time.Sleep(50 * time.Millisecond)
		}
	}
}
