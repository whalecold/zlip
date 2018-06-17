package huffman

//这个表包括literal 和 length
type Literal struct {

}

func (l *Literal)GetZoneData(liter uint16,
							data [][]uint16,
							length bool)  (uint16, uint16, uint16){

	if length == false {
		return liter, 0, liter
	} else {
		return getZoneByData(liter, data)
	}
}

func (l *Literal)GetSourceCode(code uint16, data [][]uint16)  (uint16, uint16, bool) {
	if code <= HUFFMAN_LiteralLimit {
		return 0, code, false
	} else {
		p1, p2 := getDataByZone(code, data)
		return p1, p2, true
	}
}
