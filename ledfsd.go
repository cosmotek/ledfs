package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"github.com/jgarff/rpi_ws281x/golang/ws2811"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gitlab.com/rucuriousyet/chateau-gateway/file"
)

type ColorSet struct {
	Values []string `json:"values"`
}

type ColorOptions struct {
	NumLEDs    uint32 `json:"numLeds"`
	GPIOPin    byte   `json:"gpioPin"`
	Brightness byte   `json:"brightness"`
	DMAChannel byte   `json:"dmaChannel"`
}

type LedFs struct {
	pathfs.FileSystem
	InitTime time.Time
	Files    map[string]nodefs.File
	count    uint64
}

var logger zerolog.Logger

var DefaultOptions = ColorOptions{
	NumLEDs:    24,
	GPIOPin:    18,
	Brightness: 220,
	DMAChannel: 10,
}

var DefaultColors = ColorSet{
	Values: []string{},
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

func (fs *LedFs) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	accessTime := time.Now()

	if name == "" {
		attr := fuse.Attr{
			Mode:  fuse.S_IFDIR | 0755,
			Owner: context.Owner,
		}

		attr.SetTimes(&accessTime, &fs.InitTime, &fs.InitTime)
		return &attr, fuse.OK
	}

	file, ok := fs.Files[name]
	if ok {
		attr := fuse.Attr{}
		stat := file.GetAttr(&attr)

		return &attr, stat
	}

	return nil, fuse.ENOENT
}

func (fs *LedFs) OpenDir(name string, context *fuse.Context) (c []fuse.DirEntry, code fuse.Status) {
	if name == "" {
		fileList := []fuse.DirEntry{}

		for name, file := range fs.Files {
			attr := fuse.Attr{}
			file.GetAttr(&attr)

			fileList = append(fileList, fuse.DirEntry{Name: name, Mode: attr.Mode})
		}

		return fileList, fuse.OK
	}

	return nil, fuse.ENOENT
}

func (fs *LedFs) Open(name string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	ffile, ok := fs.Files[name]
	if ok {
		return ffile, fuse.OK
	}

	return nil, fuse.ENOENT
}

func (fs *LedFs) Create(name string, flags uint32, mode uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	ffile, ok := fs.Files[name]
	if ok {
		return ffile, fuse.OK
	}

	return nil, fuse.EROFS
}

func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		logger.Fatal().Msg("Usage:\n  hello MOUNTPOINT")
	}

	logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	defer ws2811.Fini()

	err := ws2811.Init(int(DefaultOptions.GPIOPin), int(DefaultOptions.NumLEDs), int(DefaultOptions.Brightness))
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to init leds")
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for {
			<-c

			// sig is a ^C, handle it
			for i := uint32(0); i < DefaultOptions.NumLEDs; i++ {
				ws2811.SetLed(int(i), 0x000000)
			}

			if err := ws2811.Render(); err != nil {
				ws2811.Clear()
				logger.Error().Err(err).Msg("failed to render led colors")
			}

			ws2811.Fini()
			os.Exit(0)
		}
	}()

	for i := uint32(0); i < DefaultOptions.NumLEDs; i++ {
		ws2811.SetLed(int(i), 0x000000)
	}

	if err := ws2811.Render(); err != nil {
		ws2811.Clear()
		logger.Error().Err(err).Msg("failed to render led colors")
	}

	fs := &LedFs{
		FileSystem: pathfs.NewDefaultFileSystem(),
		InitTime:   time.Now(),
		Files:      map[string]nodefs.File{},
	}

	fs.Files["colors.json"] = file.NewDataFile([]byte(`{ "values": [] }`), func(data []byte) {
		colors := ColorSet{}
		err := json.Unmarshal(data, &colors)

		if err != nil {
			logger.Error().Err(err).Msg("failed to parse colors.json, reverting to default colors")
			colors = DefaultColors
		}

		logger.Debug().Interface("colors", colors).Msg("rendering led colors")
		for i, color := range colors.Values {
			ws2811.SetLed(i, colorHex(color))
		}

		if err := ws2811.Render(); err != nil {
			ws2811.Clear()
			logger.Error().Err(err).Msg("failed to render led colors")
		}
	})

	fs.Files["options.json"] = file.NewDataFile([]byte(`{ "numLeds": 24, "gpioPin": 18, "brightness": 220, "dmaChannel": 10 }`), func(data []byte) {
		options := ColorOptions{}
		err := json.Unmarshal(data, &options)

		if err != nil {
			logger.Error().Err(err).Msg("failed to parse options.json, reverting to default options")
			options = DefaultOptions
		}

		logger.Debug().Interface("options", options).Msg("resetting led options")
		ws2811.Fini()

		err = ws2811.Init(int(options.GPIOPin), int(options.NumLEDs), int(options.Brightness))
		if err != nil {
			logger.Fatal().Err(err).Msg("failed to init leds")
		}

		for i := uint32(0); i < DefaultOptions.NumLEDs; i++ {
			ws2811.SetLed(int(i), 0x000000)
		}

		if err := ws2811.Render(); err != nil {
			ws2811.Clear()
			logger.Error().Err(err).Msg("failed to render led colors")
		}
	})

	server, _, err := nodefs.MountRoot(flag.Arg(0), pathfs.NewPathNodeFs(fs, nil).Root(), nil)
	if err != nil {
		logger.Fatal().Msg(fmt.Sprintf("Mount fail: %v\n", err))
	}

	server.SetDebug(false)
	server.Serve()
}
