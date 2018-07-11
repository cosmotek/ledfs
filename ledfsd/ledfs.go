package main

import (
	"time"

	"github.com/jgarff/rpi_ws281x/golang/ws2811"
)

const RED uint32 = 8192
const BLUE uint32 = 32
const GREEN uint32 = 2097152

const pin = 18
const count = 24
const brightness = 255

func main() {
	defer ws2811.Fini()
	err := ws2811.Init(pin, count, brightness)
	if err != nil {
		panic(err)
	}

	for {
		for i := 0; i < count; i++ {
			ws2811.SetLed(i, RED)
			if err := ws2811.Render(); err != nil {
				ws2811.Clear()
				panic(err)
			}

			time.Sleep(50 * time.Millisecond)
		}

		for i := 0; i < count; i++ {
			ws2811.SetLed(i, GREEN)
			if err := ws2811.Render(); err != nil {
				ws2811.Clear()
				panic(err)
			}

			time.Sleep(50 * time.Millisecond)
		}

		for i := 0; i < count; i++ {
			ws2811.SetLed(i, BLUE)
			if err := ws2811.Render(); err != nil {
				ws2811.Clear()
				panic(err)
			}

			time.Sleep(50 * time.Millisecond)
		}
	}
}
