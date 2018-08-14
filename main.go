package main

import (
	"compress/algorithm/lz77"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"sort"
	"sync"
	"time"
	//"log"
	//"runtime/pprof"
	"net/http"
	_ "net/http/pprof"
)

var (
    Version string
    Build   string
)

func main() {

    fmt.Println("Version: ", Version)
    fmt.Println("Build time: ", Build)

	go http.ListenAndServe("0.0.0.0:8000", nil)
	cpuNum := runtime.NumCPU()
	runtime.GOMAXPROCS(cpuNum)

	time1 := time.Now().UnixNano()

	//f, err := os.Create("pprof")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//pprof.StartCPUProfile(f)
	//defer pprof.StopCPUProfile()

	fmt.Printf("cpu num : %v..\n", cpuNum)

	decode := flag.Bool("d", false, "true:decode false:encode")
	sourceFile := flag.String("source", "", "source file")
	destFile := flag.String("dest", "", "dest file")
	flag.Parse()

	ch := make(chan *Subsection, cpuNum)

	wg := &sync.WaitGroup{}

	sFile, err := os.Open(*sourceFile)
	defer sFile.Close()
	if err != nil {
		panic(err.Error())
	}

	fileLock := &sync.RWMutex{}

	reqChan := make(chan *TaskInfo, cpuNum)
	chPool := make([]chan *TaskInfo, cpuNum)
	for i := 0; i < cpuNum; i++ {
		chPool[i] = make(chan *TaskInfo)
	}

	var index int64
	if *decode == false {
		fileSize, err := sFile.Seek(0, io.SeekEnd)
		if err != nil {
			panic(err.Error())
		}
		sFile.Seek(0, io.SeekStart)
		index = fileSize / lz77.LZ77_ChunkSize
		if fileSize%lz77.LZ77_ChunkSize != 0 {
			index++
		}
		wg.Add(1)
		go dispatcher(reqChan, wg, cpuNum, fileSize, lz77.LZ77_ChunkSize)
		for i := 0; i < cpuNum; i++ {
			wg.Add(1)
			go compressTask(sFile, wg, ch, chPool[i], reqChan, lz77.LZ77_ChunkSize, fileLock)
		}

	} else {
		wg.Add(1)
			go dispatcherUn(reqChan, wg, cpuNum, sFile, ch)

		for i := 0; i < cpuNum; i++ {
			wg.Add(1)
			go unCompressTask(wg, ch, chPool[i], reqChan)
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
				if index != 0 {
					fmt.Printf("complete %.2f... \n", float64(lastWriteSequeue)/float64(index)*100)
				}

				//fmt.Printf("complete %v... size %v\n",  value.Sequence, len(value.Content))
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
WriteEnd:
	time2 := time.Now().UnixNano()
	ms := (time2 - time1) / 1e6
	fmt.Printf("cost time %vms \n", ms)

	wg.Wait()
	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)
	//fmt.Printf("MemStats Info %+v\n", memStats)
	fmt.Printf("MemStats Alloc %+v\n", memStats.Alloc)
	fmt.Printf("MemStats HeapAlloc %+v\n", memStats.HeapAlloc)
	fmt.Printf("MemStats HeapSys %+v\n", memStats.HeapSys)
	return
}
