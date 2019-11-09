package utils

// GetLowBit32 gets a bit from a uint32 starts at low position which is less important than high position, offset is
// o. 0 means lowest  position and 31 means highest position.
func GetLowBit32(num uint32, o uint) byte {
	//if bit >= 32 {
	//	panic("readBit error")
	//}
	if num&uint32(1<<o) == 0 {
		return 0
	}
	return 1
}

// GetHighBit8 gets a bit from a byte starts at high position which is more important than low position, offset is
// o. 0 means highest position and 7 means lowest position.
func GetHighBit8(b byte, o uint32) byte {
	if o > 7 {
		panic("ReadBitsHigh error offset")
	}
	move := 7 - o
	b = b >> move
	b &= 0x1
	return b
}

// GetHighBit16 gets a bit from a uint16 starts at high position which is more important than low position, offset is
// o. 0 means highest position and 15 means lowest position.
func GetHighBit16(b uint16, offset uint32) byte {
	//if offset > 15 {
	//	panic("ReadBitsHigh error offset")
	//}
	move := 15 - offset
	b = b >> move
	b &= 0x1
	return byte(b)
}

// SetHighBit8 sets the bit to s from a byte starts at high position which is more important than low position, offset is
// o. 0 means highest position and 7 means lowest position.
func SetHighBit8(b *byte, o uint32, s byte) byte {
	if o > 7 {
		panic("WriteBitsHigh error offset")
	}
	if s != 0 && s != 1 {
		panic("WriteBitsHigh error n")
	}

	i := s << uint32(7-o)
	if s == 1 {
		*b |= i
	} else {
		*b &= ^i
	}
	return *b
}

// SetHighBit16 sets the bit to s from a uint16 starts at high position which is more important than low position, offset is
// o. 0 means highest position and 15 means lowest position.
func SetHighBit16(b *uint16, o uint32, s uint16) uint16 {
	//if offset > 15 {
	//	panic("WriteBitsHigh error offset")
	//}
	//if n != 0 && n != 1 {
	//	panic("WriteBitsHigh error n")
	//}

	if s == 1 {
		i := s << uint32(15-o)
		*b = *b | i
	} else {
		i := uint16(1 << uint32(15-o))
		*b = *b & (^i)
	}
	return *b
}

// ReadBitsLen reads len bits as uint16 from bytes starts at bf bit offset.
// Return the required value、 byte slide length、bit slide length from a new byte
func ReadBitsLen(bytes []byte, bf uint32, len uint16) (result uint16, bsl uint32, offset uint32) {
	offset = bf
	if len == 0 {
		return
	}
	// all slide len
	var asl uint16
	for _, value := range bytes {
		for ; offset < 8; offset++ {
			bit := GetHighBit8(value, offset)
			result = result << 1
			// result = result ^ uint16(bit)
			result = result | uint16(bit)
			asl++
			// if the len of all slide meets the length, return
			if asl >= len {
				offset += 1
				return
			}
		}
		// clean offset when slide one byte and add byte slide len
		offset = 0
		bsl++
	}
	// it's can't happened
	panic("ReadBitsLen failed !")
}
