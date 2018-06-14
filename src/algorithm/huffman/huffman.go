package huffman

import (
	"sort"

	"container/list"
	"fmt"
)

//因为码树节点左右分支旋转不会影响压缩程度 所有huffman树有很多表示
//这里采用最不平衡的树 称之为 Deflate树 在构建的时候效率也会更高 不需要递归
func buildTree(huffmanSlice HuffmanNodeSlice) *HuffmanNode {
	//fmt.Printf("tree len : %v\n", len(huffmanSlice))
	if len(huffmanSlice) == 0 {
		panic("buildTree param length is zero!")
	}

	index := uint16(0)
	for len(huffmanSlice) != 1 {
		newNode := &HuffmanNode{
			Power: huffmanSlice[0].Power + huffmanSlice[1].Power,
			Value: index,
			LeftTree: huffmanSlice[0],
			RightTree: huffmanSlice[1],
			Leaf: false,
		}
		//fmt.Printf("huffmanSlice[0] %v   %v\n",
		//	huffmanSlice[0].Value, huffmanSlice[1].Value)
		index++
		huffmanSlice[0] = newNode
		huffmanSlice = append(huffmanSlice[:1], huffmanSlice[2:]...)
		sort.Sort(huffmanSlice)
	}


	return transDeflateTree(huffmanSlice[0], moveDeflateTree)
}


//把huffman树转成deflate树
func transDeflateTree(root *HuffmanNode, f func(ele *list.Element, deepth int)) *HuffmanNode {

	var nextLen int
	var thisLen int
	var deepth int

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
			f(ele, deepth)
			li = list.New()
			deepth ++

		}
	}
	return root
}

//把某一层的树移到最右边 这个步骤是构造deflate树的必要条件
//1、 同一层的子节点右边一定要比左边大
//2、 右边的树一定要比左边的深
func moveDeflateTree(ele *list.Element, deepth int) {
	var emptyEle *list.Element
	var leafNodeNum int
	emptyList := list.New() 	//存储空节点的队列
	for temp := ele; temp != nil; temp = temp.Prev() {
		tempNode := temp.Value.(*HuffmanNode)

		//if deepth == 3 {
		//	fmt.Printf("tempNode %v\n", tempNode.Value)
		//}

		if tempNode.Leaf == true && tempNode.LeftTree == nil {
			leafNodeNum++
		}
		if tempNode.RightTree == nil {
			emptyList.PushBack(tempNode)
		}

		if tempNode.RightTree != nil && emptyList.Front() != nil {
			emptyEle = emptyList.Front()
			emptyNode := emptyEle.Value.(*HuffmanNode)
			emptyNode.RightTree, tempNode.RightTree = tempNode.RightTree, nil
			emptyNode.LeftTree, tempNode.LeftTree = tempNode.LeftTree, nil
			emptyNode.Value, tempNode.Value = tempNode.Value, emptyNode.Value
			emptyNode.Leaf, tempNode.Leaf = tempNode.Leaf, emptyNode.Leaf
			emptyList.Remove(emptyEle)
			emptyList.PushBack(tempNode)
		}
	}



	nodeSort := make([]*HuffmanNode, 0, leafNodeNum)
	for temp := ele; temp != nil; temp = temp.Prev() {
		sortNode := temp.Value.(*HuffmanNode)
		//这里是从后往前遍历的 把自己坑了啊
		if sortNode.Leaf == true && sortNode.LeftTree == nil {
			//nodeSort = append(nodeSort, sortNode)
			nodeSort = append([]*HuffmanNode{sortNode}, nodeSort...)
		}
	}


	for index := 0; index < len(nodeSort); index++ {
		var maxNum uint16
		var maxIndex int
		var i int
		for ; i < len(nodeSort) - index; i++ {
			if nodeSort[i].Value > maxNum {
				maxNum = nodeSort[i].Value
				maxIndex = i
			}
		}
		nodeSort[i-1].Value, nodeSort[maxIndex].Value =
			nodeSort[maxIndex].Value, nodeSort[i-1].Value
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
				Value: uint16(value),
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

//根据位数来重新获得码表
func buildCodeMapByBits(bits  []byte) DeflateCodeMap {
	if len(bits) != len(distanceZone) {
		panic(fmt.Sprintf("BuildTreeByBits error length %v", len(bits)))
	}
	m := make(DeflateCodeMap)
	deepth := getMaxDeepth(bits)
	//fmt.Printf("cap ------%v \n", deepth)
	streamTemp := make([][]uint16, deepth + 1)
	for index := 0; index < len(streamTemp); index++ {
		streamTemp[index] = make([]uint16, 0, 32)
	}

	for index, value := range bits {
		//fmt.Printf("value --------- %v\n", value)
		//fmt.Printf("cap %v \n", cap(streamTemp[value]))
		if value != 0 {
			streamTemp[value] = append(streamTemp[value], uint16(index))
		}
	}

	//表示是否出现了长度不为0的值 这里的index表示长度 这里是根据deflate的规律来的
	flag := false
	var lastCode, lastLength int
	for index := 1; index < len(streamTemp); index++ {
		if len(streamTemp[index]) == 0 && flag == false{
			continue
		}
		//deflate树的最左边的节点始终为0
		lastCode = (lastCode + lastLength) << 1
		tempCode := lastCode

		for i := 0; i < len(streamTemp[index]); i++ {
			bytes := make([]byte, index)
			for t := 0; t < index; t++ {
				bytes[index-1-t] = ReadBit(tempCode, uint(t))
			}
			m[streamTemp[index][i]] = bytes
			tempCode++
		}
		lastLength = len(streamTemp[index])
		flag = true
	}
	return m
}


