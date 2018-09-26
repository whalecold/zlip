package huffman

type CCl struct {
}

func (l *CCl) GetZoneData(liter uint16, length bool) (uint16, uint16, uint16) {

	return liter, 0, liter
}

func (l *CCl) GetSourceCode(code uint16) (uint16, uint16, bool) {
	return 0, code, false
}

func (l *CCl) GetBitsLen() int {
	return HUFFMAN_CCLLen + 2
}
