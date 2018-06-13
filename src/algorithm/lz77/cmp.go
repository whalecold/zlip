package lz77

import (
	"fmt"
)


func genHashNumber(bytes []byte) uint16 {
	var hash uint32
	hash = uint32(bytes[0]) << 16 + uint32(bytes[1]) << 8 + uint32(bytes[2])
	//fmt.Printf("%v\n", uint16(hash & LZ77_WindowsMask))
	return uint16(hash & LZ77_WindowsMask)
}

//查看最长匹配串
func checkLargestCmpBytes(bytes []byte, curIndex , cmpIndex, maxSize uint64) uint64 {
	//fmt.Printf("cur Index %v cmpIndex %v\n", curIndex, cmpIndex)
	var length uint64
	temp := curIndex
	for {
		if curIndex >= maxSize || bytes[curIndex] != bytes[cmpIndex] || cmpIndex >= temp  {
			break
		}
		curIndex++
		cmpIndex++
		length++
	}
	//匹配的长度不要超过这个值
	if length >= LZ77_MaxCmpLength {
		return 0
	}
	return length
}

//更新prev和head的索引
func updateHashIndex(prev, head []uint64, hash uint16, index uint64) {
	temp := head[hash]
	head[hash] = index
	if temp != 0 {
		prev[index & LZ77_WindowsMask] = temp
	}
}

//更新bytes数组的前三位hash值
func updateHashBytes(bytes []byte, index uint64, prev, head []uint64) uint16 {
	//这里是更新接下来的匹配
	hash := genHashNumber(bytes[index-LZ77_MinCmpSize:index])
	updateHashIndex(prev, head, hash, index-LZ77_MinCmpSize)

	return hash
}

//这个是匹配是算法
//因为len(bytes)返回的是int 文件大小可能超过 所以长度用新的uint64参数表示
//第一个返回的map表示literal/length 出现的次数 第二个表示distance出现的次数 会对length和distance做一定的优化
//映射参考 doc里面的两张图
func Lz77Cmp(bytes []byte, size uint64) ([]byte, map[uint16]int, map[byte]int) {
	if len(bytes) < LZ77_MinCmpSize * 2 {
		panic("func cmp bytes need large than 3")
	}

	lliMap := make(map[uint16]int)
	disMap := make(map[byte]int)

	result := make([]byte, 0, 1024)
	result = append(result, bytes[:LZ77_MinCmpSize]...)
	prevIndex := make([]uint64, LZ77_CmpPrevSize)
	headIndex := make([]uint64, LZ77_CmpHeadSize)


	//小于三个字节
	var index uint64
	for index = LZ77_MinCmpSize; index + LZ77_MinCmpSize <= size;  {

		//每次移动窗口都要更新值
		updateHashBytes(bytes, index, prevIndex, headIndex)

		hash := genHashNumber(bytes[index:index + LZ77_MinCmpSize])
		cmpIndex := headIndex[hash]
		if cmpIndex == 0 { //没有匹配到
			result = append(result, bytes[index])
			index++
		} else {	//匹配到了

			var maxCmpStart uint64
			var maxCmpLength uint64
			var maxCmpNumber uint64
			//遍历hash值一致的链表
			for {

				length := checkLargestCmpBytes(bytes, index, cmpIndex, size)

				if maxCmpLength < length && index-cmpIndex < LZ77_MaxWindowsSize{
					maxCmpLength = length
					maxCmpStart = cmpIndex
				}

				cmpIndex = prevIndex[cmpIndex & LZ77_WindowsMask]
				if cmpIndex == 0 {
					break
				}
				//fmt.Printf("luup %v %v\n", cmpIndex, prevIndex[cmpIndex & LZ77_WindowsMask])
				maxCmpNumber++
				//限制匹配次数 不做判断会造成死循环
				if maxCmpNumber >= LZ77_MaxCmpNum {
					break
				}
			}

			//还是没有匹配到 或者 匹配到的是hash冲突的字段
			if maxCmpLength < LZ77_MinCmpSize {
				result = append(result, bytes[index])
				index++
			} else {
				str := fmt.Sprintf("(%v,%v)", maxCmpLength, index-maxCmpStart)
				temp := index + 1
				index += maxCmpLength

				for ; temp < index; temp++ {
					updateHashBytes(bytes, temp, prevIndex, headIndex)
				}
				result = append(result, []byte(str)...)
			}
		}
	}

	//如果还有剩余 直接打印出来 不匹配了
	for ; index < size; index++ {
		result = append(result, bytes[index])
	}
	//fmt.Printf("prev %v\n", prevIndex)
	//fmt.Printf("head %v\n", headIndex)
	return result
}
