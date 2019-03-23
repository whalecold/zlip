package lz77

import (
	"github.com/whalecold/zlip/pkg/huffman"
)

func unCompressSQSub(buffer []byte, h *huffman.Alg) []byte {
	lastResult := make([]byte, 0, 1024)

	var byteOffset uint32
	var bitOffset uint32
	for {
		getData, r, b, _ := h.DecodeEle(buffer[byteOffset:], bitOffset)
		byteOffset += r
		bitOffset = b
		if getData == huffman.HUFFMANCCLEndFlag {
			break
		}
		lastResult = append(lastResult, byte(getData))
	}
	return lastResult
}

func unCompressSQ(huffmanCode, sq1, sq2 []byte) ([]byte, []byte) {
	ccl := &huffman.Alg{}
	ccl.InitCCL()

	ccl.UnSerializeAndBuild(huffmanCode)

	sq1Serial := unCompressSQSub(sq1, ccl)
	sq2Serial := unCompressSQSub(sq2, ccl)

	return sq1Serial, sq2Serial
}
