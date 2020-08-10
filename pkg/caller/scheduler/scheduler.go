package scheduler

import (
	"os"
	"sync"

	"github.com/whalecold/zlip/pkg/caller/scheduler/processor"
)

// scheduler simple implement of performer processor parallel.
type scheduler struct {
	// schedulerType the task type
	schedulerType processor.CodeType
	// chunkSize fixed chunk size, the data will separated by the size.
	chunkSize int64
	// retainSize the data size which need to be process, zero means all data were processed.
	retainSize int64
	// fileSize size of file to be processed
	fileSize int64
	chanTask chan *processor.TaskProperty
	pools    []processor.Processor

	wg *sync.WaitGroup
	// dispatch the dispatch function
	dispatch dispatchFn
	sFile    *os.File
	// mutex file lock, avoid access sFile concurrently
	// TODO replace file lock
	mutex *sync.RWMutex

	collectChan chan *processor.DataChunk
}

// New new a processor scheduler
func New(source string, typ processor.CodeType, taskNum int, chunkSize int64) *scheduler {
	sFile, err := os.Open(source)
	if err != nil {
		panic(err)
	}

	sc := &scheduler{
		schedulerType: typ,
		chunkSize:     chunkSize,
		chanTask:      make(chan *processor.TaskProperty, taskNum),
		wg:            &sync.WaitGroup{},
		sFile:         sFile,
		mutex:         &sync.RWMutex{},
		collectChan:   make(chan *processor.DataChunk, taskNum),
	}

	// set basic info
	if typ == processor.EncodeType {
		sc.dispatch = sc.encodeDispatch
		sInfo, err := sc.sFile.Stat()
		if err != nil {
			panic(err)
		}
		sc.retainSize = sInfo.Size()
		sc.fileSize = sc.retainSize
	} else if typ == processor.DecodeType {
		sc.dispatch = sc.decodeDispatch
	} else {
		panic("error typ")
	}

	// init processor
	for i := 0; i < taskNum; i++ {
		sc.pools = append(sc.pools, processor.New(sc.schedulerType, chunkSize, sc.chanTask, sc.sFile, sc.mutex))
	}
	return sc
}

// getChunkCount returns the chunk count
func (sc *scheduler) getChunkCount() int64 {
	count := sc.fileSize / sc.chunkSize
	if sc.fileSize%sc.chunkSize != 0 {
		count++
	}
	return count
}

func (sc *scheduler) Run() {
	go sc.dispatch()
	for i := range sc.pools {
		sc.wg.Add(1)
		go sc.pools[i].Run(sc.wg, sc.collectChan)
	}
	sc.wg.Wait()
	close(sc.collectChan)
	_ = sc.sFile.Close()
}
