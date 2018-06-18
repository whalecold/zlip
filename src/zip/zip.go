package main

import (
	"flag"
	"os"
	"io"
	"algorithm/lz77"
	"log"
	"runtime/pprof"
)

func main() {
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

	sFile, err := os.Open(*sourceFile)
	defer sFile.Close()
	if err != nil {
		panic(err.Error())
	}

	fileSize, err := sFile.Seek(0, io.SeekEnd)
	if err != nil {
		panic(err.Error())
	}

	sFile.Seek(0, io.SeekStart)

	buffer := make([]byte, fileSize)
	if _, err = sFile.Read(buffer); err != nil {
		panic(err.Error())
	}

	var newBuffer []byte
	if *decode == true {
		//newBuffer = huffman.Decode(buffer)
		newBuffer = lz77.UnLz77Compress(buffer)
	} else {
		//newBuffer = huffman.EnCode(buffer)
		newBuffer = lz77.Lz77Compress(buffer, uint64(fileSize))
	}

	dFile, err := os.OpenFile(*destFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err.Error())
	}
	dFile.Write(newBuffer)
	defer dFile.Close()
	return
}
