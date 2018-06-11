package huffman

import (
	"algorithm/stack"
	"utils"
	//"fmt"
)

type HuffmanNode struct {
	Power 	int32 	//权重 叶子节点相当于出现次数
	Value 	byte
	LeftTree *HuffmanNode
	RightTree *HuffmanNode
	Leaf	bool 		//表示是否是叶子节点
}


type HuffmanCodeMap map[byte][]byte

//最高位的offset表示0
func (huff *HuffmanNode)readBites(b byte, offset uint32) byte {
	move := 7 - offset
	b = b >> move
	b &= 0x1
	return b
}

//return 匹配到的byte  移动的bit位数  bytes偏移位数 bit偏移位数(范围0-7)
func (huff *HuffmanNode)decodeByteFromHuffman(bytes []byte, bitOffset uint32) (byte, uint32, uint32, uint32) {

	//fmt.Printf("++++%b+++  bitOffset %v\n", bytes, bitOffset)
	tempNode := huff
	var bitLen uint32
	var byteLen uint32
	for _, value := range bytes {
		//fmt.Printf("bitOffset %v\n", bitOffset)
		for i := bitOffset; i < 8; i++ {
			bitLen++
			//fmt.Printf("%b\t", value)
			bit := huff.readBites(value, i)
			//fmt.Printf("------------------------------------%b\n", bit)

			if bit == 0 {
				tempNode = tempNode.LeftTree
			} else if bit == 1 {
				tempNode = tempNode.RightTree
			}
			if tempNode.Leaf == true {
				return tempNode.Value, bitLen, byteLen, i+1
			}
		}
		bitOffset = 0
		byteLen++
	}
	panic("decodeByteFromHuffman failed !")
}

//在树的节点没有重复的情况下 树的前序遍历数组和中序遍历数组能建立唯一的树
//所以这里产生两个数组 用来以后建立树
func (huff *HuffmanNode)genStreamByPreorder() []byte {
	//每个叶子节点都需要有值 而且必须每个都不一样
	// uint16高八位 1表示非叶子节点 第八位表示序号  高八位0表示叶子节点 低八位表示实际序号
	preorderSlice := make([]byte, 0, 512)
	stack_node := stack.NewStack()
	stack_node.Push(huff)
	for stack_node.Len() != 0 {
		node := stack_node.RPop().(*HuffmanNode)
		if node.Leaf == true {
			preorderSlice = append(preorderSlice, 0)
		} else {
			preorderSlice = append(preorderSlice, 1)
		}
		preorderSlice = append(preorderSlice, node.Value)

		if node.RightTree != nil {
			stack_node.Push(node.RightTree)
		}

		if node.LeftTree != nil {
			stack_node.Push(node.LeftTree)
		}
	}
	return preorderSlice
}

//获取中序遍历数据
func (huff *HuffmanNode)genStreamByInorder() []byte {
	inorderSlice := make([]byte, 0, 512)

	s := stack.NewStack()
	node := huff

	for node != nil || s.Len() != 0 {
		for node != nil {
			s.Push(node)
			node = node.LeftTree
		}

		if s.Len() != 0 {
			node = s.RPop().(*HuffmanNode)

			if node.Leaf == true {
				inorderSlice = append(inorderSlice, 0)
			} else {
				inorderSlice = append(inorderSlice, 1)
			}
			inorderSlice = append(inorderSlice, node.Value)

			node = node.RightTree
		}
	}
	return inorderSlice
}

//根据上面获得的两个数组来建立一个数
func buildTreeBySlice(pre, in []byte) *HuffmanNode {
	preShort := transUint16Byte(pre)
	inShort := transUint16Byte(in)

	return buildTreeByOrder(preShort, inShort)
}

//根据前序遍历和中序遍历建立一个新的树
func buildTreeByOrder(pre, in []uint16) *HuffmanNode {
	if 0 == len(pre) || 0 == len(in) {
		return nil
	}

	midNumber := pre[0]
	midIndex := 0

	root := &HuffmanNode{
		Value: byte(midNumber & 0xFF),
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

func transUint16Byte(bytes []byte) []uint16 {
	if len(bytes) % 2 == 1 {
		panic("buildTreeBySlice param error!")
	}

	uslice := make([]uint16, 0, len(bytes) / 2)

	for i:=0; i < len(bytes); i+=2 {
		value := uint16(bytes[i]) << 8 + uint16(bytes[i+1])
		//fmt.Printf("high %b low %b value : %v\n", bytes[i], bytes[1+i], value)
		uslice = append(uslice, value)
	}
	return uslice
}


//https://blog.csdn.net/gatieme/article/details/51163010
func (huff *HuffmanNode)transTreeToHuffmanCodeMap() HuffmanCodeMap {
	if huff == nil {
		return nil
	}
	m := make(HuffmanCodeMap)
	s := stack.NewStack()
	s.Push(huff)


	//var huffmanCode uint32
	//var huffmanCodeSkip uint32

	huffmanCode := make([]byte, 0, 64)
	for s.Len() != 0 {
		//fmt.Printf("len : %v\n", s.Len())
		tree := s.Pop().(*HuffmanNode)
		if tree.LeftTree != nil {
			s.Push(tree.LeftTree)
			tree.LeftTree = nil
			//huffmanCode = huffmanCode << 1
			//huffmanCode &=  ^uint32(0x1)
			huffmanCode = append(huffmanCode, 0)
		} else if tree.RightTree != nil {
			s.Push(tree.RightTree)
			tree.RightTree = nil
			//huffmanCode = huffmanCode << 1
			//huffmanCode |= 0x1

			huffmanCode = append(huffmanCode, 1)
		} else {
			s.RPop()
			if tree.Leaf == true {
				m[tree.Value] = *utils.DeepClone(&huffmanCode).(*[]byte)
				//fmt.Printf("----value %v huffmanCode : %v\n", tree.Value, huffmanCode)
			}
			if len(huffmanCode) > 0 {
				huffmanCode = huffmanCode[:len(huffmanCode)-1]
			}
		}
	}
	return m
}

func (huff *HuffmanNode)Printf() {

}

type HuffmanNodeSlice []*HuffmanNode

func (h HuffmanNodeSlice)Less(i, j int) bool {
	if h[i].Power != h[j].Power {
		return h[i].Power < h[j].Power
	} else {
		return h[i].Value < h[j].Value
	}
}

func (h HuffmanNodeSlice)Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h HuffmanNodeSlice)Len() int {
	return len(h)
}