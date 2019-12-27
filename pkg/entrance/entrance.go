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
func Entrance(source, target string, decode bool) {
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
	sc := scheduler.New(source, processor.DecodeType, cpuNum, lz77.ChunkSize)
	count := sc.GetChunkCount()

	collectChan := make(chan *processor.UnitChunk, cpuNum)
	go sc.Run(collectChan)

	//collectData(count)
	recv := make([]*processor.UnitChunk, 0, count)
	dFile, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err.Error())
	}
	defer func() {
		_ = dFile.Close()
	}()

	var lastWriteSequence int64
	for chunk := range collectChan {
		recv = append(recv, chunk)
		sort.Slice(recv, func(i, j int) bool {
			return recv[i].Sequence < recv[j].Sequence
		})

		needRemove := make([]int, 0, len(recv))
		for i, value := range recv {
			if value.Sequence == lastWriteSequence {
				_, err := dFile.Write(value.Content)
				if err != nil {
					panic(err)
				}
				lastWriteSequence++
				needRemove = append(needRemove, i)
				if count != 0 {
					fmt.Printf("complete %.2f... \n", float64(lastWriteSequence)/float64(count)*100)
				}

				//fmt.Printf("complete %v... size %v\n",  value.Sequence, len(value.Content))
				if lastWriteSequence == count {
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

	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)
	//fmt.Printf("MemStats Info %+v\n", memStats)
	fmt.Printf("MemStats Alloc %+v\n", memStats.Alloc)
	fmt.Printf("MemStats HeapAlloc %+v\n", memStats.HeapAlloc)
	fmt.Printf("MemStats HeapSys %+v\n", memStats.HeapSys)
}
