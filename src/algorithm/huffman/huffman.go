package huffman

import (
	"sort"

	"container/list"
)

//因为码树节点左右分支旋转不会影响压缩程度 所有huffman树有很多表示
//这里采用最不平衡的树 称之为 Deflate树 在构建的时候效率也会更高 不需要递归
func buildTree(huffmanSlice HuffmanNodeSlice) *HuffmanNode {
	//fmt.Printf("tree len : %v\n", len(huffmanSlice))
	if len(huffmanSlice) == 0 {
		panic("buildTree param length is zero!")
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

		//for _, v := range huffmanSlice {
		//	fmt.Printf("%p \t", *v)
		//}
		//fmt.Printf("\n")
	}


	return transDeflateTree(huffmanSlice[0])
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


//把huffman树转成deflate树
func transDeflateTree(root *HuffmanNode) *HuffmanNode {

	var nextLen int
	var thisLen int

	l := list.New()
	l.PushBack(root)
	li := list.New()
	thisLen = 1
	for {
		element := l.Front()
		if element == nil {
			break
		}

		node := element.Value.(*HuffmanNode)
		li.PushBack(node)
		//fmt.Printf("node : %v \n", node.Value)

		l.Remove(element)

		if node.LeftTree != nil {
			l.PushBack(node.LeftTree)
			nextLen++
		}

		if node.RightTree != nil {
			l.PushBack(node.RightTree)
			nextLen++
		}
		thisLen--

		if 0 == thisLen {
			thisLen = nextLen
			nextLen = 0

			ele := li.Back()
			moveDeflateTree(ele)
			li = list.New()
		}
	}
	return root
}

//把某一层的树移到最右边
func moveDeflateTree(ele *list.Element) {
	var emptyEle *list.Element
	for temp := ele; temp != nil; temp = temp.Prev() {
		tempNode := temp.Value.(*HuffmanNode)
		if tempNode.RightTree == nil && emptyEle == nil {
			emptyEle = temp
		} else if tempNode.RightTree != nil && emptyEle != nil {
			emptyNode := emptyEle.Value.(*HuffmanNode)
			emptyNode.RightTree, tempNode.RightTree = tempNode.RightTree, nil
			emptyNode.LeftTree, tempNode.LeftTree = tempNode.LeftTree, nil
			emptyNode.Value, tempNode.Value = tempNode.Value, emptyNode.Value
			emptyNode.Leaf, tempNode.Leaf = tempNode.Leaf, emptyNode.Leaf
			emptyEle = temp
		}
	}
}

//构建huffmans树
func buildHuffmanTree(bytes []byte) *HuffmanNode {
	//fmt.Printf("bytes %v\n", bytes)
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

	//for _, v := range huffmanSlice {
	//	fmt.Printf("%p \t", *v)
	//}
	//fmt.Printf("\n")

	return buildTree(huffmanSlice)
}


