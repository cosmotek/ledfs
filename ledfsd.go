package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"

	"gitlab.com/rucuriousyet/chateau-gateway/file"
)

type LedFs struct {
	pathfs.FileSystem
	InitTime time.Time
	Files    map[string]nodefs.File
	count    uint64
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

			fileList = append(fileList, fuse.DirEntry{Name: name, Mode: attr.Mode, Ino: attr.Ino})
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
		log.Fatal("Usage:\n  hello MOUNTPOINT")
	}

	fs := &LedFs{
		FileSystem: pathfs.NewDefaultFileSystem(),
		InitTime:   time.Now(),
		Files:      map[string]nodefs.File{},
	}

	fs.Files["colors.json"] = file.NewDataFile([]byte(`{ "values": [] }`), func(data []byte) {
		fmt.Println("render led colors")
	})

	fs.Files["options.json"] = file.NewDataFile([]byte(`{ "numLeds": 24, "gpioPin": 18, "brightness": 220, "dmaChannel": 10 }`), func(data []byte) {
		fmt.Println("reinit leds")
	})

	server, _, err := nodefs.MountRoot(flag.Arg(0), pathfs.NewPathNodeFs(fs, nil).Root(), nil)
	if err != nil {
		log.Fatalf("Mount fail: %v\n", err)
	}

	server.SetDebug(true)
	server.Serve()
}
