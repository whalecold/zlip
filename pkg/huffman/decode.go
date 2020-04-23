package huffman

import (
	"encoding/binary"
)

//Decode decode
func Decode(bytes []byte) []byte {

	var offset uint32

	treeLen := binary.BigEndian.Uint32(bytes[:4])
	offset += 4

	root := buildTreeBySerialize(bytes[offset:offset+treeLen], treeLen)
	root.serializeTree()
	offset += treeLen
	bitLen := binary.BigEndian.Uint32(bytes[offset : offset+4])
	offset += 4

	result := make([]byte, 0, len(bytes)-int(offset))

	var tempLen uint32
	var bitOffset uint32
	for tempLen < bitLen {
		b, l, o, bi := root.decodeByteFromHuffman(bytes[offset:], bitOffset)
		result = append(result, b)
		tempLen += l
		offset += o
		bitOffset = bi
	}

	return result
}
