package lz77

import (
	//"fmt"
	"algorithm/huffman"
	"utils"
	"encoding/binary"
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
//([]byte, map[uint16]int, map[byte]int)
func Lz77Compress(bytes []byte, size uint64) []byte {
	if len(bytes) < LZ77_MinCmpSize * 2 {
		panic("func cmp bytes need large than 3")
	}

	//lliMap := make(map[uint16]int)
	//disMap := make(map[byte]int)

	result := make([]uint16, 0, 1024)


	prevIndex := make([]uint64, LZ77_CmpPrevSize)
	headIndex := make([]uint64, LZ77_CmpHeadSize)

	disDeflateTree := &huffman.DeflateTree{}
	disDeflateTree.Init(huffman.DistanceZone)
	//disDeflateTree.BuildTree()
	//disDeflateTree.BuildMap()
	literalDeflateTree := &huffman.DeflateTree{}
	literalDeflateTree.Init(huffman.LengthZone)
	//literalDeflateTree.BuildTree()
	//literalDeflateTree.BuildMap()
	distance := &huffman.Distance{}
	literal := &huffman.Literal{}

	for i := 0; i < LZ77_MinCmpSize; i++ {
		result = append(result, uint16(bytes[i]))
		literalDeflateTree.AddElement(uint16(bytes[i]), literal,
			false)
	}
	//bytes = append(bytes, LZ77_EndFlag)

	literalDeflateTree.AddElement(LZ77_EndFlag, literal,
		false)
	//小于三个字节
	var index uint64
	for index = LZ77_MinCmpSize; index + LZ77_MinCmpSize <= size;  {

		//每次移动窗口都要更新值
		updateHashBytes(bytes, index, prevIndex, headIndex)

		hash := genHashNumber(bytes[index:index + LZ77_MinCmpSize])
		cmpIndex := headIndex[hash]
		if cmpIndex == 0 { //没有匹配到
			result = append(result, uint16(bytes[index]))
			literalDeflateTree.AddElement(uint16(bytes[index]), literal,
				false)
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
				result = append(result, uint16(bytes[index]))
				literalDeflateTree.AddElement(uint16(bytes[index]), literal,
					false)
				index++
			} else {
				tempLength := uint16(maxCmpLength)
				literalDeflateTree.AddElement(tempLength, literal,
					true)
				//fmt.Printf("pre length %v------\n", tempLength)
				utils.WriteBitsHigh16(&tempLength, 0, 1)
				//fmt.Printf("----------------------pre length %v\n", uint16(index-maxCmpStart))
				result = append(result, tempLength)
				result = append(result, uint16(index-maxCmpStart))
				disDeflateTree.AddElement(uint16(index-maxCmpStart), distance,
					false)
				//str := fmt.Sprintf("(%v,%v)", maxCmpLength, index-maxCmpStart)
				temp := index + 1
				index += maxCmpLength

				for ; temp < index; temp++ {
					updateHashBytes(bytes, temp, prevIndex, headIndex)
				}
				//result = append(result, []byte(str)...)
			}
		}
	}

	//如果还有剩余 直接打印出来 不匹配了
	for ; index < size; index++ {
		result = append(result, uint16(bytes[index]))
		literalDeflateTree.AddElement(uint16(bytes[index]), literal,
			false)
	}

	literalDeflateTree.BuildTree()
	literalDeflateTree.BuildMap()

	disDeflateTree.BuildTree()
	disDeflateTree.BuildMap()
	//fmt.Printf("prev %v\n", prevIndex)
	//fmt.Printf("head %v\n", headIndex)
	huffmanCode := make([]byte, 1, 1024)
	//var offset uint64
	var bits uint32
	var indexCode uint64
	var bit uint32

	//for _, value := range result {
	//}
	result = append(result, LZ77_EndFlag)
	for i := 0; i < len(result); i++{
		//表示长度
		if utils.ReadBitsHigh16(result[i], 0) == 1 {
			utils.WriteBitsHigh16(&result[i], 0, 0)
			//fmt.Printf("new ----- %v\n", temp)
			bit = literalDeflateTree.EnCodeElement(result[i], &huffmanCode, bits, &indexCode,
				literal, true)
			bits = bit

			i++
			bit = disDeflateTree.EnCodeElement(result[i], &huffmanCode, bits, &indexCode,
				distance, false)
			bits = bit
		} else {
			bit = literalDeflateTree.EnCodeElement(result[i], &huffmanCode, bits, &indexCode,
				literal, false)
			bits = bit
		}
	}
	literalDeflateTree.SerializeBitsStream(literal)
	disDeflateTree.SerializeBitsStream(distance)
	literalLen := make([]byte, 4)
	distanceLen := make([]byte, 4)
	binary.BigEndian.PutUint32(literalLen, literalDeflateTree.BitesLen())
	binary.BigEndian.PutUint32(distanceLen, disDeflateTree.BitesLen())
	lastResult := make([]byte, 0, 8 + literalDeflateTree.BitesLen() +
		disDeflateTree.BitesLen() + uint32(len(huffmanCode)))
	lastResult = append(lastResult, literalLen...)
	lastResult = append(lastResult, distanceLen...)
	lastResult = append(lastResult, literalDeflateTree.GetBits()...)
	lastResult = append(lastResult, disDeflateTree.GetBits()...)
	lastResult = append(lastResult, huffmanCode...)
	//fmt.Printf("lastResult... len %b\n", huffmanCode)
	//fmt.Printf("lastResult... len %v\n", literalDeflateTree.GetBits())
	//fmt.Printf("lastResult... len %v\n", disDeflateTree.GetBits())
	return lastResult
}

