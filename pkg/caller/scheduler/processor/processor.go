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

// DataChunk the whole data will be separated into several chunks by fixed size.
type DataChunk struct {
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

func New(typ CodeType, chunkSize int64, tc chan *TaskProperty, file *os.File, mutex *sync.RWMutex) Processor {
	if typ == EncodeType {
		return newEncodeProcessor(chunkSize, tc, file, mutex)
	}
	return newDecodeProcessor(chunkSize, tc, file, mutex)
}
