package scheduler

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/whalecold/zlip/pkg/lz77"

	"github.com/whalecold/zlip/pkg/caller/scheduler/processor"
)

type performDispatch func()

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
		temp := make([]byte, lz77.HeadSize)

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
			if l != lz77.HeadSize {
				panic(fmt.Errorf("expected length lz77.HeadSize but get %v", l))
			}
			return false
		}() {
			break
		}
		offset += lz77.HeadSize
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
