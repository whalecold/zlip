package huffman

import (
	"github.com/whalecold/compress/pkg/stack"
	"github.com/whalecold/compress/pkg/utils"
)

// Node huffmannode
type Node struct {
	Power     int32 //权重 叶子节点相当于出现次数
	Value     uint16
	LeftTree  *Node
	RightTree *Node
	Leaf      bool //表示是否是叶子节点
}

//NodeSlice slice
type NodeSlice []*Node

func (h NodeSlice) Less(i, j int) bool {
	if h[i].Power != h[j].Power {
		return h[i].Power < h[j].Power
	}
	return h[i].Value < h[j].Value
}

func (h NodeSlice) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h NodeSlice) Len() int {
	return len(h)
}

//CodeMap map
type CodeMap map[byte][]byte

//DeflateCodeMap deflatecodemap
type DeflateCodeMap map[uint16][]byte

//return 匹配到的byte | 移动的bit位数 | bytes偏移位数 | bit偏移位数(范围0-7)
func (huff *Node) decodeByteFromHuffman(bytes []byte, bitOffset uint32) (byte, uint32, uint32, uint32) {

	tempNode := huff
	var bitLen uint32
	var byteLen uint32
	for _, value := range bytes {
		//fmt.Printf("bitOffset %v\n", bitOffset)
		for ; bitOffset < 8; bitOffset++ {
			bitLen++
			//fmt.Printf("%b\t", value)
			bit := utils.ReadBitsHigh(value, bitOffset)
			//fmt.Printf("------------------------------------%b\n", bit)

			if bit == 0 {
				tempNode = tempNode.LeftTree
			} else if bit == 1 {
				tempNode = tempNode.RightTree
			} else {
				panic("decodeByteFromHuffman error bit")
			}
			if tempNode.Leaf == true {
				return byte(tempNode.Value), bitLen, byteLen, bitOffset + 1
			}
		}
		//备注 :这里bitOffset 需要清0 之前落了导致bug
		bitOffset = 0
		byteLen++
	}
	//走到这个肯定是程序出错了 找不到对应字符串是不能发生的
	panic("decodeByteFromHuffman failed !")
}

//上面那个是之前测试用的
//return 匹配到的区间码  | bytes偏移位数 | bit偏移位数(范围0-7)
func (huff *Node) decodeCodeDeflate(bytes []byte, bitOffset uint32) (uint16, uint32, uint32) {

	tempNode := huff
	var byteLen uint32
	for _, value := range bytes {
		for ; bitOffset < 8; bitOffset++ {
			bit := utils.ReadBitsHigh(value, bitOffset)

			if bit == 0 {
				tempNode = tempNode.LeftTree
			} else if bit == 1 {
				tempNode = tempNode.RightTree
			} else {
				panic("decodeByteFromHuffman error bit")
			}
			if tempNode.Leaf == true {
				return tempNode.Value, byteLen, bitOffset + 1
			}
		}
		//备注 :这里bitOffset 需要清0 之前落了导致bug
		bitOffset = 0
		byteLen++
	}
	//走到这个肯定是程序出错了 找不到对应字符串是不能发生的
	panic("decodeByteFromHuffman failed !")
}

//在树的节点没有重复的情况下 树的前序遍历数组和中序遍历数组能建立唯一的树
//所以这里产生两个数组 用来以后建立树
func (huff *Node) genStreamByPreorder() []byte {
	//每个叶子节点都需要有值 而且必须每个都不一样
	// uint16高八位 1表示非叶子节点 第八位表示序号  高八位0表示叶子节点 低八位表示实际序号
	preorderSlice := make([]byte, 0, 512)
	stackNode := stack.NewStack()
	stackNode.Push(huff)
	for stackNode.Len() != 0 {
		node := stackNode.RPop().(*Node)
		if node.Leaf == true {
			preorderSlice = append(preorderSlice, 0)
		} else {
			preorderSlice = append(preorderSlice, 1)
		}
		preorderSlice = append(preorderSlice, byte(node.Value))

		if node.RightTree != nil {
			stackNode.Push(node.RightTree)
		}

		if node.LeftTree != nil {
			stackNode.Push(node.LeftTree)
		}
	}
	//fmt.Printf("--- %v\n", preorderSlice)
	return preorderSlice
}

//获取中序遍历数据
func (huff *Node) genStreamByInorder() []byte {
	inorderSlice := make([]byte, 0, 512)

	s := stack.NewStack()
	node := huff

	for node != nil || s.Len() != 0 {
		for node != nil {
			s.Push(node)
			node = node.LeftTree
		}

		if s.Len() != 0 {
			node = s.RPop().(*Node)

			if node.Leaf == true {
				inorderSlice = append(inorderSlice, 0)
			} else {
				inorderSlice = append(inorderSlice, 1)
			}
			inorderSlice = append(inorderSlice, byte(node.Value))

			node = node.RightTree
		}
	}
	return inorderSlice
}

