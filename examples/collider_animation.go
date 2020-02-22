package main

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

const NUMLEDS = 18
const LEDSFILE = "/dev/leds/colors.json"
const TICKMILLIS = 24

type ColorSet struct {
	Values []string `json:"values"`
}

func main() {
	ticker := time.NewTicker(time.Millisecond * TICKMILLIS)
	inx := 0
	inx2 := NUMLEDS - 1

	for {
		<-ticker.C
		colors := make([]string, NUMLEDS)

		for i := 0; i < NUMLEDS; i++ {
			if i != inx && i != inx2 {
				colors[i] = "#000000"
			}

			if i == inx {
				colors[i] = "#ff0000"
			}

			if i == inx2 {
				colors[i] = "#0000ff"
			}
		}

		if inx == NUMLEDS-1 {
			inx = 0
		} else {
			inx += 1
		}

		if inx2 == 0 {
			inx2 = NUMLEDS - 1
		} else {
			inx2 -= 1
		}

		colorBytes, err := json.Marshal(ColorSet{colors})
		if err != nil {
			panic(err)
		}

		err = ioutil.WriteFile(LEDSFILE, colorBytes, 0644)
		if err != nil {
			panic(err)
		}
	}
}
