package caller

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"sort"
	"time"

	"github.com/whalecold/zlip/pkg/caller/scheduler"
	"github.com/whalecold/zlip/pkg/caller/scheduler/processor"
	"github.com/whalecold/zlip/pkg/lz77"
)

// Run caller
func Run(sFile, tFile string, codeType processor.CodeType) {
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
	sc := scheduler.New(sFile, codeType, cpuNum, lz77.ChunkSize)

	collectChan := make(chan *processor.UnitChunk, cpuNum)

	go sc.Run(collectChan)

	collectData(sc.GetChunkCount(), tFile, collectChan)

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

// collectData collects the data from processors and write them to the target file in order.
func collectData(count int64, tFile string, unChan chan *processor.UnitChunk) {
	// open the file
	dFile, err := os.OpenFile(tFile, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err.Error())
	}
	defer func() {
		_ = dFile.Close()
	}()

	// the next write sequence
	var writeSequence int64
	// chunk slice to store the data from processors
	cs := make([]*processor.UnitChunk, 0, count)
	for chunk := range unChan {
		// receive the data and sort out
		cs = append(cs, chunk)
		sort.Slice(cs, func(i, j int) bool {
			return cs[i].Sequence < cs[j].Sequence
		})

		for i, value := range cs {
			// write data to file in order
			if value.Sequence == writeSequence {
				_, err := dFile.Write(value.Content)
				if err != nil {
					panic(err)
				}
				writeSequence++
			} else {
				// remove the data which is written to file
				cs = cs[i:]
			}
		}
	}
}
