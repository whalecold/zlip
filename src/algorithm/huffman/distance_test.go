package huffman

import (
	"testing"
	"math/rand"
	"fmt"
	"time"
	"unsafe"
)

func TestDeflateTree2(t *testing.T) {
	//dis := []uint16{4, 4, 4, 4, 5, 6, 5, 7, 7, 9}
	//dis := []uint16{4, 4, 4, 4, 5, 6, 5, 7, 7, 9}
	dis := []uint16{67, 2231, 222, 1212, 991, 4, 4, 4, 4, 5, 6, 5, 35, 99, 99, 99, 33, 31, 90}
	tree := &DeflateTree{}
	tree.Init(DistanceZone)
	disTree := &Distance{}
	for _, value := range dis {
		tree.AddElement(value, disTree, false)
	}
	tree.BuildTree()
	tree.BuildMap()
	tree.SerializeBitsStream(disTree)
	//tree.Print()
	//fmt.Printf("tree : %v\n", tree.bits)

	newTree := &DeflateTree{}
	newTree.Init(DistanceZone)
	newTree.UnSerializeBitsStream(tree.bits, disTree)
	newTree.BuildTreeByMap()
	newTree.SerializeBitsStream(disTree)

	newTree2 := DeflateTree{}
	newTree2.Init(DistanceZone)
	newTree2.UnSerializeBitsStream(tree.bits, disTree)
	//newTree2.Print()
	//newTree.Print()

	if false == newTree.Equal(tree) {
		t.Error("过程不对")
	}

	if false == newTree2.Equal(newTree) {
		t.Error("过程不对")
	}

	fmt.Printf("size : %v\n", unsafe.Sizeof(uint16(10)))
}

//测试树的序列化和反序列化 已经对一个distance的解码与反解码
func TestDeflateTreeRandom(t *testing.T) {
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
	tree := &DeflateTree{}
	tree.Init(DistanceZone)
	disTree := &Distance{}
	for _, value := range dis {
		tree.AddElement(value, disTree, false)
	}
	tree.AddElement(3, disTree, false)
	tree.AddElement(2, disTree, false)

	tree.BuildTree()
	tree.BuildMap()
	tree.SerializeBitsStream(disTree)
	//tree.Print()
	//fmt.Printf("tree : %v\n", tree.bits)

	newTree := &DeflateTree{}
	newTree.Init(DistanceZone)
	newTree.UnSerializeBitsStream(tree.bits, disTree)
	newTree.BuildTreeByMap()
	newTree.SerializeBitsStream(disTree)
	//newTree.Print()

	newTree2 := DeflateTree{}
	newTree2.Init(DistanceZone)
	newTree2.UnSerializeBitsStream(tree.bits, disTree)
	//newTree2.Print()
	if false == newTree.Equal(tree) {
		t.Error("过程不对")
	}
	if false == newTree2.Equal(newTree) {
		t.Error("过程不对")
	}

	index := rand.Uint32() % uint32(len(dis))
	testDis := dis[index]
	//testDis = 3
	code := make([]byte, 1, 32)

	var offset uint64
	newTree.EnCodeElement(testDis, &code, 0, &offset,
		disTree, false)
	fmt.Printf("offset ..  %v\n", offset)
	getData, _, _, _ := newTree.DecodeEle(code, 0, disTree)

	if getData != testDis {
		t.Error("过程不对")
	}
}

//测试一个串的数字解码与反解码
func TestDeflateTree3(t *testing.T)  {
	length := 10000
	dis := make([]uint16, 0, length)

	rand.Seed(time.Now().Unix())
	//fmt.Printf("%v \n", rand.Uint32())
	for i := 0; i < length; i++ {
		dis = append(dis, uint16(rand.Uint32() % 32768))
	}
	tree := &DeflateTree{}
	tree.Init(DistanceZone)
	disTree := &Distance{}
	for _, value := range dis {
		tree.AddElement(value, disTree, false)
	}
	tree.AddElement(3, disTree, false)
	tree.AddElement(2, disTree, false)
	tree.AddElement(4, disTree, false)
	tree.AddElement(5, disTree, false)
	tree.BuildTree()
	tree.BuildMap()
	tree.SerializeBitsStream(disTree)

	newTree := &DeflateTree{}
	newTree.Init(DistanceZone)
	newTree.UnSerializeBitsStream(tree.bits, disTree)
	newTree.BuildTreeByMap()
	newTree.SerializeBitsStream(disTree)

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
		bit := tree.EnCodeElement(value, &code, bits, &indexCode,
					disTree, false)
		bits = bit
	}


	result := make([]uint16, 0, 32)
	var resubyteoffset uint32
	var bitoffset uint32
	for i := 0; i < times; i++ {
		getData, r, b, _:= newTree.DecodeEle(code[resubyteoffset:], bitoffset, disTree)
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