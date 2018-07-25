package huffman

import (
	"fmt"
	"testing"
)

type TestNode struct {
	Data   uint16
	Length bool
}

func TestLiteralTree1(t *testing.T) {
	//dis := []uint16{4, 4, 4, 4, 5, 6, 5, 7, 7, 9}
	//dis := []uint16{4, 4, 4, 4, 5, 6, 5, 7, 7, 9}
	dis := []uint16{'h', 'e', 'l', 'l'}
	tree := &DeflateTree{condition: &Literal{extraCode: LengthZone}}
	tree.Init()
	for _, value := range dis {
		tree.AddElement(value, false)
	}
	tree.AddElement(3, true)
	tree.AddElement(100, true)
	tree.BuildTree()
	tree.BuildMap()
	tree.SerializeBitsStream()
	tree.Print()
	//fmt.Printf("tree : %v\n", tree.bits)

	newTree := &DeflateTree{condition: &Literal{extraCode: LengthZone}}
	newTree.Init()

	newTree.UnSerializeBitsStream(tree.bits)
	newTree.BuildTreeByMap()

	code := make([]byte, 1, 32)
	var offset uint64
	tree.EnCodeElement(100, &code, 0, &offset, true)

	getData, _, _, _ := newTree.DecodeEle(code, 0)
	fmt.Printf("get c %v\n", getData)
}

func TestLiteralTree2(t *testing.T) {
	//dis := []uint16{4, 4, 4, 4, 5, 6, 5, 7, 7, 9}
	//dis := []uint16{4, 4, 4, 4, 5, 6, 5, 7, 7, 9}
	literalSlice := []*TestNode{{'1', false}, {'h', false},
		{'e', false},
		{126, true}, {11, true}, {'w', false},
		{256, false}}
	tree := &DeflateTree{condition: &Literal{extraCode: LengthZone}}
	tree.Init()

	for _, value := range literalSlice {
		tree.AddElement(value.Data, value.Length)
	}
	tree.BuildTree()
	tree.BuildMap()
	tree.SerializeBitsStream()
	tree.Print()
	//fmt.Printf("tree : %v\n", tree.bits)

	newTree := &DeflateTree{condition: &Literal{extraCode: LengthZone}}
	newTree.Init()

	newTree.UnSerializeBitsStream(tree.bits)
	newTree.BuildTreeByMap()

	code := make([]byte, 1, 32)
	//var offset uint64
	var bits uint32
	var indexCode uint64
	for _, value := range literalSlice {
		bit := tree.EnCodeElement(value.Data, &code, bits, &indexCode, value.Length)
		bits = bit
	}

	var resubyteoffset uint32
	var bitoffset uint32
	for {
		getData, r, b, l := newTree.DecodeEle(code[resubyteoffset:], bitoffset)
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
