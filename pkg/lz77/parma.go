package lz77

const (
	LZ77_MaxWindowsSize =  32 * 1024		//滑动窗口的大小
	//LZ77_MaxWindowsSize = 32 //滑动窗口的大小
	LZ77_CmpHeadSize    = LZ77_MaxWindowsSize
	LZ77_CmpPrevSize    = LZ77_MaxWindowsSize
	LZ77_WindowsMask    = LZ77_MaxWindowsSize - 1
	LZ77_MinCmpSize     = 3
	LZ77_MaxCmpNum      = 8   //最大比较次数
	LZ77_MaxCmpLength   = 256 //最大比较长度
	LZ77_EndFlag        = 256 //結束
	LZ77_HeadInfo       = "zls1129@gmail.com version 1.0.1"

	LZ77_ChunkSize = 1024 * 1024 * 5
)

const (
	RLC_Length    = 3  //游程编码最小匹配长度
	RLC_Special   = 17 //特殊字符
	RLC_MaxLength = 18 //最大长度
	RLC_Zero      = 0  //特殊字符
	RLC_EndFlag   = 255
)
