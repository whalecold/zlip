package lz77

import (
	"encoding/binary"

	"github.com/whalecold/compress/pkg/huffman"
)

func genHashNumber(bytes []byte) uint16 {
	var hash uint32
	hash = uint32(bytes[0])<<16 + uint32(bytes[1])<<8 + uint32(bytes[2])
	//fmt.Printf("%v\n", uint16(hash & LZ77_WindowsMask))
	return uint16(hash & LZ77WindowsMask)
}

//查看最长匹配串 这个函数都是数组操作 不知道为什么会这么耗cpu~~
func checkLargestCmpBytes(bytes []byte, curIndex, cmpIndex, maxSize uint64) uint64 {
	//fmt.Printf("cur Index %v cmpIndex %v\n", curIndex, cmpIndex)
	var length uint64
	temp := curIndex
	for {
		if curIndex >= maxSize || bytes[curIndex] != bytes[cmpIndex] ||
			cmpIndex >= temp || length >= LZ77MaxCmpLength-1 {
			break
		}
		curIndex++
		cmpIndex++
		length++
	}
	////匹配的长度不要超过这个值 把这个判断放上面可以减少判断次数
	//if length >= LZ77_MaxCmpLength {
	//	return 0
	//}
	return length
}

//更新prev和head的索引
func updateHashIndex(prev, head []uint64, hash uint16, index uint64) {
	temp := head[hash]
	head[hash] = index
	if temp != 0 {
		prev[index&LZ77WindowsMask] = temp
	}
}

//更新bytes数组的前三位hash值
func updateHashBytes(bytes []byte, index uint64, prev, head []uint64) uint16 {
	//这里是更新接下来的匹配
	hash := genHashNumber(bytes[index-LZ77MinCmpSize : index])
	updateHashIndex(prev, head, hash, index-LZ77MinCmpSize)

	return hash
}

//Compress nil
//这个是匹配是算法
//因为len(bytes)返回的是int 文件大小可能超过 所以长度用新的uint64参数表示
//第一个返回的map表示literal/length 出现的次数 第二个表示distance出现的次数 会对length和distance做一定的优化
//映射参考 doc里面的两张图
//([]byte, map[uint16]int, map[byte]int)
func Compress(bytes []byte, outBuffer *[]byte, size uint64) ([]byte, int) {
	if len(bytes) < LZ77MinCmpSize*2 {
		panic("func cmp bytes need large than 3")
	}

	//进行第一步压缩 得到两个码表序列和压缩后的码流
	cl1Bits, cl2Bits, huffmanCode := compressCl(bytes, size)
	//fmt.Printf("cl1Bits %v  cl2Bits %v\n", len(cl1Bits), len(cl2Bits))
	//游程编码压缩
	sq1 := RLC(cl1Bits)
	sq2 := RLC(cl2Bits)

	//fmt.Printf("sq1 %v  sq2 %v\n", len(sq1), len(sq2))

	sq1Bits, sq2Bits, huffman3 := compressCCl(sq1, sq2)

	//fmt.Printf("sq1Bits %v  sq2Bits %v huffman3 %v\n", len(sq1Bits), len(sq2Bits), len(huffman3))

	/* 压缩格式 单位 byte
	|  huffman3Len(2) | sq1BitsLen(2) |
	|sq2BitsLen(2) | huffman3 | sq1Bits | sq2Bits | huffmanCode...|
	*/
	//headInfoLen := make([]byte, 2)
	huffman3Len := make([]byte, 2)
	sq1BitsLen := make([]byte, 2)
	sq2BitsLen := make([]byte, 2)

	//binary.BigEndian.PutUint16(headInfoLen, uint16(len(LZ77_HeadInfo)))
	binary.BigEndian.PutUint16(huffman3Len, uint16(len(huffman3)))
	binary.BigEndian.PutUint16(sq1BitsLen, uint16(len(sq1Bits)))
	binary.BigEndian.PutUint16(sq2BitsLen, uint16(len(sq2Bits)))

	//lastResult = append(lastResult, headInfoLen...)
	//lastResult = append(lastResult, []byte(LZ77_HeadInfo)...)
	*outBuffer = append(*outBuffer, huffman3Len...)
	*outBuffer = append(*outBuffer, sq1BitsLen...)
	*outBuffer = append(*outBuffer, sq2BitsLen...)

	*outBuffer = append(*outBuffer, huffman3...)
	*outBuffer = append(*outBuffer, sq1Bits...)
	*outBuffer = append(*outBuffer, sq2Bits...)

	*outBuffer = append(*outBuffer, huffmanCode...)
	return *outBuffer, len(*outBuffer)
}

