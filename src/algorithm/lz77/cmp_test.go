package lz77

import "testing"

func TestCmp(t *testing.T) {
	Cmp([]byte("1232131"))
	Cmp([]byte("122232111231"))
}