func UnLz77Compress(bytes []byte) []byte {
	if len(bytes) < 8 {
		panic("UnLz77Compress error param len ")
	}

	var offset uint64
	literalBitsLen := binary.BigEndian.Uint32(bytes[:4])
	distanceBitsLen := binary.BigEndian.Uint32(bytes[4:8])

	disDeflateTree := &huffman.DeflateTree{}
	disDeflateTree.Init(huffman.DistanceZone)
	literalDeflateTree := &huffman.DeflateTree{}
	literalDeflateTree.Init(huffman.LengthZone)
	distance := &huffman.Distance{}
	literal := &huffman.Literal{}
	offset += 8

	literalDeflateTree.UnSerializeBitsStream(bytes[offset:offset+uint64(literalBitsLen)], literal)

	//fmt.Printf("read %v\n", bytes[offset:offset+uint64(literalBitsLen)])
	literalDeflateTree.BuildTreeByMap()
	//literalDeflateTree.Print()
	offset += uint64(literalBitsLen)
	disDeflateTree.UnSerializeBitsStream(bytes[offset:offset+uint64(distanceBitsLen)], distance)
	//fmt.Printf("read %v\n", bytes[offset:offset+uint64(distanceBitsLen)])
	disDeflateTree.BuildTreeByMap()
	//disDeflateTree.Print()
	offset += uint64(distanceBitsLen)

	lastResult := make([]byte, 0, 1024)

	buffer := bytes[offset:]
	//fmt.Printf("lastResult... len %b\n", buffer)
	var resubyteoffset uint32
	var bitoffset uint32
	for {
		getData, r, b, l:= literalDeflateTree.DecodeEle(buffer[resubyteoffset:], bitoffset, literal)
		resubyteoffset += r
		bitoffset = b
		//fmt.Printf("dara %v, %v\n", getData, l)
		if l == true {
			length := uint64(getData)
			getData, r, b, _ = disDeflateTree.DecodeEle(buffer[resubyteoffset:], bitoffset, distance)
			resubyteoffset += r
			bitoffset = b
			//fmt.Printf("dara %v, %v\n", getData, l)
			nowLen := uint64(len(lastResult))
			for i := uint64(0); i < length; i++ {
				lastResult = append(lastResult, lastResult[nowLen-uint64(getData)+i])
			}
		} else if getData == 256 {
			//fmt.Printf("end buffer \n")
			break
		} else {
			lastResult = append(lastResult, byte(getData))
		}
	}
	return lastResult
}
