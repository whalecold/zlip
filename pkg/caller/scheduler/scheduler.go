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
	// processorNum the count of parallel processor
	processorNum int
	// chunkSize fixed chunk size, the data will separated by the size.
	chunkSize int64
	// retainSize the data size which need to be process, zero means all data were processed.
	retainSize int64
	// fileSize size of file to be processed
	fileSize int64
	//
	chanTask      chan *processor.TaskProperty
	ProcessorPool []processor.Processor

	wg *sync.WaitGroup
	// dispatch the dispatch function
	dispatch performDispatch
	sFile    *os.File
	// mutex file lock, avoid access sFile concurrently
	// TODO replace file lock
	mutex *sync.RWMutex

	collectChan chan *processor.UnitChunk
}

// New new a processor scheduler
func New(source string, typ processor.CodeType, num int, chunkSize int64) *scheduler {
	sFile, err := os.Open(source)
	if err != nil {
		panic(err)
	}

	sc := &scheduler{
		schedulerType: typ,
		processorNum:  num,
		chunkSize:     chunkSize,
		chanTask:      make(chan *processor.TaskProperty, num),
		wg:            &sync.WaitGroup{},
		sFile:         sFile,
		mutex:         &sync.RWMutex{},
		collectChan:   make(chan *processor.UnitChunk, num),
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
	for i := 0; i < num; i++ {
		sc.ProcessorPool = append(sc.ProcessorPool, processor.New(sc.schedulerType, chunkSize, sc.chanTask, sc.sFile, sc.mutex))
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
	for i := range sc.ProcessorPool {
		sc.wg.Add(1)
		go sc.ProcessorPool[i].Run(sc.wg, sc.collectChan)
	}
	sc.wg.Wait()
	close(sc.collectChan)
	_ = sc.sFile.Close()
}
