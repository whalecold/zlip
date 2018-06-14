package huffman

import (
	"testing"
)

func TestReadBit(t *testing.T) {
	num := 0x2107
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
