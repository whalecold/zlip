package lz77

import (
	"compress/algorithm/huffman"
	"compress/algorithm/stack"
)

func dealWithBytesAndStack(stackNode *stack.Stack, r *[]byte) {
	temp := make([]byte, 0, RLC_MaxLength)
	for node := stackNode.RPop(); node != nil; {
		temp = append(temp, node.(byte))
		node = stackNode.RPop()
	}
	if temp[0] == RLC_Zero && len(temp) >= RLC_Length {
		tempLen := len(temp)
		for tempLen > huffman.HUFFMAN_CCLLen {
			*r = append(*r, RLC_Special)
			*r = append(*r, huffman.HUFFMAN_CCLLen)
			tempLen -= huffman.HUFFMAN_CCLLen
		}

		if tempLen != 0 {
			*r = append(*r, RLC_Special)
			*r = append(*r, byte(tempLen))
		}

	} else {
		for _, value := range temp {
			*r = append(*r, value)
		}
	}
}

//游程编码
//run length coding 这里感觉非0的重复不会很多 只对0进行编码
//简单处理下 17表示0 后面的数字表示0重复的个数 多于3个重复才开始编码
func RLC(bytes []byte) []byte {
	result := make([]byte, 0, len(bytes))
	stackNode := stack.NewStack()

	stackNode.Push(bytes[0])

	for i := 1; i < len(bytes); i++ {
		lastNode := stackNode.Pop()
		if lastNode != nil {
			lastData := lastNode.(byte)
			if lastData != bytes[i] {
				dealWithBytesAndStack(stackNode, &result)
			}
		}
		stackNode.Push(bytes[i])
	}
	dealWithBytesAndStack(stackNode, &result)
	return bytes
}

func UnRLC(bytes []byte) []byte {
	result := make([]byte, 0, RLC_MaxLength)
	for i := 0; i < len(bytes); i++ {
		if bytes[i] == RLC_Special {
			tempLen := bytes[i+1]
			for k := byte(0); k < tempLen; k++ {
				result = append(result, RLC_Zero)
			}
			i++
		} else {
			result = append(result, bytes[i])
		}
	}
	return bytes
}
