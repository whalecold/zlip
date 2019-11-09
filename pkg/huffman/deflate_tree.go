package huffman

import (
	"fmt"
	"sort"

	"github.com/whalecold/zlip/pkg/utils"
)

//DelateBitsStreamInfo nil
type DelateBitsStreamInfo struct {
	Value   uint32 //区间值
	BitsLen byte   //树的深度 也可以称为树的信息
}

//DelateBitsArray nil
type DelateBitsArray []*DelateBitsStreamInfo

//DeflateCommon nil
//通用接口
type DeflateCommon interface {
	GetZoneData(uint16, bool) (uint16, uint16, uint16)
	//bool 表示是否是长度
	GetSourceCode(uint16) (uint16, uint16, bool)
	GetBitsLen() int
}

//DeflateTree deflate tree
type DeflateTree struct {
	//m map[uint16]*Node 	//这个是距离出现次数的映射表
	elementSlice []*Node  //优化压缩耗时
	huffmanSlice [][]byte //这个是区间码到huffman字节码的映射表
	node         *Node    //deflate树的根节点
	bits         []byte   //deflate树转成bits流的长度 用来存文件
	condition    DeflateCommon
}

//BitesLen bit len
func (deflate *DeflateTree) BitesLen() uint32 {
	return uint32(len(deflate.bits))
}

//GetBits get bits
func (deflate *DeflateTree) GetBits() []byte {
	return deflate.bits
}

//Init  init
func (deflate *DeflateTree) Init() {
	//deflat.m = make(map[uint16]*Node)
	deflate.elementSlice = make([]*Node, deflate.condition.GetBitsLen())
	deflate.huffmanSlice = make([][]byte, deflate.condition.GetBitsLen())
}

//AddElement add element
func (deflate *DeflateTree) AddElement(element uint16, length bool) {

	//zone, _, _ := getZoneByData(element, deflate.extraCode)
	zone, _, _ := deflate.condition.GetZoneData(element, length)
	//fmt.Printf("zone-------- %v\n", zone)
	ele := deflate.elementSlice[zone]
	if ele == nil {
		deflate.elementSlice[zone] = &Node{
			Power:     1,
			Value:     zone,
			LeftTree:  nil,
			RightTree: nil,
			Leaf:      true,
		}
	} else {
		ele.Power++
	}
	//if m, ok := deflate.m[zone]; ok {
	//	m.Power++
	//} else {
	//	deflate.m[zone] = &Node{
	//		Power: 1,
	//		Value: zone,
	//		LeftTree: nil,
	//		RightTree: nil,
	//		Leaf: true,
	//	}
	//}
}

//BuildTree build tree
//建立deflate树
func (deflate *DeflateTree) BuildTree() {
	huffmanSlice := make(NodeSlice, 0, len(deflate.elementSlice))
	//for _, v := range deflate.m {
	//	huffmanSlice = append(huffmanSlice, v)
	//}
	for _, v := range deflate.elementSlice {
		if v != nil {
			huffmanSlice = append(huffmanSlice, v)
		}
	}
	sort.Sort(huffmanSlice)
	deflate.node = buildTree(huffmanSlice)
}

//SerializeBitsStream serial bits
//根据字码映射表获取字节流 相当于是序列化
func (deflate *DeflateTree) SerializeBitsStream() {
	//max := 0
	deflate.bits = make([]byte, deflate.condition.GetBitsLen())
	for k, v := range deflate.huffmanSlice {
		if int(k) >= deflate.condition.GetBitsLen() {
			panic(fmt.Sprintf("BuildBitsStream error should be less than %v", deflate.condition.GetBitsLen()))
		}
		//fmt.Printf("+++++++ %v  %v", k, v)
		deflate.bits[k] = byte(len(v))
		//if len(v) > max {
		//	max = len(v)
		//}
	}
	//fmt.Printf("max-------%v\n", max)
}

//BuildMap build map
//根据deflate获取 字码映射表
func (deflate *DeflateTree) BuildMap() {
	deflate.huffmanSlice = deflate.node.transTreeToDeflateCodeMap(deflate.condition.GetBitsLen())
}

//EnCodeElement encode element
//returs offset 表示bytes[0]的位偏移 return  | bit偏移 | byte 偏移
func (deflate *DeflateTree) EnCodeElement(ele uint16,
	bytes *[]byte,
	offset uint32,
	dataSet *uint64,
	length bool) uint32 {
	if offset > 7 {
		panic("EnCodeDistance error param offset")
	}
	zone, bitLen, lower := deflate.condition.GetZoneData(ele, length)
	zoneBits := deflate.huffmanSlice[zone]
	//if !ok {
	//	deflate.Print()
	//	panic(fmt.Sprintf("EnCodeElement error para %v", zone))
	//}
	for _, value := range zoneBits {
		utils.SetHighBit8(&(*bytes)[*dataSet], offset, value)
		offset++
		if checkBytesFull(bytes, &offset) {
			*dataSet++
		}
	}

	if bitLen != 0 {
		sur := ele % lower
		for i := 16 - bitLen; i < 16; i++ {
			v := utils.GetHighBit16(sur, uint32(i))
			utils.SetHighBit8(&(*bytes)[*dataSet], offset, v)
			offset++
			if checkBytesFull(bytes, &offset) {
				*dataSet++
			}
		}
	}
	return offset
}

// DecodeEle encode ele
// 传入参数 bytes bits流 offset第一个字节之后的偏移位置
// return 第1个返回实际距离 第2个参数表示返回字节偏移  第二个参数表示bits偏移
// return 匹配到的区间码  | bytes偏移位数 | bit偏移位数(范围0-7) | bool表示是否是长度 true 是
func (deflate *DeflateTree) DecodeEle(bytes []byte,
	offset uint32) (uint16, uint32, uint32, bool) {
	code, off, bits := deflate.node.decodeCodeDeflate(bytes, offset)
	//bitsLen, lower := getDataByZone(code, deflate.extraCode)
	bitsLen, lower, flag := deflate.condition.GetSourceCode(code)

	dis, o, bits := utils.ReadBitsLen(bytes[off:], bits, bitsLen)
	//fmt.Printf("dis %v  lower %v\n", dis, lower)
	dis += lower
	return dis, off + o, bits, flag
}

// UnSerializeBitsStream unserialize
// 根据位数来重新获得码表映射
func (deflate *DeflateTree) UnSerializeBitsStream(bits []byte) {
	if len(bits) != deflate.condition.GetBitsLen() {
		panic(fmt.Sprintf("BuildTreeByBits error length %v shoud be %v", len(bits), deflate.condition.GetBitsLen()))
	}
	deflate.huffmanSlice = buildCodeMapByBits(bits)
}

// BuildTreeByMap build tree
// 根据码表map建立deflate树
func (deflate *DeflateTree) BuildTreeByMap() {
	deflate.node = buildDeflatTreeByMap(deflate.huffmanSlice)
}

//Print print
func (deflate *DeflateTree) Print() {
	//fmt.Printf("start ")
	for k, v := range deflate.huffmanSlice {
		fmt.Printf("%v -- %b\n", k, v)
	}
}

//Equal equal
func (deflate *DeflateTree) Equal(other *DeflateTree) bool {
	if deflate.huffmanSlice == nil || other.huffmanSlice == nil {
		return false
	}

	for key, bytes := range deflate.huffmanSlice {
		obytes := other.huffmanSlice[key]
		//if !ok {
		//	return false
		//}
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
