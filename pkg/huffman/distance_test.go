package huffman

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
	"unsafe"
)

func TestDeflateTree2(t *testing.T) {
	// dis := []uint16{4, 4, 4, 4, 5, 6, 5, 7, 7, 9}
	// dis := []uint16{4, 4, 4, 4, 5, 6, 5, 7, 7, 9}
	dis := []uint16{67, 2231, 222, 1212, 991, 4, 4, 4, 4, 5, 6, 5, 35, 99, 99, 99, 33, 31, 90}
	tree := &DeflateTree{condition: &Distance{extraCode: DistanceZone}}
	tree.Init()
	for _, value := range dis {
		tree.AddElement(value, false)
	}
	tree.BuildTree()
	tree.BuildMap()
	tree.SerializeBitsStream()
	// tree.Print()
	// fmt.Printf("tree : %v\n", tree.bits)

	newTree := &DeflateTree{condition: &Distance{extraCode: DistanceZone}}
	newTree.Init()
	newTree.UnSerializeBitsStream(tree.bits)
	newTree.BuildTreeByMap()
	newTree.SerializeBitsStream()

	newTree2 := DeflateTree{condition: &Distance{extraCode: DistanceZone}}
	newTree2.Init()
	newTree2.UnSerializeBitsStream(tree.bits)
	// newTree2.Print()
	// newTree.Print()

	if !newTree.Equal(tree) {
		t.Error("equal failed")
	}

	if !newTree2.Equal(newTree) {
		t.Error("equal failed")
	}

	fmt.Printf("size : %v\n", unsafe.Sizeof(uint16(10)))
}

// 测试树的序列化和反序列化 已经对一个distance的解码与反解码
func TestDeflateTreeRandom(t *testing.T) {
	// dis := []uint16{4, 4, 4, 4, 5, 6, 5, 7, 7, 9}
	// dis := []uint16{4, 4, 4, 4, 5, 6, 5, 7, 7, 9}
	// dis := []uint16{67, 2231, 222, 1212, 991, 4, 4, 4, 4, 5, 6, 5, 35, 99, 99, 99, 33, 31, 90}
	length := 10000
	dis := make([]uint16, 0, length)

	rand.Seed(time.Now().Unix())
	fmt.Printf("%v \n", rand.Uint32())
	for i := 0; i < length; i++ {
		dis = append(dis, uint16(rand.Uint32()%32768))
	}
	tree := &DeflateTree{condition: &Distance{extraCode: DistanceZone}}
	tree.Init()

	for _, value := range dis {
		tree.AddElement(value, false)
	}
	tree.AddElement(3, false)
	tree.AddElement(2, false)

	tree.BuildTree()
	tree.BuildMap()
	tree.SerializeBitsStream()
	// tree.Print()
	// fmt.Printf("tree : %v\n", tree.bits)

	newTree := &DeflateTree{condition: &Distance{extraCode: DistanceZone}}
	newTree.Init()
	newTree.UnSerializeBitsStream(tree.bits)
	newTree.BuildTreeByMap()
	newTree.SerializeBitsStream()
	// newTree.Print()

	newTree2 := DeflateTree{condition: &Distance{extraCode: DistanceZone}}
	newTree2.Init()
	newTree2.UnSerializeBitsStream(tree.bits)
	// newTree2.Print()
	if !newTree.Equal(tree) {
		t.Error("not match")
	}
	if !newTree2.Equal(newTree) {
		t.Error("not match")
	}

	index := rand.Uint32() % uint32(len(dis))
	testDis := dis[index]
	// testDis = 3
	code := make([]byte, 1, 32)

	var offset uint64
	newTree.EnCodeElement(testDis, &code, 0, &offset, false)
	fmt.Printf("offset ..  %v\n", offset)
	getData, _, _, _ := newTree.DecodeEle(code, 0)

	if getData != testDis {
		t.Error("data not match")
	}
}

// 测试一个串的数字解码与反解码
func TestDeflateTree3(t *testing.T) {
	length := 10000
	dis := make([]uint16, 0, length)

	rand.Seed(time.Now().Unix())
	//fmt.Printf("%v \n", rand.Uint32())
	for i := 0; i < length; i++ {
		dis = append(dis, uint16(rand.Uint32()%32768))
	}
	tree := &DeflateTree{condition: &Distance{extraCode: DistanceZone}}
	tree.Init()
	for _, value := range dis {
		tree.AddElement(value, false)
	}
	tree.AddElement(3, false)
	tree.AddElement(2, false)
	tree.AddElement(4, false)
	tree.AddElement(5, false)
	tree.BuildTree()
	tree.BuildMap()
	tree.SerializeBitsStream()

	newTree := &DeflateTree{condition: &Distance{extraCode: DistanceZone}}
	newTree.Init()
	newTree.UnSerializeBitsStream(tree.bits)
	newTree.BuildTreeByMap()
	newTree.SerializeBitsStream()

	code := make([]byte, 1, 32)
	testDis := make([]uint16, 0, 12)
	times := 20
	testDis = append(testDis, 3)
	testDis = append(testDis, 2)
	testDis = append(testDis, 2)
	testDis = append(testDis, 4)
	testDis = append(testDis, 5)
	testDis = append(testDis, 2)
	for i := 0; i < times; i++ {
		index := rand.Uint32() % uint32(len(dis))
		testDis = append(testDis, dis[index])
	}

	fmt.Printf("%b  | data : %v\n", code, testDis)

	var indexCode uint64
	var bits uint32
	for _, value := range testDis {
		bit := tree.EnCodeElement(value, &code, bits, &indexCode, false)
		bits = bit
	}

	result := make([]uint16, 0, 32)
	var resultByteOffset uint32
	var bitOffset uint32
	for i := 0; i < times; i++ {
		getData, r, b, _ := newTree.DecodeEle(code[resultByteOffset:], bitOffset)
		resultByteOffset += r
		bitOffset = b
		result = append(result, getData)
	}

	for index, value := range result {
		if value != testDis[index] {
			t.Errorf("complier error！！！")
			break
		}
	}
}

func TestDistanceGetMaxLen(t *testing.T) {
	literalAlgorithm := &Alg{}
	literalAlgorithm.InitLiteral()
	rand.Seed(time.Now().Unix())

	for i := 0; i <= 285; i++ {
		// randomTimes := rand.Uint32() % 200
		// for j := uint32(0); j < randomTimes; j++ {
		for j := 0; j < i*10000; j++ {
			literalAlgorithm.AddElement(uint16(i), false)
		}
	}
	literalAlgorithm.AddElement(uint16(9), false)
	literalAlgorithm.BuildHuffmanMap()
	literalAlgorithm.SerializeBitsStream()
}
