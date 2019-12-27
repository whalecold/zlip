package processor

import (
	"os"
	"sync"
)

type CodeType string

const (
	EncodeType CodeType = "encode"
	DecodeType CodeType = "decode"
)

// UnitChunk the whole data will be separated into several chunks by fixed size, the scheduler dispatches
// those chunks to processors, the processors performer task parallel. After processing the data, scheduler will
// collects and sorts out the result. The out data may be out of order as the parallel performing, so we need mark sequence
// to every chunk for ordering result.
type UnitChunk struct {
	// Sequence sequence of chunk
	Sequence int64
	Content  []byte
}

// TaskProperty sub task property
type TaskProperty struct {
	// Offset the offset of the file to starts reading.
	Offset int64
	// ProcessSize the data size of this task will process.
	ProcessSize int64
	// Index the task index
	Index int64
}

type Processor interface {
	Run(wg *sync.WaitGroup, ch chan *UnitChunk)
}

func New(typ CodeType, chunkSize int64, tc chan *TaskProperty, file *os.File) Processor {
	if typ == EncodeType {
		return newEncodeProcessor(chunkSize, tc, file)
	}
	return newDecodeProcessor(chunkSize, tc, file)
}