//待优化这种序列化方式占用的空间有点高 这只是自己想的存储码表的方式 实际不采用这种方式
func (huff *Node) serializeTree() []byte {
	pre := huff.genStreamByPreorder()
	in := huff.genStreamByInorder()
	serialize := make([]byte, 0, len(pre)+len(in))
	serialize = append(serialize, pre...)
	serialize = append(serialize, in...)
	//fmt.Printf("pre : %v  in %v\n", pre, in)
	return serialize
}

//根据上面获得的两个数组来建立一个数
func buildTreeBySlice(pre, in []byte) *Node {
	preShort := transUint16Byte(pre)
	inShort := transUint16Byte(in)

	return buildTreeByOrder(preShort, inShort)
}

//反序列化
func buildTreeBySerialize(serial []byte, size uint32) *Node {
	preShort := transUint16Byte(serial[:size/2])
	inShort := transUint16Byte(serial[size/2:])

	return buildTreeByOrder(preShort, inShort)
}

//根据前序遍历和中序遍历建立一个新的树
func buildTreeByOrder(pre, in []uint16) *Node {
	if 0 == len(pre) || 0 == len(in) {
		return nil
	}

	midNumber := pre[0]
	midIndex := 0

	root := &Node{
		Value: uint16(midNumber & 0xFF),
	}

	//这里的1表示是否是叶子节点
	//fmt.Printf("mid hight %v\n", midNumber)
	if (midNumber & uint16(0xff00)) == 0 {
		root.Leaf = true
	} else {
		root.Leaf = false
	}

	for i := 0; i < len(in); i++ {
		if midNumber == in[i] {
			midIndex = i
			break
		}
	}
	//fmt.Printf("value : %v leaf %v mid %v  minNumber %v\n", root.Value, root.Leaf, midIndex, midNumber)

	if midIndex == len(in) {
		return root
	}

	leftChild := buildTreeByOrder(pre[1:midIndex+1], in[:midIndex])
	rightChild := buildTreeByOrder(pre[midIndex+1:], in[midIndex+1:])

	root.LeftTree = leftChild
	root.RightTree = rightChild

	return root
}

//因为是两个字节表示一个数
func transUint16Byte(bytes []byte) []uint16 {
	if len(bytes)%2 == 1 {
		panic("buildTreeBySlice param error!")
	}

	uslice := make([]uint16, 0, len(bytes)/2)

	for i := 0; i < len(bytes); i += 2 {
		value := uint16(bytes[i])<<8 + uint16(bytes[i+1])
		//fmt.Printf("high %b low %b value : %v\n", bytes[i], bytes[1+i], value)
		uslice = append(uslice, value)
	}
	return uslice
}

//调用这个函数会破坏树的结构 后序遍历
//https://blog.csdn.net/gatieme/article/details/51163010
func (huff *Node) transTreeToHuffmanCodeMap() CodeMap {
	if huff == nil {
		return nil
	}
	m := make(CodeMap)
	s := stack.NewStack()
	s.Push(huff)

	huffmanCode := make([]byte, 0, 64)
	for s.Len() != 0 {
		tree := s.Pop().(*Node)
		if tree.LeftTree != nil {
			s.Push(tree.LeftTree)
			tree.LeftTree = nil
			huffmanCode = append(huffmanCode, 0)
		} else if tree.RightTree != nil {
			s.Push(tree.RightTree)
			tree.RightTree = nil
			huffmanCode = append(huffmanCode, 1)
		} else {
			s.RPop()
			if tree.Leaf == true {
				m[byte(tree.Value)] = *utils.DeepClone(&huffmanCode).(*[]byte)
			}
			if len(huffmanCode) > 0 {
				huffmanCode = huffmanCode[:len(huffmanCode)-1]
			}
		}
	}
	return m
}

//调用这个函数会破坏树的结构
func (huff *Node) transTreeToDeflateCodeMap(length int) [][]byte {
	if huff == nil {
		return nil
	}
	result := make([][]byte, length)
	s := stack.NewStack()
	s.Push(huff)

	huffmanCode := make([]byte, 0, 64)
	for s.Len() != 0 {
		tree := s.Pop().(*Node)
		if tree.LeftTree != nil {
			s.Push(tree.LeftTree)
			tree.LeftTree = nil
			huffmanCode = append(huffmanCode, 0)
		} else if tree.RightTree != nil {
			s.Push(tree.RightTree)
			tree.RightTree = nil
			huffmanCode = append(huffmanCode, 1)
		} else {
			s.RPop()
			if tree.Leaf == true {
				result[tree.Value] = *utils.DeepClone(&huffmanCode).(*[]byte)
			}
			if len(huffmanCode) > 0 {
				huffmanCode = huffmanCode[:len(huffmanCode)-1]
			}
		}
	}
	return result
}

// Printf printf
func (huff *Node) Printf() {

}
