package huffman

import (
	"sort"

	"container/list"
	"whalecold/compress/pkg/utils"
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
			Power:     huffmanSlice[0].Power + huffmanSlice[1].Power,
			Value:     index,
			LeftTree:  huffmanSlice[0],
			RightTree: huffmanSlice[1],
			Leaf:      false,
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
			deepth++

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
	emptyList := list.New() //存储空节点的队列
	for temp := ele; temp != nil; temp = temp.Prev() {
		tempNode := temp.Value.(*HuffmanNode)

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
		for ; i < len(nodeSort)-index; i++ {
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
				Power:     1,
				Value:     uint16(value),
				LeftTree:  nil,
				RightTree: nil,
				Leaf:      true,
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
//
//
//https://blog.csdn.net/jison_r_wang/article/details/52071841
/*
	deflate树:1、右边一定比左边深
			  2、同一深度的子节点 右边的值一定比左边的大
	N 表示节点
    l 表示子节点 后面的数字表示具体的值
				                  根节点
							0/                  \1
							N                  --N--
						0/     \1         0/            \1
						l:4    l:12       N              N
										0/  \1       0/        \1
										l:3  l:9     l:19       N
														   0/       \1
												          N          N
													    0/  \1    0/    \1
                                                      l:1   l:23  l:17  l:27
						原码          码表       长度
					    4		--   00			2
						12 		--   01			2
						3 		--   100		3
						9 		--   101		3
						19		--	 110		3
						1		-- 	 11100		5
						23 		-- 	 11101		5
						17 		-- 	 11110		5
						27 		-- 	 11111		5
	在序列化的时候记录码表的长度 然后用下标记录原码 这里就有两份信息了 长度所需要的比特位小于原来的值 就达到了压缩的目的
	原码的具体信息在表 distanceZone 里面的具体数值是现成找来的 总共有30个 返回是 [0, 29]
	所以记录下来的值就是如下
	0 5 0 3 | 2 0 0 0 | 0 0 0 0 | 2 0 0 0 | 0 5 0 0 | 0 0 0 5 | 0 0 0 5 | 0 0

	然后是根据这个信息反推出原码映射表 这里用到了deflate的两个特性
	1、整棵树最左边叶子节点的码字为0（码字长度视情况而定）
	2、树深为（n+1）时，该层最左面的叶子节点的值为，树深为n的这一层最左面的叶子节点的值加上该层所有叶子节点的个数，然后变长一位（即左移一位）。
 可以对照上面的树验证一下

	按照上面的规则 长度(len)最短之中下标(index)最小的那个值就可以得到一个映射  index   -- 0 * len (这里的*表示个数)
也就是 4 -- 00 ， 然后根据`同一深度的子节点 右边的值一定比左边的大` 这一特性 找到 12  -- 01 的映射
	接下来看树深(n) 加1的值 根据上面的规则 3 -- (00 + 这一层的个数) << 1  得到 100  下面的以此类推

*/
func buildCodeMapByBits(bits []byte) [][]byte {

	m := make([][]byte, len(bits))
	deepth := getMaxDeepth(bits)
	streamTemp := make([][]uint16, deepth+1)
	for l := 0; l < len(streamTemp); l++ {
		streamTemp[l] = make([]uint16, 0, 32)
	}

	for sourceCode, huffmanLen := range bits {
		if huffmanLen != 0 {
			streamTemp[huffmanLen] = append(streamTemp[huffmanLen], uint16(sourceCode))
		}
	}

	//表示是否出现了长度不为0的值 这里的index表示huffman码的长度 这里是根据deflate的规律来的
	flag := false
	//记录树深(n-1)的第一个码的值和长度
	var lastCode, lastLength uint32
	for huffmanLen := 1; huffmanLen < len(streamTemp); huffmanLen++ {
		if len(streamTemp[huffmanLen]) == 0 && flag == false {
			continue
		}
		//deflate树的最左边的节点始终为0
		lastCode = (lastCode + lastLength) << 1
		tempCode := lastCode

		for i := 0; i < len(streamTemp[huffmanLen]); i++ {
			bytes := make([]byte, huffmanLen)
			for t := 0; t < huffmanLen; t++ {
				bytes[huffmanLen-1-t] = utils.ReadBitLow(tempCode, uint(t))
			}
			m[streamTemp[huffmanLen][i]] = bytes
			tempCode++
		}
		lastLength = uint32(len(streamTemp[huffmanLen]))
		flag = true
	}
	return m
}

//建立deflate树根据map
func buildDeflatTreeByMap(m [][]byte) *HuffmanNode {

	root := &HuffmanNode{}
	var temp *HuffmanNode
	for k, v := range m {
		temp = root
		for index, bit := range v {
			if bit == 0 {
				if temp.LeftTree != nil && index == len(v)-1 {
					panic("buildDeflatTreeByMap error LeftTree")
				}
				if temp.LeftTree != nil && index != len(v)-1 && temp.LeftTree.Leaf == true {
					panic("buildDeflatTreeByMap error LeftTree 2")
				}
				if temp.LeftTree != nil {
					temp = temp.LeftTree
				} else {
					newTemp := &HuffmanNode{}
					temp.LeftTree = newTemp
					temp = newTemp
					if index == len(v)-1 {
						temp.Value = uint16(k)
						temp.Leaf = true
					}
				}
			} else if bit == 1 {
				if temp.RightTree != nil && index == len(v)-1 {
					panic("buildDeflatTreeByMap error RightTree")
				}
				if temp.RightTree != nil && index != len(v)-1 && temp.RightTree.Leaf == true {
					panic("buildDeflatTreeByMap error RightTree 2")
				}
				if temp.RightTree != nil {
					temp = temp.RightTree
				} else {
					newTemp := &HuffmanNode{}
					temp.RightTree = newTemp
					temp = newTemp
					if index == len(v)-1 {
						temp.Value = uint16(k)
						temp.Leaf = true
					}
				}
			} else {
				panic("buildDeflatTreeByMap error")
			}
		}
	}
	return root
}
