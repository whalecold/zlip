package huffman

import (
	"sort"
	"fmt"
	"utils"
)

type DelateBitsStreamInfo struct {
	Value uint32	//区间值
	BitsLen byte 	//树的深度 也可以称为树的信息
}

type DelateBitsArray []*DelateBitsStreamInfo


//通用接口
type DeflateCommon interface{
	GetZoneData(uint16, [][]uint16, bool)  (uint16, uint16, uint16)
	//bool 表示是否是长度
	GetSourceCode(uint16, [][]uint16)  (uint16, uint16, bool)
	GetBitsLen(data [][]uint16) int
}

type DeflateTree struct {
	m map[uint16]*HuffmanNode 	//这个是距离出现次数的映射表
	dishuffMap DeflateCodeMap	//这个是区间码到huffman字节码的映射表
	node *HuffmanNode 	//deflate树的根节点
	bits []byte //deflate树转成bits流的长度 用来存文件
	extraCode [][]uint16	//码表
}

func (deflate *DeflateTree)BitesLen() uint32 {
	return uint32(len(deflate.bits))
}

func (deflate *DeflateTree)GetBits() []byte {
	return deflate.bits
}

func (deflat *DeflateTree)Init(extra [][]uint16) {
	deflat.m = make(map[uint16]*HuffmanNode)
	deflat.dishuffMap = make(map[uint16][]byte)
	deflat.extraCode = extra
}

func (deflate *DeflateTree)AddElement(element uint16, common DeflateCommon,
										length bool) {

	//zone, _, _ := getZoneByData(element, deflate.extraCode)
	zone, _, _ := common.GetZoneData(element, deflate.extraCode, length)
	//fmt.Printf("zone-------- %v\n", zone)
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
//建立deflate树
func (deflate *DeflateTree)BuildTree() {
	huffmanSlice := make(HuffmanNodeSlice, 0, len(deflate.m))
	for _, v := range deflate.m {
		huffmanSlice = append(huffmanSlice, v)
	}
	sort.Sort(huffmanSlice)
	deflate.node = buildTree(huffmanSlice)
}

//根据字码映射表获取字节流 相当于是序列化
func (deflate *DeflateTree)SerializeBitsStream(common DeflateCommon) {
	deflate.bits = make([]byte, common.GetBitsLen(deflate.extraCode))
	for k, v := range deflate.dishuffMap {
		if int(k) >= common.GetBitsLen(deflate.extraCode) {
			panic("BuildBitsStream error should be less than 30")
		}
		//fmt.Printf("+++++++ %v  %v", k, v)
		deflate.bits[k] = byte(len(v))
	}
	//fmt.Printf("--%v\n", deflate.bits)
}

//根据deflate获取 字码映射表
func (deflate *DeflateTree)BuildMap() {
	deflate.dishuffMap = deflate.node.transTreeToDeflateCodeMap()
}

//returs offset 表示bytes[0]的位偏移 return  | bit偏移 | byte 偏移
func (deflate *DeflateTree)EnCodeElement(ele uint16,
										bytes *[]byte,
										offset uint32,
										dataSet *uint64,
										common DeflateCommon,
									 	length bool ) (uint32){
	if offset > 7 {
		panic("EnCodeDistance error param offset")
	}
	zone, bitLen, lower := common.GetZoneData(ele, deflate.extraCode, length)
	zoneBits, ok  := deflate.dishuffMap[zone]
	if !ok {
		deflate.Print()
		panic(fmt.Sprintf("EnCodeElement error para %v", zone))
	}
	for _, value := range zoneBits {
		utils.WriteBitsHigh(&(*bytes)[*dataSet], offset, value)
		offset++
		if checkBytesFull(bytes, &offset) == true {
			*dataSet ++
		}
	}

	if bitLen != 0 {
		sur := ele % lower
		for i := 16-bitLen; i < 16; i++ {
			v := utils.ReadBitsHigh16(sur, uint32(i))
			utils.WriteBitsHigh(&(*bytes)[*dataSet], offset, v)
			offset++
			if checkBytesFull(bytes, &offset) == true {
				*dataSet++
			}
		}
	}
	return  offset
}

//传入参数 bytes bits流 offset第一个字节之后的偏移位置
//return 第1个返回实际距离 第2个参数表示返回字节偏移  第二个参数表示bits偏移
//return 匹配到的区间码  | bytes偏移位数 | bit偏移位数(范围0-7) | bool表示是否是长度 true 是
func (deflate *DeflateTree)DecodeEle(bytes []byte,
									offset uint32,
									common DeflateCommon) (uint16, uint32, uint32, bool) {
	code, off, bits := deflate.node.decodeCodeDeflate(bytes, offset)
	//bitsLen, lower := getDataByZone(code, deflate.extraCode)
	bitsLen, lower, flag := common.GetSourceCode(code, deflate.extraCode)

	dis, o, bits := utils.ReadBitsLen(bytes[off:], bits, bitsLen)
	//fmt.Printf("dis %v  lower %v\n", dis, lower)
	dis += lower
	return dis, off + o, bits, flag
}

//根据位数来重新获得码表映射
func (deflate *DeflateTree)UnSerializeBitsStream(bits  []byte, common DeflateCommon) {
	if len(bits) != common.GetBitsLen(deflate.extraCode) {
		panic(fmt.Sprintf("BuildTreeByBits error length %v", len(bits)))
	}
	deflate.dishuffMap = buildCodeMapByBits(bits)
}

//根据码表map建立deflate树
func (deflate *DeflateTree)BuildTreeByMap() {
	deflate.node = buildDeflatTreeByMap(deflate.dishuffMap)
}

func (deflate *DeflateTree)Print() {
	//fmt.Printf("start ")
	for k, v := range deflate.dishuffMap {
		fmt.Printf("%v -- %b\n", k, v)
	}
}

func (deflate *DeflateTree)Equal(other *DeflateTree) bool {
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
