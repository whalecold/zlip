package huffman

import "fmt"

// deflate 构建distance树

// DistanceZone distance zone data
// {min max distance, bits, code}
var DistanceZone = [][]uint16{
	{1, 1, 0, 0}, {2, 2, 0, 1}, {3, 3, 0, 2}, {4, 4, 0, 3},
	{5, 6, 1, 4}, {7, 8, 1, 5},
	{9, 12, 2, 6}, {13, 16, 2, 7}, {17, 24, 3, 8}, {25, 32, 3, 9},
	{33, 48, 4, 10}, {49, 64, 4, 11}, {65, 96, 5, 12},
	{97, 128, 5, 13}, {129, 192, 6, 14}, {193, 256, 6, 15},
	{257, 384, 7, 16}, {385, 512, 7, 17}, {513, 768, 8, 18},
	{769, 1024, 8, 19}, {1025, 1536, 9, 20}, {1537, 2048, 9, 21},
	{2049, 3072, 10, 22}, {3073, 4096, 10, 23}, {4097, 6144, 11, 24},
	{6145, 8192, 11, 25}, {8193, 12288, 12, 26},
	{12289, 16384, 12, 27}, {16385, 24576, 13, 28},
	{24577, 32768, 13, 29}}

// LengthZone length zone
// {min max length, bits, code}
var LengthZone = [][]uint16{
	{3, 3, 0, 257}, {4, 4, 0, 258}, {5, 5, 0, 259}, {6, 6, 0, 260},
	{7, 7, 0, 261}, {8, 8, 0, 262},
	{9, 9, 0, 263}, {10, 10, 0, 264}, {11, 12, 1, 265}, {13, 14, 1, 266},
	{15, 16, 1, 267}, {17, 18, 1, 268}, {19, 22, 2, 269},
	{23, 26, 2, 270}, {27, 30, 2, 271}, {31, 34, 2, 272},
	{35, 42, 3, 273}, {43, 50, 3, 274}, {51, 58, 3, 275},
	{59, 66, 3, 276}, {67, 82, 4, 277}, {83, 98, 4, 278},
	{99, 114, 4, 279}, {115, 130, 4, 280}, {131, 162, 5, 281},
	{163, 194, 5, 282}, {195, 226, 5, 283},
	{227, 257, 5, 284}, {258, 258, 0, 285}}

// getZoneByData {zone, bits lower}
func getZoneByData(distance uint16, data [][]uint16) (uint16, uint16, uint16) {
	for _, value := range data {
		if distance <= uint16(value[1]) {
			return value[3], value[2], value[0]
		}
	}
	panic(fmt.Sprintf("getZoneByDistance : error param %v "+
		"(check if init)", distance))
}

// GetZoneByDis get zone by dis
// {zone, bits lower}
func GetZoneByDis(distance uint16) (uint16, uint16, uint16) {
	return getZoneByData(distance, DistanceZone)
}

// GetZoneByLength get zone
func GetZoneByLength(distance uint16) (uint16, uint16, uint16) {
	return getZoneByData(distance, LengthZone)
}

// getDataByZone 返回 bits扩展位置 最小值
func getDataByZone(zone uint16, data [][]uint16) (uint16, uint16) {
	for _, value := range data {
		if value[3] == zone {
			return value[2], value[0]
		}
	}
	panic(fmt.Sprintf("getDistanceByZone : error param %v", zone))
}

// GetDisByData get dis by data
func GetDisByData(zone uint16) (uint16, uint16) {
	return getDataByZone(zone, DistanceZone)
}

// GetLengthByData get length
func GetLengthByData(zone uint16) (uint16, uint16) {
	return getDataByZone(zone, LengthZone)
}

func getMaxDepth(bits []byte) int {
	var max byte
	for _, value := range bits {
		if value > max {
			max = value
		}
	}
	return int(max)
}

// CompareTwoBytes comapre
func CompareTwoBytes(l, m []byte) bool {
	if len(l) != len(m) {
		return false
	}
	for index, value := range l {
		if value != m[index] {
			return false
		}
	}
	return true
}

func checkBytesFull(bytes *[]byte, offset *uint32) bool {
	if *offset == 8 {
		*offset = 0
		*bytes = append(*bytes, byte(0))
		return true
	}
	return false
}
