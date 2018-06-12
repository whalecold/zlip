package lz77

const (
	LZ77_MaxWindowsSize = 1024 * 32		//滑动窗口的大小
	LZ77_CmpHeadSize 	= LZ77_MaxWindowsSize
	LZ77_CmpPrevSize 	= LZ77_MaxWindowsSize
	LZ77_WindowsMask 	= 0x7FFF
)