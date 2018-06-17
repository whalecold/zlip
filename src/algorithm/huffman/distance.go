package huffman

type Distance struct {

}

func (d *Distance)GetZoneData(distance uint16,
								data [][]uint16,
								length bool)  (uint16, uint16, uint16){

	return getZoneByData(distance, data)
}

func (d *Distance)GetSourceCode(code uint16, data [][]uint16)  (uint16, uint16, bool) {
	p1, p2 := getDataByZone(code, data)
	return p1, p2, false
}