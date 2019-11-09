package utils

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestGetLowBit32(t *testing.T) {
	num := uint32(0x2107)
	//10 0001 0000 0111
	if GetLowBit32(num, 0) != 1 ||
		GetLowBit32(num, 1) != 1 ||
		GetLowBit32(num, 2) != 1 ||
		GetLowBit32(num, 3) != 0 ||
		GetLowBit32(num, 4) != 0 ||
		GetLowBit32(num, 5) != 0 ||
		GetLowBit32(num, 6) != 0 ||
		GetLowBit32(num, 7) != 0 ||
		GetLowBit32(num, 9) != 0 ||
		GetLowBit32(num, 13) != 1 ||
		GetLowBit32(num, 8) != 1 {

		t.Errorf("GetLowBit32 falied")
	}
}

func TestGetHighBit8(t *testing.T) {
	var i int
	fmt.Printf("sizeof %v\n", unsafe.Sizeof(i))
	//0xB6
	//1011 0110
	if GetHighBit8(0xB6, 0) != 1 ||
		GetHighBit8(0xB6, 1) != 0 ||
		GetHighBit8(0xB6, 2) != 1 ||
		GetHighBit8(0xB6, 3) != 1 ||
		GetHighBit8(0xB6, 4) != 0 ||
		GetHighBit8(0xB6, 5) != 1 ||
		GetHighBit8(0xB6, 6) != 1 ||
		GetHighBit8(0xB6, 7) != 0 {
		t.Errorf("GetHighBit8 failed")
	}
}

func TestGetUint16FromBytes(t *testing.T) {
	bytes := []byte{255, 127, 98}
	if value, byteLen, offset := GetUint16FromBytes(bytes, 0, 8); value != 255 || byteLen != 1 || offset != 0 {
		t.Errorf("GetUint16FromBytes failed %v %v %v", value, byteLen, offset)
	}
	if value, byteLen, offset := GetUint16FromBytes(bytes, 0, 9); value != 510 || byteLen != 1 || offset != 1 {
		t.Errorf("GetUint16FromBytes failed %v %v %v", value, byteLen, offset)
	}
	if value, byteLen, offset := GetUint16FromBytes(bytes, 0, 11); value != 2043 || byteLen != 1 || offset != 3 {
		t.Errorf("GetUint16FromBytes failed %v %v %v", value, byteLen, offset)
	}
	if value, byteLen, offset := GetUint16FromBytes(bytes, 1, 8); value != 254 || byteLen != 1 || offset != 1 {
		t.Errorf("GetUint16FromBytes failed %v %v %v", value, byteLen, offset)
	}
	if value, byteLen, offset := GetUint16FromBytes(bytes, 3, 11); value != 2015 || byteLen != 1 || offset != 6 {
		t.Errorf("GetUint16FromBytes failed %v %v %v", value, byteLen, offset)
	}
}
