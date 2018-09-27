package utils

//ReadBitLow read
//备注：这里面的判断都可以不加 最好是在debuf版本加 release版本不加
//优化压缩耗时
//bit 范围 0-31 0表示最低位 从最低位开始读
func ReadBitLow(num uint32, bit uint) byte {
	//if bit >= 32 {
	//	panic("readBit error")
	//}
	temp := uint32(1 << bit)
	if num&temp == 0 {
		return 0
	}
	return 1
}

//ReadBitsHigh read
//offset0表示最高位 从高位开始读 [0,7]
func ReadBitsHigh(b byte, offset uint32) byte {
	//if offset > 7 {
	//	panic("ReadBitsHigh error offset")
	//}
	move := 7 - offset
	b = b >> move
	b &= 0x1
	return b
}

//ReadBitsHigh16 read
//offset0表示最高位 从高位开始读 [0,15]
func ReadBitsHigh16(b uint16, offset uint32) byte {
	//if offset > 15 {
	//	panic("ReadBitsHigh error offset")
	//}
	move := 15 - offset
	b = b >> move
	b &= 0x1
	return byte(b)
}

//WriteBitsHigh write
//offset0表示最高位 从高位开始写 [0,7] 把某一位置为 n
func WriteBitsHigh(b *byte, offset uint32, n byte) byte {
	//把这两个注释了 测试的使用再加上吧
	//if offset > 7 {
	//	panic("WriteBitsHigh error offset")
	//}
	//if n != 0 && n != 1 {
	//	panic("WriteBitsHigh error n")
	//}

	i := n << uint32(7-offset)
	if n == 1 {
		*b |= i
	} else {
		*b &= ^i
	}
	return *b
}

//WriteBitsHigh16 write
//offset0表示最高位 从高位开始写 [0,15] 把某一位置为 n
func WriteBitsHigh16(b *uint16, offset uint32, n uint16) uint16 {
	//if offset > 15 {
	//	panic("WriteBitsHigh error offset")
	//}
	//if n != 0 && n != 1 {
	//	panic("WriteBitsHigh error n")
	//}

	if n == 1 {
		i := n << uint32(15-offset)
		*b = *b | i
	} else {
		i := uint16(1 << uint32(15-offset))
		*b = *b & (^i)
	}
	return *b
}

//ReadBitsLen read
//从bytes bitOffset readLen位数据 返回 值 byte偏移 bits偏移
func ReadBitsLen(bytes []byte, bitOffset uint32, readLen uint16) (uint16, uint32, uint32) {
	if readLen == 0 {
		return 0, 0, bitOffset
	}
	var byteLen uint32
	var getLen uint16
	var result uint16
	for _, value := range bytes {
		for ; bitOffset < 8; bitOffset++ {
			bit := ReadBitsHigh(value, bitOffset)
			result = result << 1
			//result = result ^ uint16(bit)
			result = result | uint16(bit)
			getLen++
			if getLen >= readLen {
				return result, byteLen, bitOffset + 1
			}
		}
		bitOffset = 0
		byteLen++
	}
	//走到这个肯定是程序出错了 找不到对应字符串是不能发生的
	panic("ReadBitsLen failed !")
}

//hotfix test
//branchTest
