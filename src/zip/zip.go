package main

import (
	"flag"
	"os"
	"io"
	"algorithm/huffman"
)

func main() {

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
		newBuffer = huffman.Decode(buffer)
	} else {
		newBuffer = huffman.EnCode(buffer)
	}

	dFile, err := os.OpenFile(*destFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err.Error())
	}
	dFile.Write(newBuffer)
	defer dFile.Close()
	return
}
