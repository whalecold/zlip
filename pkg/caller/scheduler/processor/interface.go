package processor

import "sync"

type Processor interface {
	Run(wg *sync.WaitGroup, ch chan *DataChunk)
}
