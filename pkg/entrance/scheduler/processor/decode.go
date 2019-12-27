package processor

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/whalecold/zlip/pkg/lz77"
)

func newDecodeProcessor(chunkSize int64, tc chan *TaskProperty, file *os.File) Processor {
	return &decodeProcessor{
		chunkSize: chunkSize,
		taskChan:  tc,
		sFile:     file,
	}
}

type decodeProcessor struct {
	chunkSize int64
	taskChan  chan *TaskProperty
	sFile     *os.File
}

func (de *decodeProcessor) Run(wg *sync.WaitGroup, ch chan *UnitChunk) {
	// buffer to store task data
	cache := make([]byte, de.chunkSize*2)
	for t := range de.taskChan {
		chunk := &UnitChunk{Sequence: t.Index}
		if t.ProcessSize > int64(len(cache)) {
			cache = make([]byte, t.ProcessSize)
		}

		// reset start seek
		_, err := de.sFile.Seek(t.Offset, io.SeekStart)
		if err != nil {
			panic(err)
		}

		// read data from file
		readBuffer := cache[:t.ProcessSize]
		l, err := de.sFile.Read(readBuffer)
		if err != nil {
			panic(err)
		}
		if int64(l) != t.ProcessSize {
			panic(fmt.Errorf("l %v should equal with process size %v", l, t.ProcessSize))
		}
		//outputBuffer := unCompressBuffer[:0]

		//temp  := lz77.UnCompress(t.UnCompress[:t.ReadLen])
		//chunk.Content = *utils.DeepClone(&temp).(*[]byte)
		chunk.Content = lz77.UnCompress(readBuffer)
		//fmt.Printf("t --- %v  index %v 00 %v\n", len(t.UnCompress[:t.ReadLen]), t.Index, len(chunk.Content))
		ch <- chunk
	}
	wg.Done()
}
