package huffman

import (
	"fmt"
	"encoding/binary"
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

	for _, value := range in{
		fmt.Printf("%v\t", value)
	}
	fmt.Printf("\n")

	result := buildTreeBySlice(pre, in)
	pre111 := result.genStreamByInorder()
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

func EnCode(bytes []byte) []byte {

	root := buildHuffmanTree(bytes)
	serial := root.serializeTree()
	m := root.transTreeToHuffmanCodeMap()

	treeLen := len(serial)
	fmt.Printf("treeLen %v\n", treeLen)
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
			//0异或任何值的时候得到的结构都为另外一个值 尾部肯定为0 所谓得到的结果尾部也肯定为v
			tempBype ^= v
			//if v == 1 {
			//	tempBype |= 0x1
			//} else {
			//	tempBype &=  ^byte(0x1)
			//}
			bitIndex++
			if bitIndex == 8 {
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

func Decode(bytes []byte) []byte {

	var offset uint32

	treeLen := binary.BigEndian.Uint32(bytes[:4])
	offset += 4

	//fmt.Printf("treeLen %v\n", treeLen)

	root := buildTreeBySerialize(bytes[offset:offset + treeLen], treeLen)
	root.serializeTree()
	offset += treeLen
	bitLen := binary.BigEndian.Uint32(bytes[offset:offset+4])
	offset += 4

	result := make([]byte, 0, len(bytes) - int(offset))


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