package huffman

import (
	"testing"
	"fmt"
)

type TestNode struct {
	Data uint16
	Length bool
}

func TestLiteralTree1(t *testing.T) {
	//dis := []uint16{4, 4, 4, 4, 5, 6, 5, 7, 7, 9}
	//dis := []uint16{4, 4, 4, 4, 5, 6, 5, 7, 7, 9}
	dis := []uint16{'h', 'e', 'l', 'l'}
	tree := &DeflateTree{}
	tree.Init(LengthZone)
	literalTree := &Literal{}
	for _, value := range dis {
		tree.AddElement(value, literalTree, false)
	}
	tree.AddElement(3, literalTree, true)
	tree.AddElement(100, literalTree, true)
	tree.BuildTree()
	tree.BuildMap()
	tree.SerializeBitsStream(literalTree)
	tree.Print()
	//fmt.Printf("tree : %v\n", tree.bits)

	newTree := &DeflateTree{}
	newTree.Init(LengthZone)

	newTree.UnSerializeBitsStream(tree.bits, literalTree)
	newTree.BuildTreeByMap()

	code := make([]byte, 1, 32)
	var offset uint64
	tree.EnCodeElement(100, &code, 0, &offset,
		literalTree, true)

	getData, _, _, _ := newTree.DecodeEle(code, 0, literalTree)
	fmt.Printf("get c %v\n", getData)
}

func TestLiteralTree2(t *testing.T) {
	//dis := []uint16{4, 4, 4, 4, 5, 6, 5, 7, 7, 9}
	//dis := []uint16{4, 4, 4, 4, 5, 6, 5, 7, 7, 9}
	literalSlice := []*TestNode{{'1', false},{'h', false},
	{'e', false},
	{126, true}, {11, true}, {'w', false},
	{256, false}}
	tree := &DeflateTree{}
	tree.Init(LengthZone)
	literalTree := &Literal{}
	for _, value := range literalSlice {
		tree.AddElement(value.Data, literalTree, value.Length)
	}
	tree.BuildTree()
	tree.BuildMap()
	tree.SerializeBitsStream(literalTree)
	tree.Print()
	//fmt.Printf("tree : %v\n", tree.bits)

	newTree := &DeflateTree{}
	newTree.Init(LengthZone)

	newTree.UnSerializeBitsStream(tree.bits, literalTree)
	newTree.BuildTreeByMap()

	code := make([]byte, 1, 32)
	//var offset uint64
	var bits uint32
	var indexCode uint64
	for _, value := range literalSlice {
		bit := tree.EnCodeElement(value.Data, &code, bits, &indexCode,
			literalTree, value.Length)
		bits = bit
	}

	var resubyteoffset uint32
	var bitoffset uint32
	for {
		getData, r, b, l:= newTree.DecodeEle(code[resubyteoffset:], bitoffset, literalTree)
		resubyteoffset += r
		bitoffset = b
		if l == true {
			fmt.Printf("length %v\n", getData)
		} else {
			fmt.Printf("length %c\n", getData)
		}

		if getData == 256 {
			fmt.Printf("end %c\n", getData)
			break
		}
		//result = append(result, getData)
	}
	//getData, _, _, _ := newTree.DecodeEle(code, 0, literalTree)
	//fmt.Printf("get c %v\n", getData)
}
