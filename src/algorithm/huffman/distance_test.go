package huffman

import (
	"testing"
	"math/rand"
	"fmt"
	"time"
)

func TestDeflateDisTree2(t *testing.T) {
	//dis := []uint16{4, 4, 4, 4, 5, 6, 5, 7, 7, 9}
	//dis := []uint16{4, 4, 4, 4, 5, 6, 5, 7, 7, 9}
	dis := []uint16{67, 2231, 222, 1212, 991, 4, 4, 4, 4, 5, 6, 5, 35, 99, 99, 99, 33, 31, 90}
	tree := &DeflateDisTree{}
	tree.Init()
	for _, value := range dis {
		tree.AddDisElement(value)
	}
	tree.BuildTree()
	tree.BuildMap()
	tree.SerializeBitsStream()
	//tree.Print()
	//fmt.Printf("tree : %v\n", tree.bits)

	newTree := &DeflateDisTree{}
	newTree.Init()
	newTree.UnSerializeBitsStream(tree.bits)
	newTree.BuildTreeByMap()
	newTree.SerializeBitsStream()

	newTree2 := DeflateDisTree{}
	newTree2.Init()
	newTree2.UnSerializeBitsStream(tree.bits)
	//newTree2.Print()
	//newTree.Print()

	if false == newTree.Equal(tree) {
		t.Error("过程不对")
	}

	if false == newTree2.Equal(newTree) {
		t.Error("过程不对")
	}
}

//测试树的序列化和反序列化 已经对一个distance的解码与反解码
func TestDeflateDisTreeRandom(t *testing.T) {
	//dis := []uint16{4, 4, 4, 4, 5, 6, 5, 7, 7, 9}
	//dis := []uint16{4, 4, 4, 4, 5, 6, 5, 7, 7, 9}
	//dis := []uint16{67, 2231, 222, 1212, 991, 4, 4, 4, 4, 5, 6, 5, 35, 99, 99, 99, 33, 31, 90}
	length := 10000
	dis := make([]uint16, 0, length)

	rand.Seed(time.Now().Unix())
	fmt.Printf("%v \n", rand.Uint32())
	for i := 0; i < length; i++ {
		dis = append(dis, uint16(rand.Uint32() % 32768))
	}
	tree := &DeflateDisTree{}
	tree.Init()
	for _, value := range dis {
		tree.AddDisElement(value)
	}
	tree.BuildTree()
	tree.BuildMap()
	tree.SerializeBitsStream()
	//tree.Print()
	//fmt.Printf("tree : %v\n", tree.bits)

	newTree := &DeflateDisTree{}
	newTree.Init()
	newTree.UnSerializeBitsStream(tree.bits)
	newTree.BuildTreeByMap()
	newTree.SerializeBitsStream()
	//newTree.Print()

	newTree2 := DeflateDisTree{}
	newTree2.Init()
	newTree2.UnSerializeBitsStream(tree.bits)
	//newTree2.Print()
	if false == newTree.Equal(tree) {
		t.Error("过程不对")
	}
	if false == newTree2.Equal(newTree) {
		t.Error("过程不对")
	}

	index := rand.Uint32() % uint32(len(dis))
	testDis := dis[index]
	code := make([]byte, 1, 32)

	var offset uint64
	newTree.EnCodeDistance(testDis, &code, 0, &offset)

	getData, _, _ := newTree.DecodeDistance(code, 0)

	if getData != testDis {
		t.Error("过程不对")
	}
}

//测试一个串的数字解码与反解码
func TestDeflateDisTree3(t *testing.T)  {
	length := 10000
	dis := make([]uint16, 0, length)

	rand.Seed(time.Now().Unix())
	//fmt.Printf("%v \n", rand.Uint32())
	for i := 0; i < length; i++ {
		dis = append(dis, uint16(rand.Uint32() % 32768))
	}
	tree := &DeflateDisTree{}
	tree.Init()
	for _, value := range dis {
		tree.AddDisElement(value)
	}
	tree.BuildTree()
	tree.BuildMap()
	tree.SerializeBitsStream()

	newTree := &DeflateDisTree{}
	newTree.Init()
	newTree.UnSerializeBitsStream(tree.bits)
	newTree.BuildTreeByMap()
	newTree.SerializeBitsStream()

	code := make([]byte, 1, 32)
	testDis := make([]uint16, 0, 12)
	times := 20
	for i := 0; i < times; i++ {
		index := rand.Uint32() % uint32(len(dis))
		testDis = append(testDis, dis[index])
	}

	fmt.Printf("%b  | data : %v\n", code, testDis)

	var indexCode uint64
	var bits uint32
	for _, value := range testDis {
		bit := tree.EnCodeDistance(value, &code, bits, &indexCode)
		bits = bit
	}

	result := make([]uint16, 0, 32)
	var resubyteoffset uint32
	var bitoffset uint32
	for i := 0; i < times; i++ {
		getData, r, b := newTree.DecodeDistance(code[resubyteoffset:], bitoffset)
		resubyteoffset += r
		bitoffset = b
		result = append(result, getData)
	}

	for index, value := range result {
		if value != testDis[index] {
			t.Errorf("编解码错误！！！")
			break
		}
	}
}