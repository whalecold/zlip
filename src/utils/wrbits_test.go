package utils

import (
	"testing"
	"fmt"
	"unsafe"
)

func TestReadBitLow(t *testing.T) {
	num := uint32(0x2107)
	//10 0001 0000 0111
	if ReadBitLow(num, 0) != 1 ||
		ReadBitLow(num, 1) != 1 ||
		ReadBitLow(num, 2) != 1 ||
		ReadBitLow(num, 3) != 0||
		ReadBitLow(num, 4) != 0||
		ReadBitLow(num, 5) != 0||
		ReadBitLow(num, 6) != 0||
		ReadBitLow(num, 7) != 0||
		ReadBitLow(num, 9) != 0||
		ReadBitLow(num, 13) != 1||
		ReadBitLow(num, 8) != 1{

		t.Error("函数测试没通过") // 如果不是如预期的那么就报错
	}
}

func TestReadBitsHigh(t *testing.T) {
	var i int
	fmt.Printf("sizeof %v\n", unsafe.Sizeof(i))
	//0xB6
	//1011 0110
	if ReadBitsHigh(0xB6, 0) != 1 ||
		ReadBitsHigh(0xB6, 1) != 0 ||
		ReadBitsHigh(0xB6, 2) != 1 ||
		ReadBitsHigh(0xB6, 3) != 1 ||
		ReadBitsHigh(0xB6, 4) != 0 ||
		ReadBitsHigh(0xB6, 5) != 1 ||
		ReadBitsHigh(0xB6, 6) != 1 ||
		ReadBitsHigh(0xB6, 7) != 0 {
		t.Error("函数测试没通过")
	}
}

