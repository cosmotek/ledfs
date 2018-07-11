package neopixel

/*
#cgo CFLAGS: -std=c99
#cgo LDFLAGS: -lws2811
#include "ws2811.go.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

// Init initializes the LED set for communication
func Init(gpioPin, ledCount, brightness int) error {
	C.ledstring.channel[0].gpionum = C.int(gpioPin)
	C.ledstring.channel[0].count = C.int(ledCount)
	C.ledstring.channel[0].brightness = C.uint8_t(brightness)

	if res := int(C.ws2811_init(&C.ledstring)); res != 0 {
		return errors.New(fmt.Sprintf("failed to initialize neopixels: %v", res))
	}

	return nil
}

// Close gracefully closes the LED communication when everything
// is completed and cleans up the connection
func Close() {
	C.ws2811_fini(&C.ledstring)
}

// Show renders the current color buffer
func Show() error {
	if res := int(C.ws2811_render(&C.ledstring)); res != 0 {
		return errors.New(fmt.Sprintf("failed to show neopixels: %v", res))
	}

	return nil
}

func Wait() error {
	if res := int(C.ws2811_wait(&C.ledstring)); res == 0 {
		return errors.New(fmt.Sprintf("failed to wait neopixels", res))
	}

	return nil
}

// SetPixelColor writes an index and color value to the
// buffer, to be drawn when Show() is called next.
func SetPixelColor(index int, value uint32) {
	C.ws2811_set_led(&C.ledstring, C.int(index), C.uint32_t(value))
}

// Clear clears the color buffer. Make sure to follow
// this call up with Show()
func Clear() {
	C.ws2811_clear(&C.ledstring)
}

func SetBitmap(a []uint32) {
	C.ws2811_set_bitmap(&C.ledstring, unsafe.Pointer(&a[0]), C.int(len(a)*4))
}
