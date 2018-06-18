package lz77

import (
	"testing"
	"fmt"
)

func TestRLC(t *testing.T) {
	slice := []byte{0, 0, 0, 0, 0, 0, 1, 2, 0, 0, 0, 0, 0, 4, 1, 3, 4, 0, 0, 0}
	re := RLC(slice)
	fmt.Printf("result %v\n", slice)
	fmt.Printf("result %v\n", re)
	next := UnRLC(re)
	fmt.Printf("result %v\n", next)
}
