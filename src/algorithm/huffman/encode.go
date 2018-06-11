package huffman

import (
	"fmt"
	"encoding/binary"
)

func EnCode(bytes []byte) []byte {

	root := buildHuffmanTree(bytes)
	serial := root.serializeTree()
	m := root.transTreeToHuffmanCodeMap()

	treeLen := len(serial)
	//fmt.Printf("treeLen %v\n", treeLen)
	head := make([]byte, 4)
	binary.BigEndian.PutUint32(head, uint32(treeLen))


	huffmanBuffer := make([]byte, 0, len(bytes) + len(head) + treeLen + 4)
	huffmanBuffer = append(huffmanBuffer, head...)

	huffmanBuffer = append(huffmanBuffer, serial...)

	var tempBype byte
	var bitIndex uint8

	bitBuffer := make([]byte, 0, len(bytes))
	var bitLen uint32 //表示压缩之后占用的比特位
	for _, value := range bytes {
		huffmanCode := m[value]
		//fmt.Printf("huffmanCode %b\n", huffmanCode)
		for _, v := range huffmanCode {
			//fmt.Printf("temp data %b\n", v)
			bitLen++
			tempBype = tempBype << 1
			//0异或任何值的时候得到的都为另外一个值 tempByte最后一位肯定为0 所以得到的结果尾部也肯定为v
			tempBype ^= v
			bitIndex++
			if bitIndex == 8 {
				//这个字节写满了 需要重刷新数据
				bitBuffer = append(bitBuffer, tempBype)
				tempBype = 0
				bitIndex = 0
			}
		}
	}

	//表示tempBype中还有数据
	if bitIndex != 0 {
		//fmt.Printf("tail-----%08b index %v\n", tempBype, bitIndex)
		bitIndex = 8 - bitIndex
		tempBype = tempBype << bitIndex
		bitBuffer = append(bitBuffer, tempBype)
	}

	bitLenByte := make([]byte, 4)
	binary.BigEndian.PutUint32(bitLenByte, bitLen)


	huffmanBuffer = append(huffmanBuffer, bitLenByte...)


	huffmanBuffer = append(huffmanBuffer, bitBuffer...)
	//fmt.Printf("++++++bitLen %v  stream %08b\n", bitLen, bitBuffer)
	return huffmanBuffer
}
