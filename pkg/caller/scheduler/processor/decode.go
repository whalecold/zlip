package processor

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/whalecold/zlip/pkg/lz77"
)

func newDecodeProcessor(chunkSize int64, tc chan *TaskProperty, file *os.File, mutex *sync.RWMutex) Processor {
	return &decodeProcessor{
		chunkSize: chunkSize,
		taskChan:  tc,
		sFile:     file,
		mutex:     mutex,
	}
}

type decodeProcessor struct {
	chunkSize int64
	taskChan  chan *TaskProperty
	sFile     *os.File
	mutex     *sync.RWMutex
}

func (de *decodeProcessor) Run(wg *sync.WaitGroup, ch chan *UnitChunk) {
	// buffer to store task data
	cache := make([]byte, de.chunkSize*2)
	for t := range de.taskChan {
		chunk := &UnitChunk{Sequence: t.Index}
		if t.ProcessSize > int64(len(cache)) {
			cache = make([]byte, t.ProcessSize)
		}

		// read data from file
		readBuffer := cache[:t.ProcessSize]
		if func() bool {
			de.mutex.Lock()
			defer de.mutex.Unlock()
			// reset start seek
			_, err := de.sFile.Seek(t.Offset, io.SeekStart)
			if err != nil {
				panic(err)
			}
			l, err := de.sFile.Read(readBuffer)
			if err == io.EOF {
				return true
			}
			if err != nil {
				panic(err)
			}
			if int64(l) != t.ProcessSize {
				panic(fmt.Errorf("l %v should equal with process size %v", l, t.ProcessSize))
			}
			return false
		}() {
			break
		}
		chunk.Content = lz77.UnCompress(readBuffer)
		ch <- chunk
	}
	wg.Done()
}
