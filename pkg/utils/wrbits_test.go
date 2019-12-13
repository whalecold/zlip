package utils

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestGetLowBit32(t *testing.T) {
	tests := []struct {
		in       uint
		expected byte
	}{
		{0, 1},
		{1, 1},
		{2, 1},
		{3, 0},
		{4, 0},
		{5, 0},
		{6, 0},
		{7, 0},
		{13, 1},
		{8, 1},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			// 10 0001 0000 0111
			if out := GetLowBit32(0x2107, tt.in); out != tt.expected {
				t.Errorf("expected (%v) but got (%v)", tt.expected, out)
			}
		})
	}
}

func TestGetHighBit8(t *testing.T) {
	var i int
	fmt.Printf("sizeof %v\n", unsafe.Sizeof(i))

	tests := []struct {
		in       uint32
		expected byte
	}{
		{0, 1},
		{1, 0},
		{2, 1},
		{3, 1},
		{4, 0},
		{5, 1},
		{6, 1},
		{7, 0},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			// 0xB6
			// 1011 0110
			if out := GetHighBit8(0xB6, tt.in); out != tt.expected {
				t.Errorf("expected (%v) but got (%v)", tt.expected, out)
			}
		})
	}
}

func TestGetUint16FromBytes(t *testing.T) {
	bytes := []byte{255, 127, 98}
	type output struct {
		result uint16
		bsl    uint32
		offset uint32
	}
	tests := []struct {
		name string
		bf   uint32
		len  uint16
		out  output
	}{
		{"0", 0, 8, output{255, 1, 0}},
		{"1", 0, 9, output{510, 1, 1}},
		{"2", 0, 11, output{2043, 1, 3}},
		{"3", 1, 8, output{254, 1, 1}},
		{"4", 3, 11, output{2015, 1, 6}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, byteLen, offset := GetUint16FromBytes(bytes, tt.bf, tt.len)
			if value != tt.out.result || byteLen != tt.out.bsl || offset != tt.out.offset {
				t.Errorf("expected (%v %v %v) but got (%v %v %v)",
					tt.out.result, tt.out.bsl, tt.out.offset, value, byteLen, offset)
			}
		})
	}
}
