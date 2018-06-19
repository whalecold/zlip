package main

import (
	"sync"
	"algorithm/lz77"
	"io"
	"os"
	"encoding/binary"
	"fmt"
)

func compressCor(sFile *os.File, wg sync.WaitGroup, ch chan *Subsection, offset, index, readSize int64) {
	buffer := make([]byte, readSize)
	chunk := &Subsection{Sequence:index}
	sFile.Seek(offset, io.SeekStart)
	if _, err := sFile.Read(buffer); err != nil {
		panic(err.Error())
	}

	chunk.Content = lz77.Lz77Compress(buffer, uint64(readSize))
	lenInfo := make([]byte, 4)
	binary.BigEndian.PutUint32(lenInfo, uint32(len(chunk.Content)))
	fmt.Printf("seq %v  length %v\n", index, len(chunk.Content))

	chunk.Content = append(lenInfo, chunk.Content...)

	ch <- chunk
	wg.Done()
}
