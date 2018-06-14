package huffman

import (
	"sort"
	"fmt"
)

type DelateBitsStreamInfo struct {
	Value uint32	//区间值
	BitsLen byte 	//树的深度 也可以称为树的信息
}

type DelateBitsArray []*DelateBitsStreamInfo


type DeflateDisTree struct {
	m map[uint16]*HuffmanNode 	//这个是距离出现次数的映射表
	dishuffMap DeflateCodeMap	//这个是区间码到huffman字节码的映射表
	node *HuffmanNode 	//deflate树的根节点
	bits []byte //deflate树转成bits流的长度 用来存文件
}

func (deflat *DeflateDisTree)Init() {
	deflat.m = make(map[uint16]*HuffmanNode)
	deflat.dishuffMap = make(map[uint16][]byte)
}

func (deflate *DeflateDisTree)AddDisElement(distance uint16) {

	zone, _, _ := GetZoneByDis(distance)
	if m, ok := deflate.m[zone]; ok {
		m.Power++
	} else {
		deflate.m[zone] = &HuffmanNode{
			Power: 1,
			Value: zone,
			LeftTree: nil,
			RightTree: nil,
			Leaf: true,
		}
	}
}

func (deflate *DeflateDisTree)BuildTree() {
	huffmanSlice := make(HuffmanNodeSlice, 0, len(deflate.m))
	for _, v := range deflate.m {
		huffmanSlice = append(huffmanSlice, v)
	}
	sort.Sort(huffmanSlice)
	deflate.node = buildTree(huffmanSlice)
}

func (deflate *DeflateDisTree)BuildBitsStream() {
	deflate.bits = make([]byte, len(distanceZone))
	for k, v := range deflate.dishuffMap {
		if int(k) >= len(deflate.bits) {
			panic("BuildBitsStream error should be less than 30")
		}
		deflate.bits[k] = byte(len(v))
	}
}

func (deflate *DeflateDisTree)BuildMap() {
	deflate.dishuffMap = deflate.node.transTreeToDeflateCodeMap()
}

//returs offset 表示bytes[0]的位偏移 return  | bit偏移 | byte 偏移
func (deflate *DeflateDisTree)EnCodeDistance(dis uint16, bytes *[]byte, offset uint32, dataSet *uint64) (uint32){
	if offset > 7 {
		panic("EnCodeDistance error param offset")
	}
	zone, bitLen, lower := GetZoneByDis(dis)
	zoneBits, ok  := deflate.dishuffMap[zone]
	if !ok {
		panic("EnCodeDistance error param")
	}
	for _, value := range zoneBits {
		WriteBitsHigh(&(*bytes)[*dataSet], offset, value)
		offset++
		if checkBytesFull(bytes, &offset) == true {
			*dataSet ++
		}
	}

	sur := dis % lower
	//fmt.Printf("sur %v lower %v\n", sur, lower)
	for i := 16-bitLen; i < 16; i++ {
		v := readBitsHigh16(sur, uint32(i))
		WriteBitsHigh(&(*bytes)[*dataSet], offset, v)
		offset++
		if checkBytesFull(bytes, &offset) == true {
			*dataSet++
		}
	}
	return  offset
}

//传入参数 bytes bits流 offset第一个字节之后的偏移位置
//return 第1个返回实际距离 第2个参数表示返回字节偏移  第二个参数表示bits偏移
//return 匹配到的区间码  | bytes偏移位数 | bit偏移位数(范围0-7)
func (deflate *DeflateDisTree)DecodeDistance(bytes []byte, offset uint32) (uint16, uint32, uint32) {
	code, off, bits := deflate.node.decodeCodeDeflate(bytes, offset)
	bitsLen, lower := GetDisByData(code)
	dis, o, bits := readBitsLen(bytes[off:], bits, bitsLen)
	//fmt.Printf("dis %v  lower %v\n", dis, lower)
	dis += lower
	return dis, off + o, bits
}

//根据位数来重新获得码表
func (deflate *DeflateDisTree)BuildCodeMapByBits(bits  []byte) {
	deflate.dishuffMap = buildCodeMapByBits(bits)
}

func (deflate *DeflateDisTree)BuildTreeByMap() {
	deflate.node = buildDeflatTreeByMap(deflate.dishuffMap)
}

func (deflate *DeflateDisTree)Print() {
	for k, v := range deflate.dishuffMap {
		fmt.Printf("%v -- %b\n", k, v)
	}
}

func (deflate *DeflateDisTree)Equal(other *DeflateDisTree) bool {
	if deflate.dishuffMap == nil || other.dishuffMap == nil {
		return false
	}

	for key, bytes := range deflate.dishuffMap {
		obytes, ok := other.dishuffMap[key]
		if !ok {
			return false
		}
		if len(bytes) != len(obytes) {
			return false
		}

		for index, b := range bytes {
			if b != obytes[index] {
				return false
			}
		}
	}
	return true
}
