package huffman

import (
	"sort"
)

//因为码树节点左右分支旋转不会影响压缩程度 所有huffman树有很多表示
//这里采用最不平衡的树 称之为 Deflate树 在构建的时候效率也会更高 不需要递归
func buildTree(huffmanSlice HuffmanNodeSlice) *HuffmanNode {
	//fmt.Printf("tree len : %v\n", len(huffmanSlice))
	if len(huffmanSlice) == 0 {
		return nil
	}

	index := byte(0)
	for len(huffmanSlice) != 1 {
		newNode := &HuffmanNode{
			Power: huffmanSlice[0].Power + huffmanSlice[1].Power,
			Value: index,
			LeftTree: huffmanSlice[0],
			RightTree: huffmanSlice[1],
			Leaf: false,
		}
		index++
		huffmanSlice[0] = newNode
		huffmanSlice = append(huffmanSlice[:1], huffmanSlice[2:]...)
		sort.Sort(huffmanSlice)
	}


	return huffmanSlice[0]
	//if len(huffmanSlice) < 2 {
	//	panic("buildTree need params longer than 2")
	//}
	//
	//root := huffmanSlice[0]
	//
	//for i := 1; i < len(huffmanSlice); i++ {
	//	node := &HuffmanNode{
	//		Power: huffmanSlice[0].Power + huffmanSlice[1].Power,
	//		Value: '0',
	//		LeftTree: huffmanSlice[0],
	//		RightTree: huffmanSlice[1],
	//		Leaf: false,
	//	}
	//	node.LeftTree = huffmanSlice[i]
	//	node.RightTree = root
	//	root = node
	//}
	//return root
}

//构建huffmans树
func buildHuffmanTree(bytes []byte) *HuffmanNode {

	huffmanMap := make(map[byte]*HuffmanNode)
	for _, value := range bytes {
		if m, ok := huffmanMap[value]; ok {
			m.Power++
		} else {
			huffmanMap[value] = &HuffmanNode{
				Power: 1,
				Value: value,
				LeftTree: nil,
				RightTree: nil,
				Leaf: true,
			}
		}
	}

	huffmanSlice := make(HuffmanNodeSlice, 0, len(huffmanMap))
	for _, v := range huffmanMap {
		huffmanSlice = append(huffmanSlice, v)
	}
	sort.Sort(huffmanSlice)

	return buildTree(huffmanSlice)
}
