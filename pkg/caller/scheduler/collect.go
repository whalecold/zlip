package scheduler

import (
	"os"
	"sort"

	"github.com/whalecold/zlip/pkg/caller/scheduler/processor"
)

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
