package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
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
	})

	fs.Files["options.json"] = file.NewDataFile([]byte(`{ "numLeds": 24, "gpioPin": 18, "brightness": 220, "dmaChannel": 10 }`), func(data []byte) {
		options := ColorOptions{}
		err := json.Unmarshal(data, &options)

		if err != nil {
			logger.Error().Err(err).Msg("failed to parse options.json, reverting to default options")
			options = DefaultOptions
		}

		logger.Debug().Interface("options", options).Msg("resetting led options")
	})

	server, _, err := nodefs.MountRoot(flag.Arg(0), pathfs.NewPathNodeFs(fs, nil).Root(), nil)
	if err != nil {
		logger.Fatal().Msg(fmt.Sprintf("Mount fail: %v\n", err))
	}

	server.SetDebug(false)
	server.Serve()
}
