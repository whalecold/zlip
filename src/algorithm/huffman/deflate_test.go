package huffman

import (
	"testing"
)

func TestReadBit(t *testing.T) {
	num := 0x2107
	//10 0001 0000 0111
	if ReadBit(num, 0) != 1 ||
		ReadBit(num, 1) != 1 ||
		ReadBit(num, 2) != 1 ||
		ReadBit(num, 3) != 0||
		ReadBit(num, 4) != 0||
		ReadBit(num, 5) != 0||
		ReadBit(num, 6) != 0||
		ReadBit(num, 7) != 0||
		ReadBit(num, 9) != 0||
		ReadBit(num, 13) != 1||
		ReadBit(num, 8) != 1{

		t.Error("函数测试没通过") // 如果不是如预期的那么就报错
	}
}
