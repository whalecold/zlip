package entrance

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"sort"
	"time"

	"github.com/whalecold/zlip/pkg/entrance/scheduler"
	"github.com/whalecold/zlip/pkg/entrance/scheduler/processor"
	"github.com/whalecold/zlip/pkg/lz77"
)

// Entrance entrance
func Entrance(source, target string, codeType processor.CodeType) {
	go func() {
		if err := http.ListenAndServe("0.0.0.0:8000", nil); err != nil {
			panic(err)
		}
	}()
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
	flag.Parse()

	// perform scheduler
	sc := scheduler.New(source, codeType, cpuNum, lz77.ChunkSize)
	count := sc.GetChunkCount()

	collectChan := make(chan *processor.UnitChunk, cpuNum)
	go sc.Run(collectChan)

	collectData(count, target, collectChan)

	time2 := time.Now().UnixNano()
	ms := (time2 - time1) / 1e6
	fmt.Printf("cost time %vms \n", ms)

	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)
	//fmt.Printf("MemStats Info %+v\n", memStats)
	fmt.Printf("MemStats Alloc %+v\n", memStats.Alloc)
	fmt.Printf("MemStats HeapAlloc %+v\n", memStats.HeapAlloc)
	fmt.Printf("MemStats HeapSys %+v\n", memStats.HeapSys)
}

func collectData(count int64, target string, cc chan *processor.UnitChunk) {
	// open the file
	dFile, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err.Error())
	}
	defer func() {
		_ = dFile.Close()
	}()

	// the next write sequence
	var writeSequence int64
	chunkSlice := make([]*processor.UnitChunk, 0, count)
	for chunk := range cc {
		// receive the data and sort out
		chunkSlice = append(chunkSlice, chunk)
		sort.Slice(chunkSlice, func(i, j int) bool {
			return chunkSlice[i].Sequence < chunkSlice[j].Sequence
		})

		for i, value := range chunkSlice {
			if value.Sequence == writeSequence {
				_, err := dFile.Write(value.Content)
				if err != nil {
					panic(err)
				}
				writeSequence++
			} else {
				// remove the data which is written to file
				chunkSlice = chunkSlice[i:]
			}
		}
	}
}
