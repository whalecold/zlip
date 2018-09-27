package huffman

//Distance distance
type Distance struct {
	extraCode [][]uint16 //码表
}

//GetZoneData get zone data
func (d *Distance) GetZoneData(distance uint16, length bool) (uint16, uint16, uint16) {
	return getZoneByData(distance, d.extraCode)
}

//GetSourceCode source code
func (d *Distance) GetSourceCode(code uint16) (uint16, uint16, bool) {
	p1, p2 := getDataByZone(code, d.extraCode)
	return p1, p2, false
}

//GetBitsLen get bits len
func (d *Distance) GetBitsLen() int {
	return len(d.extraCode)
}
