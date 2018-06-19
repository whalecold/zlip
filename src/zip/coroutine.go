package main

import (
	"sync"
	"algorithm/lz77"
	"io"
	"os"
	"encoding/binary"
)

func compressCor(sFile *os.File, wg sync.WaitGroup, ch chan *Subsection, offset, index, readSize int64, lock *sync.RWMutex) {

	buffer := make([]byte, readSize)
	chunk := &Subsection{Sequence:index}
	//syscall.Flock(int(sFile.Fd()), syscall.LOCK_SH)		这里的文件锁不知道为什么不能用了
	lock.Lock()
	sFile.Seek(offset, io.SeekStart)
	if _, err := sFile.Read(buffer); err != nil {
		panic(err.Error())
	}
	lock.Unlock()

	chunk.Content = lz77.Lz77Compress(buffer, uint64(readSize))
	lenInfo := make([]byte, 4)
	binary.BigEndian.PutUint32(lenInfo, uint32(len(chunk.Content)))

	chunk.Content = append(lenInfo, chunk.Content...)

	ch <- chunk
	wg.Done()
}
