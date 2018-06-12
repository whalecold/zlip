package lz77

import (
	"hash"
	"encoding/binary"
	"fmt"
	"crypto/sha1"
)


var hashIns hash.Hash

func genHashNumber(bytes []byte) uint16 {
	ret := hashIns.Sum(bytes)
	ret = ret[:8]
	fmt.Printf("%v   %v  %x\n", bytes, ret, binary.LittleEndian.Uint64(ret))
	return uint16(binary.LittleEndian.Uint64(ret) & uint64(LZ77_WindowsMask))
}

//这个是匹配是算法

func Cmp(bytes []byte) {
	if len(bytes) < 3 {
		panic("func cmp bytes need large than 3")
	}
	hashIns = sha1.New()
	index := genHashNumber(bytes[:3])

	fmt.Printf(" hash %v \n", index)

	//prevIndex := make([]uint64, LZ77_CmpPrevSize)
	//headIndex := make([]uint64, LZ77_CmpPrevSize)
	//
	////小于三个字节
	//for i := 3; i < len(bytes); i++ {
	//
	//}
}
