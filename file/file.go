package file

import (
	"fmt"
	"sync"
	"time"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
)

type dataFile struct {
	data          []byte
	flushCallback func(data []byte)
	crTime        time.Time

	lock sync.Mutex
	nodefs.File
}

func NewDataFile(data []byte, cb func(data []byte)) nodefs.File {
	f := new(dataFile)
	f.data = data
	f.crTime = time.Now()

	f.flushCallback = cb
	f.File = nodefs.NewDefaultFile()

	return f
}

func (f *dataFile) String() string {
	l := len(f.data)
	if l > 10 {
		l = 10
	}

	return fmt.Sprintf("dataFile(%x)", f.data[:l])
}

func (f *dataFile) GetAttr(out *fuse.Attr) fuse.Status {
	out.Mode = fuse.S_IFREG | 0644
	out.Size = uint64(len(f.data))

	out.SetTimes(&f.crTime, &f.crTime, &f.crTime)
	return fuse.OK
}

func (f *dataFile) Read(buf []byte, off int64) (res fuse.ReadResult, code fuse.Status) {
	end := int(off) + int(len(buf))
	if end > len(f.data) {
		end = len(f.data)
	}

	return fuse.ReadResultData(f.data[off:end]), fuse.OK
}

func (f *dataFile) Allocate(off uint64, size uint64, mode uint32) (code fuse.Status) {
	return fuse.OK
}

func (f *dataFile) Write(content []byte, off int64) (uint32, fuse.Status) {
	f.lock.Lock()
	f.data = content

	f.lock.Unlock()
	return uint32(len(content)), fuse.OK
}

func (f *dataFile) Flush() fuse.Status {
	f.flushCallback(f.data)
	return fuse.OK
}

func (f *dataFile) Fsync(flags int) (code fuse.Status) {
	return fuse.OK
}

func (f *dataFile) Truncate(size uint64) (code fuse.Status) {
	return fuse.OK
}
