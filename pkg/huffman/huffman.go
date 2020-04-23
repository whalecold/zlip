package huffman

import (
	"container/list"
	"sort"

	"github.com/whalecold/zlip/pkg/utils"
)

// 因为码树节点左右分支旋转不会影响压缩程度 所有 huffman 树有很多表示
// 这里采用最不平衡的树 称之为 Deflate 树 在构建的时候效率也会更高 不需要递归
func buildTreeFromNodes(huffmanNodes []*treeNode) *treeNode {
	if len(huffmanNodes) == 0 {
		panic("buildTreeFromNodes param length is zero!")
	}
	index := uint16(0)
	for len(huffmanNodes) != 1 {
		sort.Slice(huffmanNodes, func(i, j int) bool {
			if huffmanNodes[i].Power != huffmanNodes[j].Power {
				return huffmanNodes[i].Power < huffmanNodes[j].Power
			}
			return huffmanNodes[i].Value < huffmanNodes[j].Value
		})
		newNode := &treeNode{
			Power:     huffmanNodes[0].Power + huffmanNodes[1].Power,
			LeftTree:  huffmanNodes[0],
			Value:     index,
			RightTree: huffmanNodes[1],
		}
		index++
		huffmanNodes[0] = newNode
		huffmanNodes = append(huffmanNodes[:1], huffmanNodes[2:]...)
	}
	return moveHuffmanToDeflateTree(huffmanNodes[0])
}

func cleanChildNode(node *treeNode) *treeNode {
	node.LeftTree = nil
	node.RightTree = nil
	return node
}

// moveHuffmanToDeflateTree transfers huffman tree to deflate
func moveHuffmanToDeflateTree(root *treeNode) *treeNode {
	var curCol, nextCol []*treeNode

	l := list.New()
	l.PushBack(root)
	for l.Len() != 0 {
		queueLen := l.Len()
		// stores cur col and next col spin nodes which has no affect on compress effectiveness.
		curCol = make([]*treeNode, 0, queueLen)
		nextCol = make([]*treeNode, 0, 2*queueLen)
		for i := 0; i < queueLen; i++ {

			element := l.Front()
			node := element.Value.(*treeNode)
			l.Remove(element)

			if node.LeftTree != nil {
				nextCol = append(nextCol, node.LeftTree)
			}
			if node.RightTree != nil {
				nextCol = append(nextCol, node.RightTree)
			}

			// cleans the left and right tree to reconnect the next col nodes.
			curCol = append(curCol, cleanChildNode(node))
		}

		// sorts the nodes, the make sure leaf node stores in the front of none leaf node,
		// if the leafs has the same attributes, sort by value.
		sort.SliceStable(nextCol, func(i, j int) bool {
			if !nextCol[i].Leaf && nextCol[j].Leaf {
				return false
			}
			if nextCol[i].Leaf && !nextCol[j].Leaf {
				return true
			}
			return nextCol[i].Value < nextCol[j].Value
		})

		// reconnects the cur col and next col in reverse order.
		for i := 0; i < len(curCol); i++ {
			nextColIndex := len(nextCol) - 1 - 2*i
			curColIndex := len(curCol) - i - 1

			connect := func(child **treeNode, nextIndex int) bool {
				if nextColIndex >= 0 {
					*child = nextCol[nextColIndex]
					l.PushFront(nextCol[nextColIndex])
					return true
				}
				return false
			}
			if !connect(&curCol[curColIndex].RightTree, nextColIndex) {
				break
			}
			nextColIndex -= 1
			if !connect(&curCol[curColIndex].LeftTree, nextColIndex) {
				break
			}
		}
	}
	return root
}

func transferBytesToTreeNodes(bytes []byte) []*treeNode {
	huffmanMap := make(map[byte]*treeNode)
	for _, value := range bytes {
		if m, ok := huffmanMap[value]; ok {
			m.Power++
		} else {
			huffmanMap[value] = &treeNode{
				Power: 1,
				Value: uint16(value),
				Leaf:  true,
			}
		}
	}

	huffmanNodes := make([]*treeNode, 0, len(huffmanMap))
	for _, v := range huffmanMap {
		huffmanNodes = append(huffmanNodes, v)
	}
	return huffmanNodes
}

