package huffman

import (
	"fmt"
)

//deflate 构建distance树

//{distance, code, bits}
var distanceZone = [][]int{
					{1, 0, 0}, {2, 0, 1}, {3, 0, 2}, {4, 0, 3}, {6, 1, 4}, {8, 1, 5},
					{12, 2, 6}, {16, 2, 7}, {24, 3, 8}, {32, 3, 9},
					{48, 4, 10}, {64, 4, 11}, {96, 5, 12},
					{128, 5, 13}, {192, 6, 14}, {256, 6, 15}, {384, 7, 16}, {512, 7, 17}, {768, 8, 18},
					{1024, 8, 19}, {1536, 9, 20}, {2048, 9, 21}, {3072, 10, 22}, {4096, 10, 23}, {6144, 11, 24}, {8192, 11, 25},
					{12288, 12, 26}, {16384, 12, 27}, {24576, 13, 28}, {32768, 13, 29}}


func init() {

	//var distanceZone := [[1, 0], [2, 1]]
	//s := [1, 2]
}

//{code, bits}
func getZoneByDistance(distance int) (int, int){
	for _, value := range distanceZone {
		if distance <= value[0] {
			return value[2], value[1]
		}
	}
	panic(fmt.Sprintf("getZoneByDistance : error param %v", distance))
}