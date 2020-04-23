package huffman

//Literal literal
//这个表包括literal 和 length
type Literal struct {
	extraCode [][]uint16 //码表
}

//GetZoneData get zone data
func (l *Literal) GetZoneData(liter uint16, length bool) (uint16, uint16, uint16) {
	if !length {
		return liter, 0, liter
	}
	return getZoneByData(liter, l.extraCode)
}

//GetSourceCode get source code
func (l *Literal) GetSourceCode(code uint16) (uint16, uint16, bool) {
	if code <= LiteralBoundary {
		return 0, code, false
	}
	p1, p2 := getDataByZone(code, l.extraCode)
	return p1, p2, true
}

//GetBitsLen get bits len
func (l *Literal) GetBitsLen() int {
	return 256 + 1 + len(l.extraCode)
}
