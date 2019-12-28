package caller

import (
	"flag"
	"fmt"
	"net/http"
	"runtime"
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

	go sc.Run()

	sc.CollectData(tFile)

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