//UnCompress nil
/* 压缩格式 单位 byte
|huffman3Len(2) | sq1BitsLen(2) |
|sq2BitsLen(2) | huffman3 | sq1Bits | sq2Bits | huffmanCode...|
*/
func UnCompress(bytes []byte) []byte {
	if len(bytes) < 8 {
		panic("UnLz77Compress error param len ")
	}

	var offset uint64
	//headLen := binary.BigEndian.Uint16(bytes[:2])
	//offset += 2
	//headInfo := bytes[offset:offset+uint64(headLen)]
	//fmt.Printf("head info %v\n", string(headInfo))
	//offset += uint64(headLen)
	//fmt.Printf("UnLz77Compress  %v \n", len(bytes))

	huffman3Len := binary.BigEndian.Uint16(bytes[offset : offset+2])
	offset += 2
	sq1BitsLen := binary.BigEndian.Uint16(bytes[offset : offset+2])
	offset += 2
	sq2BitsLen := binary.BigEndian.Uint16(bytes[offset : offset+2])
	offset += 2
	huffmanCode := bytes[offset : offset+uint64(huffman3Len)]
	offset += uint64(huffman3Len)
	sq1Bits := bytes[offset : offset+uint64(sq1BitsLen)]
	offset += uint64(sq1BitsLen)
	sq2Bits := bytes[offset : offset+uint64(sq2BitsLen)]
	offset += uint64(sq2BitsLen)

	//fmt.Printf("sq1Bits %v sq2Bits %v huffmanCode %v\n", len(sq1Bits), len(sq2Bits), len(huffmanCode))
	sq1Serial, sq2Serial := unCompressSQ(huffmanCode, sq1Bits, sq2Bits)
	//fmt.Printf("sq1Serial %v sq2Serial %v\n", len(sq1Serial), len(sq2Serial))
	sq1Serial = UnRLC(sq1Serial)
	sq2Serial = UnRLC(sq2Serial)
	//fmt.Printf("sq1Serial %v sq2Serial %v\n", len(sq1Serial), len(sq2Serial))

	cl1 := &huffman.Alg{}
	cl1.InitDis()
	cl2 := &huffman.Alg{}
	cl2.InitLiteral()

	cl2.UnSerializeAndBuild(sq2Serial)

	cl1.UnSerializeAndBuild(sq1Serial)
	//fmt.Printf("read %v\n", bytes[offset:offset+uint64(distanceBitsLen)])
	//disDeflateTree.Print()

	outBuffer := make([]byte, 0, LZ77ChunkSize)
	buffer := bytes[offset:]
	//fmt.Printf("lastResult... len %b\n", buffer)
	var resubyteoffset uint32
	var bitoffset uint32
	for {
		getData, r, b, l := cl2.DecodeEle(buffer[resubyteoffset:], bitoffset)
		resubyteoffset += r
		bitoffset = b
		//fmt.Printf("dara %v, %v\n", getData, l)
		if l == true {
			length := uint64(getData)
			getData, r, b, _ = cl1.DecodeEle(buffer[resubyteoffset:], bitoffset)
			resubyteoffset += r
			bitoffset = b
			//fmt.Printf("dara %v, %v\n", getData, l)
			nowLen := uint64(len(outBuffer))
			for i := uint64(0); i < length; i++ {
				outBuffer = append(outBuffer, outBuffer[nowLen-uint64(getData)+i])
			}
		} else if getData == huffman.HUFFMANEndFlag {
			//fmt.Printf("end buffer \n")
			break
		} else {
			outBuffer = append(outBuffer, byte(getData))
		}
	}
	return outBuffer
}
