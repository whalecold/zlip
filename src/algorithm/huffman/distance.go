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

//returs bitstream bitslen uint64只是个放bits的流载体
func (deflate *DeflateDisTree)EnCodeDistance(dis uint16) (uint64, int){
	zone, bitLen, lower := GetZoneByDis(dis)
	zoneBits, ok  := deflate.dishuffMap[zone]
	if !ok {
		panic("EnCodeDistance error param")
	}
	var bits uint64
	var length int
	//var offset uint
	for _, value := range zoneBits {
		bits = bits << 1
		bits ^= uint64(value)
		length++
	}
	sur := dis % lower
	bits = bits << bitLen
	bits += uint64(sur)
	return bits, length + int(bitLen)
}

//根据位数来重新获得码表
func (deflate *DeflateDisTree)BuildCodeMapByBits(bits  []byte) {
	deflate.dishuffMap = buildCodeMapByBits(bits)
}

func (deflate *DeflateDisTree)BuildTreeByMap() {
	//deflate.node = buildDeflatTreeByMap(deflate.dishuffMap)
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
