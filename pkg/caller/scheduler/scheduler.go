package scheduler

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sort"
	"sync"

	"github.com/whalecold/zlip/pkg/caller/scheduler/processor"
)

type performDispatch func()

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
	mutex    *sync.RWMutex

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

func (sc *scheduler) encodeDispatch() {
	// the data size has been processed
	var processedSize int64
	for index := int64(0); sc.retainSize > 0; index++ {
		tp := &processor.TaskProperty{
			Index:  index,
			Offset: processedSize,
		}

		if sc.retainSize >= sc.chunkSize {
			tp.ProcessSize = sc.chunkSize
			sc.retainSize -= sc.chunkSize
		} else {
			tp.ProcessSize = sc.retainSize
			sc.retainSize = 0
		}
		// only send the task when process size is not empty
		if tp.ProcessSize != 0 {
			sc.chanTask <- tp
		}

		processedSize += tp.ProcessSize
	}
	close(sc.chanTask)
}

func (sc *scheduler) decodeDispatch() {

	var offset int64
	for index := int64(0); ; index++ {
		temp := make([]byte, 4)

		if func() bool {
			sc.mutex.Lock()
			defer sc.mutex.Unlock()
			if _, err := sc.sFile.Seek(offset, io.SeekStart); err != nil {
				panic(err)
			}
			// read head info to the temp
			l, err := sc.sFile.Read(temp)
			if err == io.EOF {
				return true
			}
			if err != nil {
				panic(err)
			}
			if l != 4 {
				panic(fmt.Errorf("expected length 4 but get %v", l))
			}
			return false
		}() {
			break
		}
		offset += 4
		tp := &processor.TaskProperty{
			Index:       index,
			Offset:      offset,
			ProcessSize: int64(binary.BigEndian.Uint32(temp)),
		}
		sc.chanTask <- tp
		offset += tp.ProcessSize
	}
	close(sc.chanTask)
}

// CollectData collects the data from processors and write them to the target file in order.
func (sc *scheduler) CollectData(tFile string) {
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
	cs := make([]*processor.UnitChunk, 0, sc.getChunkCount())
	for chunk := range sc.collectChan {
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
