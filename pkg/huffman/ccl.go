package huffman

// CCl no use
type CCl struct {
}

// GetZoneData nil
func (l *CCl) GetZoneData(liter uint16, length bool) (uint16, uint16, uint16) {

	return liter, 0, liter
}

//GetSourceCode nil
func (l *CCl) GetSourceCode(code uint16) (uint16, uint16, bool) {
	return 0, code, false
}

//GetBitsLen nil
func (l *CCl) GetBitsLen() int {
	return HUFFMANCCLLen + 2
}
