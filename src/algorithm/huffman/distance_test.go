package huffman

import (
	"testing"
	"math/rand"
	"time"
)

func TestDeflateDisTree(t *testing.T) {
	//dis := []uint16{4, 4, 4, 4, 5, 6, 5, 7, 7, 9}
	//dis := []uint16{4, 4, 4, 4, 5, 6, 5, 7, 7, 9}
	//dis := []uint16{67, 2231, 222, 1212, 991, 4, 4, 4, 4, 5, 6, 5, 35, 99, 99, 99, 33, 31, 90}
	dis := make([]uint16, 0, 32)

	rand.Seed(time.Now().Unix())
	//fmt.Printf("%v \n", rand.Uint32())
	for i := 0; i < 10000; i++ {
		dis = append(dis, uint16(rand.Uint32() % 32768))
	}
	tree := &DeflateDisTree{}
	tree.Init()
	for _, value := range dis {
		tree.AddDisElement(value)
	}
	tree.BuildTree()
	tree.BuildMap()
	tree.BuildBitsStream()
	//tree.Print()
	//fmt.Printf("tree : %v\n", tree.bits)

	newTree := DeflateDisTree{}
	newTree.Init()
	newTree.BuildCodeMapByBits(tree.bits)
	//newTree.Print()

	if false == newTree.Equal(tree) {
		t.Error("过程不对")
	}
	//str := fmt.Sprintf("%b\n", 11)
	//fmt.Printf("result %v--\n", []byte(str))
}
