package huffman

import (
	"encoding/binary"
	"fmt"
)

func decodeTest(bytes []byte) []byte {
	root := buildHuffmanTree(bytes)
	//root1 := buildHuffmanTree(bytes)
	//root2 := buildHuffmanTree(bytes)
	m := root.transTreeToHuffmanCodeMap()
	for k, v := range m {
		fmt.Printf("key : %c  value : %b\n", k, v)
	}
	fmt.Printf("-------------next-----------------\n")
	headTest := buildHuffmanTree(bytes)

	pre := headTest.genStreamByPreorder()
	in := headTest.genStreamByInorder()

	for _, value := range in {
		fmt.Printf("%v\t", value)
	}
	fmt.Printf("\n")

	result := buildTreeBySlice(pre, in)
	pre111 := result.genStreamByInorder()
	for _, value := range pre111 {
		fmt.Printf("%v\t", value)
	}
	fmt.Printf("\n")

	mm := result.transTreeToHuffmanCodeMap()
	for k, v := range mm {
		fmt.Printf("key : %c  value : %b\n", k, v)
	}
	return bytes
}

//Decode decode
func Decode(bytes []byte) []byte {

	var offset uint32

	treeLen := binary.BigEndian.Uint32(bytes[:4])
	offset += 4

	//fmt.Printf("treeLen %v\n", treeLen)

	root := buildTreeBySerialize(bytes[offset:offset+treeLen], treeLen)
	root.serializeTree()
	offset += treeLen
	bitLen := binary.BigEndian.Uint32(bytes[offset : offset+4])
	offset += 4

	result := make([]byte, 0, len(bytes)-int(offset))

	var tempLen uint32
	var bitOffset uint32
	//fmt.Printf("   len %v  data %b\n", len(bytes[offset:]), bytes[offset:])
	for tempLen < bitLen {
		b, l, o, bi := root.decodeByteFromHuffman(bytes[offset:], bitOffset)
		//fmt.Printf("\n------------------result %c  offset %v\n", b, bi)
		result = append(result, b)
		tempLen += l
		offset += o
		bitOffset = bi
	}

	return result
}
