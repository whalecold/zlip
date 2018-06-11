package huffman

import (
	"fmt"
)


func Decode(bytes []byte) []byte {
	root := buildHuffmanTree(bytes)
	//root1 := buildHuffmanTree(bytes)
	//root2 := buildHuffmanTree(bytes)
	m := root.transTreeToHuffmanCodeMap()
	for k, v := range m {
		fmt.Printf("key : %c  value : %b\n", k, v)
	}
	fmt.Printf("-------------next-----------------\n")
	headTest := buildHuffmanTree(bytes)

	pre := headTest.genSliceByPreorder()
	in := headTest.genSliceByInorder()

	for _, value := range in{
		fmt.Printf("%v\t", value)
	}
	fmt.Printf("\n")

	result := buildTreeBySlice(pre, in)
	pre111 := result.genSliceByInorder()
	for _, value := range pre111{
		fmt.Printf("%v\t", value)
	}
	fmt.Printf("\n")

	mm := result.transTreeToHuffmanCodeMap()
	for k, v := range mm {
		fmt.Printf("key : %c  value : %b\n", k, v)
	}
	return bytes
}