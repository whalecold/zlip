package lz77

import (
	"github.com/whalecold/zlip/pkg/huffman"
	"github.com/whalecold/zlip/pkg/utils"
)

// 先进行第一步压缩 生成cl1 cl2 和 压缩后的bits
func compressCl(bytes []byte, size uint64) ([]byte, []byte, []byte) {

	result := make([]uint16, 0, ChunkSize)
	prevIndex := make([]uint64, CmpPrevSize)
	headIndex := make([]uint64, CmpHeadSize)

	//distance
	cl1 := &huffman.Alg{}
	cl1.InitDis()
	//literal/length
	cl2 := &huffman.Alg{}
	cl2.InitLiteral()

	for i := 0; i < MinCmpSize; i++ {
		result = append(result, uint16(bytes[i]))
		cl2.AddElement(uint16(bytes[i]), false)
	}
	//bytes = append(bytes, LZ77_EndFlag)

	cl2.AddElement(huffman.EndFlag, false)
	//小于三个字节
	var index uint64
	for index = MinCmpSize; index+MinCmpSize <= size; {

		//每次移动窗口都要更新值
		updateHashBytes(bytes, index, prevIndex, headIndex)

		hash := genHashNumber(bytes[index : index+MinCmpSize])
		cmpIndex := headIndex[hash]
		if cmpIndex == 0 { //没有匹配到
			result = append(result, uint16(bytes[index]))
			cl2.AddElement(uint16(bytes[index]), false)
			index++
		} else { //匹配到了

			var maxCmpStart uint64
			var maxCmpLength uint64
			var maxCmpNumber uint64
			//遍历hash值一致的链表
			for {

				length := checkLargestCmpBytes(bytes, index, cmpIndex, size)

				if maxCmpLength < length && index-cmpIndex < MaxWindowsSize {
					maxCmpLength = length
					maxCmpStart = cmpIndex
				}

				cmpIndex = prevIndex[cmpIndex&WindowsMask]
				if cmpIndex == 0 {
					break
				}
				//fmt.Printf("luup %v %v\n", cmpIndex, prevIndex[cmpIndex & LZ77_WindowsMask])
				maxCmpNumber++
				//限制匹配次数 不做判断会造成死循环
				if maxCmpNumber >= MaxCmpNum {
					break
				}
			}

			//还是没有匹配到 或者 匹配到的是hash冲突的字段
			if maxCmpLength < MinCmpSize {
				result = append(result, uint16(bytes[index]))
				cl2.AddElement(uint16(bytes[index]), false)
				index++
			} else {
				tempLength := uint16(maxCmpLength)
				cl2.AddElement(tempLength, true)
				//fmt.Printf("pre length %v------\n", tempLength)
				utils.SetHighBit16(&tempLength, 0, 1)
				//fmt.Printf("----------------------pre length %v\n", uint16(index-maxCmpStart))
				result = append(result, tempLength)
				result = append(result, uint16(index-maxCmpStart))
				cl1.AddElement(uint16(index-maxCmpStart), false)
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
		cl2.AddElement(uint16(bytes[index]), false)
	}

	cl2.BuildHuffmanMap()

	cl1.BuildHuffmanMap()
	huffmanCode := make([]byte, 1, 1024)
	//var offset uint64
	var bits uint32
	var indexCode uint64
	var bit uint32

	result = append(result, huffman.EndFlag)
	for i := 0; i < len(result); i++ {
		//表示长度
		if utils.GetHighBit16(result[i], 0) == 1 {
			utils.SetHighBit16(&result[i], 0, 0)
			//fmt.Printf("new ----- %v\n", temp)
			bit = cl2.EnCodeElement(result[i], &huffmanCode, bits, &indexCode, true)
			bits = bit

			i++
			bit = cl1.EnCodeElement(result[i], &huffmanCode, bits, &indexCode, false)
			bits = bit
		} else {
			bit = cl2.EnCodeElement(result[i], &huffmanCode, bits, &indexCode, false)
			bits = bit
		}
	}
	cl2Bits, _ := cl2.SerializeBitsStream()
	cl1Bits, _ := cl1.SerializeBitsStream()
	return cl1Bits, cl2Bits, huffmanCode
}

func compressCClSub(cl []byte, huff *huffman.Alg) []byte {
	var bits uint32
	var indexCode uint64
	var bit uint32
	huffmanCode := make([]byte, 1, len(cl))
	for i := 0; i < len(cl); i++ {
		bit = huff.EnCodeElement(uint16(cl[i]), &huffmanCode, bits, &indexCode, false)
		bits = bit
	}
	return huffmanCode
}

//把cl1 cl2 进行第二部压缩 返回ccl 和bits
func compressCCl(cl1 []byte, cl2 []byte) ([]byte, []byte, []byte) {
	ccl := &huffman.Alg{}
	ccl.InitCCL()
	for _, value := range cl1 {
		ccl.AddElement(uint16(value), false)
	}
	for _, value := range cl2 {
		ccl.AddElement(uint16(value), false)
	}
	ccl.AddElement(huffman.CCLEndFlag, false)

	ccl.BuildHuffmanMap()

	cl1 = append(cl1, huffman.CCLEndFlag)
	cl2 = append(cl2, huffman.CCLEndFlag)

	sq1 := compressCClSub(cl1, ccl)
	sq2 := compressCClSub(cl2, ccl)
	huff, _ := ccl.SerializeBitsStream()
	return sq1, sq2, huff
}
