package main

import (
	"flag"
	"os"
	"io"
	"algorithm/lz77"
	"log"
	"runtime/pprof"
	"sync"
	"sort"
	"encoding/binary"
	"fmt"
	"runtime"
	"utils"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	f, err := os.Create("pprof")
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	decode := flag.Bool("d", false, "true:decode false:encode")
	sourceFile := flag.String("source", "", "source file")
	destFile := flag.String("dest", "", "dest file")
	flag.Parse()

	ch := make(chan *Subsection, 10)
	wg := sync.WaitGroup{}

	sFile, err := os.Open(*sourceFile)
	defer sFile.Close()
	if err != nil {
		panic(err.Error())
	}


	var index int64
	if *decode == false {
		fileSize, err := sFile.Seek(0, io.SeekEnd)
		if err != nil {
			panic(err.Error())
		}
		sFile.Seek(0, io.SeekStart)

		var offset int64
		for fileSize > lz77.LZ77_ChunkSize {
			wg.Add(1)
			go compressCor(sFile, wg, ch, offset, index, lz77.LZ77_ChunkSize)
			index++
			offset+=lz77.LZ77_ChunkSize
			fileSize -= lz77.LZ77_ChunkSize
		}

		if fileSize != 0 {
			wg.Add(1)
			go compressCor(sFile, wg, ch, offset, index, fileSize)
			index++
		}
	} else {

		for {
			temp := make([]byte, 4)
			_, err := sFile.Read(temp)
			if err != nil && err == io.EOF{
				break
			} else if err != nil {
				panic("read file error")
			}
			contentLen := binary.BigEndian.Uint32(temp)
			fmt.Printf("seq %v  length %v\n", index, contentLen)
			temp = make([]byte, contentLen)
			_, err = sFile.Read(temp)
			if err != nil && err == io.EOF{
				break
			} else if err != nil {
				panic("read file error")
			}

			newBuffer := utils.DeepClone(&temp).(*[]byte)
			wg.Add(1)
			go func(index int64, newBuffer []byte, wg sync.WaitGroup, ch chan *Subsection) {
				chunk := &Subsection{Sequence:index}
				chunk.Content = lz77.UnLz77Compress(newBuffer)
				ch <- chunk
				wg.Done()
			}(index, *newBuffer, wg, ch)
			index++
		}
		//newBuffer = huffman.EnCode(buffer)
		//newBuffer = lz77.Lz77Compress(buffer, uint64(fileSize))
	}

	recv := make(SubsectionSlice, 0, index)
	dFile, err := os.OpenFile(*destFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err.Error())
	}
	defer dFile.Close()

	var lastWriteSequeue int64
	for b := range ch {
		recv = append(recv, b)
		sort.Sort(recv)
		for _, value := range recv {
			if value.Sequence == lastWriteSequeue {
				dFile.Write(value.Content)
				lastWriteSequeue++
				if lastWriteSequeue == index {
					goto WriteEnd
				}
			} else {
				break
			}
		}
	}

	WriteEnd:
	return
}