func buildHuffmanTreeFromBytes(bytes []byte) *treeNode {
	nodes := transferBytesToTreeNodes(bytes)
	return buildTreeFromNodes(nodes)
}

// 根据位数来重新获得码表
//
//
// https://blog.csdn.net/jison_r_wang/article/details/52071841
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

	然后是根据这个信息反推出原码映射表 这里用到了 deflate 的两个特性
	1、整棵树最左边叶子节点的码字为0（码字长度视情况而定）
	2、树深为 n+1 时，该层最左面的叶子节点的值为，树深为 n 的这一层最左面的叶子节点的值加上该层所有叶子节点的个数，然后变长一位（即左移一位）。
 可以对照上面的树验证一下

	按照上面的规则 长度 len 最短之中下标 index 最小的那个值就可以得到一个映射  index   -- 0 * len (这里的*表示个数)
也就是 4 -- 00 ， 然后根据`同一深度的子节点 右边的值一定比左边的大` 这一特性 找到 12  -- 01 的映射
	接下来看树深 n 加1的值 根据上面的规则 3 -- (00 + 这一层的个数) << 1  得到 100  下面的以此类推

*/

func buildCodeMapFromBits(bits []byte) [][]byte {

	codeMap := make([][]byte, len(bits))
	bitStream := make([][]uint16, getMaxDepth(bits)+1)
	for l := 0; l < len(bitStream); l++ {
		bitStream[l] = make([]uint16, 0, 32)
	}

	for sourceCode, huffmanLen := range bits {
		if huffmanLen != 0 {
			bitStream[huffmanLen] = append(bitStream[huffmanLen], uint16(sourceCode))
		}
	}

	// 表示是否出现了长度不为0的值 这里的 index 表示 huffman 码的长度 这里是根据 deflate 的规律来的
	flag := false
	// 记录树深 n-1 的第一个码的值和长度
	var lastCode, lastLength uint32
	for huffmanLen := 1; huffmanLen < len(bitStream); huffmanLen++ {
		if len(bitStream[huffmanLen]) == 0 && !flag {
			continue
		}
		// deflate 树的最左边的节点始终为0
		lastCode = (lastCode + lastLength) << 1
		tempCode := lastCode

		for i := 0; i < len(bitStream[huffmanLen]); i++ {
			bytes := make([]byte, huffmanLen)
			for t := 0; t < huffmanLen; t++ {
				bytes[huffmanLen-1-t] = utils.GetLowBit32(tempCode, uint(t))
			}
			codeMap[bitStream[huffmanLen][i]] = bytes
			tempCode++
		}
		lastLength = uint32(len(bitStream[huffmanLen]))
		flag = true
	}
	return codeMap
}

// buildDeflateTreeFromMap generator a deflate tree by map
func buildDeflateTreeFromMap(codeMap [][]byte) *treeNode {
	root := &treeNode{}
	var temp *treeNode
	for k, v := range codeMap {
		temp = root
		for index, bit := range v {
			if bit == 0 {
				if temp.LeftTree != nil && index == len(v)-1 {
					panic("buildDeflateTreeFromMap error LeftTree")
				}
				if temp.LeftTree != nil && index != len(v)-1 && temp.LeftTree.Leaf {
					panic("buildDeflateTreeFromMap error LeftTree 2")
				}
				if temp.LeftTree != nil {
					temp = temp.LeftTree
				} else {
					newTemp := &treeNode{}
					temp.LeftTree = newTemp
					temp = newTemp
					if index == len(v)-1 {
						temp.Value = uint16(k)
						temp.Leaf = true
					}
				}
			} else if bit == 1 {
				if temp.RightTree != nil && index == len(v)-1 {
					panic("buildDeflateTreeFromMap error RightTree")
				}
				if temp.RightTree != nil && index != len(v)-1 && temp.RightTree.Leaf {
					panic("buildDeflateTreeFromMap error RightTree 2")
				}
				if temp.RightTree != nil {
					temp = temp.RightTree
				} else {
					newTemp := &treeNode{}
					temp.RightTree = newTemp
					temp = newTemp
					if index == len(v)-1 {
						temp.Value = uint16(k)
						temp.Leaf = true
					}
				}
			} else {
				panic("buildDeflateTreeFromMap error")
			}
		}
	}
	return root
}
