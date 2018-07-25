package main

import (
	"algorithm/lz77"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sync"
)

type TaskInfo struct {
	Offset     int64
	Index      int64
	ReadSize   int64
	ReqCh      chan *TaskInfo
	UnCompress []byte
	ReadLen    int64
}

func compressCor(sFile *os.File, wg sync.WaitGroup, ch chan<- *Subsection, offset, index, readSize int64, lock *sync.RWMutex) {

	buffer := make([]byte, readSize)
	chunk := &Subsection{Sequence: index}
	//syscall.Flock(int(sFile.Fd()), syscall.LOCK_SH)		这里的文件锁不知道为什么不能用了
	lock.Lock()
	sFile.Seek(offset, io.SeekStart)
	fmt.Printf("read len %v   \n", len(buffer))
	if _, err := sFile.Read(buffer); err != nil {
		panic(err.Error())
	}
	lock.Unlock()

	outBuffer := make([]byte, 0, 1024*1024)
	chunk.Content, _ = lz77.Lz77Compress(buffer, &outBuffer, uint64(readSize))
	lenInfo := make([]byte, 4)
	binary.BigEndian.PutUint32(lenInfo, uint32(len(chunk.Content)))

	chunk.Content = append(lenInfo, chunk.Content...)

	ch <- chunk
	wg.Done()
}

//任务池中的一个子任务
func compressTask(sFile *os.File, wg *sync.WaitGroup, ch chan *Subsection, inCh chan *TaskInfo, reqCh chan *TaskInfo, readSize int64, lock *sync.RWMutex) {
	buffer := make([]byte, readSize)
	compressBuffer := make([]byte, readSize)
	reqCh <- &TaskInfo{ReqCh: inCh}
	for recv := range inCh {
		chunk := &Subsection{Sequence: recv.Index}
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
		if _, err := sFile.Read(readBuffer); err != nil {
			panic(err.Error())
		}
		lock.Unlock()

		outBuffer := compressBuffer[:0]
		chunk.Content, _ = lz77.Lz77Compress(readBuffer, &outBuffer, uint64(recv.ReadSize))
		lenInfo := make([]byte, 4)
		binary.BigEndian.PutUint32(lenInfo, uint32(len(chunk.Content)))
		chunk.Content = append(lenInfo, chunk.Content...)

		ch <- chunk
		reqCh <- &TaskInfo{ReqCh: inCh}
	}
	wg.Done()
}

//想了一下好像也没有必要统一调度 各自协程解锁也可以解决这个问题啊
func dispatcher(dispChan chan *TaskInfo, wg *sync.WaitGroup, cupNum int, fileSize, chunkSize int64) {
	var endNum int
	var index int64
	var offset int64
	for req := range dispChan {
		taskInfo := &TaskInfo{Index: index, Offset: offset}
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

func dispatcherUn(dispChan chan *TaskInfo, wg *sync.WaitGroup, cupNum int, sFile *os.File, needClose chan *Subsection) {

	var endNum int
	fileEnd := false
	var index int64

	for req := range dispChan {
		if fileEnd == true {
			endNum++
			close(req.ReqCh)
		} else {
			temp := make([]byte, 4)
			_, err := sFile.Read(temp)
			if err != nil && err == io.EOF {
				endNum++
				fileEnd = true
				close(req.ReqCh)
			} else if err != nil {
				panic("read file error")
			} else {
				contentLen := binary.BigEndian.Uint32(temp)
				//fmt.Printf("read len %v  facet %v index %v\n", contentLen, len(req.UnCompress), index)
				readBuffer := req.UnCompress[:contentLen]
				_, err = sFile.Read(readBuffer)
				if err != nil && err == io.EOF {
					endNum++
					fileEnd = true
					close(req.ReqCh)
				} else if err != nil {
					panic("read file error")
				} else {
					req.ReqCh <- &TaskInfo{UnCompress: req.UnCompress, Index: index, ReadLen: int64(contentLen)}
					index++
				}
			}
		}

		if endNum >= cupNum {
			close(needClose)
			break
		}
	}

	wg.Done()
}

func unCompressTask(wg *sync.WaitGroup, ch chan *Subsection, inCh chan *TaskInfo, reqCh chan *TaskInfo) {
	//unCompressBuffer := make([]byte, lz77.LZ77_ChunkSize * 2)
	readBuffer := make([]byte, lz77.LZ77_ChunkSize*2)
	reqCh <- &TaskInfo{ReqCh: inCh, UnCompress: readBuffer}
	for recv := range inCh {
		chunk := &Subsection{Sequence: recv.Index}
		//outputBuffer := unCompressBuffer[:0]

		//temp  := lz77.UnLz77Compress(recv.UnCompress[:recv.ReadLen])
		//chunk.Content = *utils.DeepClone(&temp).(*[]byte)
		chunk.Content = lz77.UnLz77Compress(recv.UnCompress[:recv.ReadLen])
		//fmt.Printf("recv --- %v  index %v 00 %v\n", len(recv.UnCompress[:recv.ReadLen]), recv.Index, len(chunk.Content))
		ch <- chunk
		reqCh <- &TaskInfo{ReqCh: inCh, UnCompress: readBuffer}
	}
	wg.Done()
}
