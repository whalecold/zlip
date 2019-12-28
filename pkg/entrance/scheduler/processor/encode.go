package processor

import (
	"encoding/binary"
	"io"
	"os"
	"sync"

	"github.com/whalecold/zlip/pkg/lz77"
)

func newEncodeProcessor(chunkSize int64, tc chan *TaskProperty, file *os.File) Processor {
	return &encodeProcessor{
		chanTask:  tc,
		chunkSize: chunkSize,
		sFile:     file,
	}
}

type encodeProcessor struct {
	chanTask  chan *TaskProperty
	chunkSize int64
	sFile     *os.File
}

func (en *encodeProcessor) Run(wg *sync.WaitGroup, ch chan *UnitChunk) {
	// buffer to store task data
	cache := make([]byte, en.chunkSize)
	// stores the compressing result
	encodeBuffer := make([]byte, en.chunkSize)
	for t := range en.chanTask {

		chunk := &UnitChunk{Sequence: t.Index}
		if t.ProcessSize > en.chunkSize {
			panic("error t.ProcessSize")
		}

		// reset start seek
		_, err := en.sFile.Seek(t.Offset, io.SeekStart)
		if err != nil {
			panic(err)
		}

		// read data from file
		readBuffer := cache[:t.ProcessSize]
		if l, err := en.sFile.Read(readBuffer); err != nil || int64(l) != t.ProcessSize {
			panic(err)
		}

		// perform compress
		result := encodeBuffer[:0]
		chunk.Content, _ = lz77.Compress(readBuffer, &result, uint64(t.ProcessSize))

		// add metadata info to the head
		lenInfo := make([]byte, 4)
		binary.BigEndian.PutUint32(lenInfo, uint32(len(chunk.Content)))
		chunk.Content = append(lenInfo, chunk.Content...)
		ch <- chunk
	}
	wg.Done()
}
