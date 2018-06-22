package main

import (
	"flag"
	"os"
	"io"
	"algorithm/lz77"
	"sync"
	"sort"
	"encoding/binary"
	"runtime"
	"utils"
	"fmt"
	"time"
)

func main() {

	cpuNum := runtime.NumCPU()
	runtime.GOMAXPROCS(cpuNum)

	time1 := time.Now().UnixNano()
	//f, err := os.Create("pprof")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//pprof.StartCPUProfile(f)
	//defer pprof.StopCPUProfile()

	decode := flag.Bool("d", false, "true:decode false:encode")
	sourceFile := flag.String("source", "", "source file")
	destFile := flag.String("dest", "", "dest file")
	flag.Parse()

	ch := make(chan *Subsection)

	wg := sync.WaitGroup{}

	sFile, err := os.Open(*sourceFile)
	defer sFile.Close()
	if err != nil {
		panic(err.Error())
	}

	fileLock := &sync.RWMutex{}

	var index int64
	if *decode == false {
		fileSize, err := sFile.Seek(0, io.SeekEnd)
		if err != nil {
			panic(err.Error())
		}
		sFile.Seek(0, io.SeekStart)
		//var offset int64
		//for fileSize > lz77.LZ77_ChunkSize {
		//	wg.Add(1)
		//	go compressCor(sFile, wg, ch, offset, index, lz77.LZ77_ChunkSize, fileLock)
		//	index++
		//	offset+=lz77.LZ77_ChunkSize
		//	fileSize -= lz77.LZ77_ChunkSize
		//}
		//
		//if fileSize != 0 {
		//	wg.Add(1)
		//	go compressCor(sFile, wg, ch, offset, index, fileSize, fileLock)
		//	index++
		//}


		index = fileSize / lz77.LZ77_ChunkSize
		if fileSize % lz77.LZ77_ChunkSize != 0 {
			index++
		}
		reqChan := make(chan *TaskInfo, cpuNum)
		chPool := make([]chan *TaskInfo, cpuNum)
		for i := 0; i < cpuNum; i++ {
			chPool[i] = make(chan *TaskInfo)
		}
		wg.Add(1)
		go dispatcher(reqChan, wg, cpuNum, fileSize, lz77.LZ77_ChunkSize)
		//time.Sleep(time.Second)
		for i := 0; i < cpuNum; i ++ {
			go compressTask(sFile, wg , ch , chPool[i], reqChan, lz77.LZ77_ChunkSize, fileLock)
		}

	} else {

		for {
			temp := make([]byte, 4)
			_, err := sFile.Read(temp)
			if err != nil && err == io.EOF{
				break
			} else if err != nil {
				panic("read file error")
			}
			contentLen := binary.BigEndian.Uint32(temp)
			temp = make([]byte, contentLen)
			_, err = sFile.Read(temp)
			if err != nil && err == io.EOF{
				break
			} else if err != nil {
				panic("read file error")
			}

			newBuffer := utils.DeepClone(&temp).(*[]byte)
			wg.Add(1)
			go func(index int64, newBuffer []byte, wg sync.WaitGroup, ch chan *Subsection) {
				chunk := &Subsection{Sequence:index}
				chunk.Content = lz77.UnLz77Compress(newBuffer)
				ch <- chunk
				wg.Done()
			}(index, *newBuffer, wg, ch)
			index++
		}
	}

	recv := make(SubsectionSlice, 0, index)
	dFile, err := os.OpenFile(*destFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err.Error())
	}
	defer dFile.Close()

	var lastWriteSequeue int64
	for b := range ch {
		recv = append(recv, b)
		sort.Sort(recv)

		needRemove := make([]int, 0, len(recv))
		for i, value := range recv {
			if value.Sequence == lastWriteSequeue {
				dFile.Write(value.Content)
				lastWriteSequeue++
				needRemove = append(needRemove, i)
				fmt.Printf("complete %.2f... size %v\n", float64(lastWriteSequeue)/float64(index) * 100, len(value.Content))
				if lastWriteSequeue == index {
					goto WriteEnd
				}
			} else {
				break
			}
		}

		if len(needRemove) != 0 {
			for i := len(needRemove); i > 0; i-- {
				recv = append(recv[:i-1], recv[i:]...)
			}
		}
	}

	wg.Wait()
	WriteEnd:
		time2 := time.Now().UnixNano()
		ms := (time2 - time1) / 1e6
		fmt.Printf("cost time %vms \n", ms)
	return
}
