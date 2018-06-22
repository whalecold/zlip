package main

import (
	"sync"
	"algorithm/lz77"
	"io"
	"os"
	"encoding/binary"
	"fmt"
)

type TaskInfo struct {
	Offset  int64
	Index 	int64
	ReadSize int64
	ReqCh 	chan *TaskInfo
}


func compressCor(sFile *os.File, wg sync.WaitGroup, ch chan<- *Subsection, offset, index, readSize int64, lock *sync.RWMutex) {

	buffer := make([]byte, readSize)
	chunk := &Subsection{Sequence:index}
	//syscall.Flock(int(sFile.Fd()), syscall.LOCK_SH)		这里的文件锁不知道为什么不能用了
	lock.Lock()
	sFile.Seek(offset, io.SeekStart)
	fmt.Printf("read len %v   \n", len(buffer))
	if _, err := sFile.Read(buffer); err != nil {
		panic(err.Error())
	}
	lock.Unlock()

	outBuffer := make([]byte, 0, 1024 * 1024)
	chunk.Content = lz77.Lz77Compress(buffer, outBuffer, uint64(readSize))
	lenInfo := make([]byte, 4)
	binary.BigEndian.PutUint32(lenInfo, uint32(len(chunk.Content)))

	chunk.Content = append(lenInfo, chunk.Content...)

	ch <- chunk
	wg.Done()
}

//任务池中的一个子任务
func compressTask(sFile *os.File, wg sync.WaitGroup, ch chan *Subsection, inCh chan *TaskInfo, reqCh chan  *TaskInfo, readSize int64, lock *sync.RWMutex) {
	buffer := make([]byte, readSize)
	compressBuffer := make([]byte, readSize)
	reqCh <- &TaskInfo{ReqCh:inCh}
	for recv := range inCh {
		chunk := &Subsection{Sequence:recv.Index}
		var readBuffer []byte
		if recv.ReadSize == readSize {
			readBuffer = buffer
		} else if recv.ReadSize < readSize {
			readBuffer = buffer[:recv.ReadSize]
		} else {
			panic("error recv.ReadSize")
		}

		lock.Lock()
		sFile.Seek(recv.Offset, io.SeekStart)
		fmt.Printf("read len %v   \n", len(readBuffer))
		if _, err := sFile.Read(readBuffer); err != nil {
			panic(err.Error())
		}
		lock.Unlock()

		outBuffer := compressBuffer[:0]
		chunk.Content = lz77.Lz77Compress(readBuffer, outBuffer, uint64(recv.ReadSize))
		lenInfo := make([]byte, 4)
		binary.BigEndian.PutUint32(lenInfo, uint32(len(chunk.Content)))
		fmt.Printf("compressTask len %v ---\n", len(chunk.Content))
		chunk.Content = append(lenInfo, chunk.Content...)

		ch <- chunk
		reqCh <- &TaskInfo{ReqCh:inCh}
	}
	wg.Done()
}

//想了一下好像也没有必要统一调度 各自协程解锁也可以解决这个问题啊- -！ 尴尬
func dispatcher(dispChan chan *TaskInfo, wg sync.WaitGroup, cupNum int, fileSize, chunkSize int64) {
	var endNum int
	var index int64
	var offset int64
	for req := range dispChan {
		taskInfo := &TaskInfo{Index:index, Offset:offset}
		if fileSize == 0 {
			close(req.ReqCh)
			endNum++
		} else if fileSize > chunkSize {
			taskInfo.ReadSize = chunkSize
			fileSize -= chunkSize
		} else {
			taskInfo.ReadSize = fileSize
			fileSize = 0
		}

		//fmt.Printf("ReadSize %v index %v offset %v\n", taskInfo.ReadSize, taskInfo.Index, taskInfo.Offset)
		if taskInfo.ReadSize != 0 {
			req.ReqCh <- taskInfo
		}

		if endNum >= cupNum {
			break
		}
		index++
		offset += chunkSize
	}
	wg.Done()
}
