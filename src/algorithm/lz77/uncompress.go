package lz77

import (
	"algorithm/huffman"
)

func unCompressSQSub(buffer []byte, h *huffman.HuffmanAlg) []byte {
	lastResult := make([]byte, 0, 1024)

	var byteOffset uint32
	var bitOffset uint32
	for {
		getData, r, b, _ := h.DecodeEle(buffer[byteOffset:], bitOffset)
		byteOffset += r
		bitOffset = b
		if getData == huffman.HUFFMAN_CCLEndFlag {
			break
		}
		lastResult = append(lastResult, byte(getData))
	}
	return lastResult
}

func unCompressSQ(huffmanCode, sq1, sq2 []byte) ([]byte, []byte) {
	ccl := &huffman.HuffmanAlg{}
	ccl.InitCCL()

	ccl.UnSerializeAndBuild(huffmanCode)

	sq1Serial := unCompressSQSub(sq1, ccl)
	sq2Serial := unCompressSQSub(sq2, ccl)

	return sq1Serial, sq2Serial
}
