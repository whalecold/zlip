package lz77

const (
	//LZ77_MaxWindowsSize =  32 * 1024		//滑动窗口的大小
	LZ77_MaxWindowsSize =  32 		//滑动窗口的大小
	LZ77_CmpHeadSize 	= LZ77_MaxWindowsSize
	LZ77_CmpPrevSize 	= LZ77_MaxWindowsSize
	LZ77_WindowsMask 	= LZ77_MaxWindowsSize - 1
	LZ77_MinCmpSize 	= 3
	LZ77_MaxCmpNum 		= 16 					//最大比较次数
	LZ77_MaxCmpLength 	= 256 					//最大比较长度
)