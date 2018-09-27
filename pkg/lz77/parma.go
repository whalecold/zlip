package lz77

const (
	//LZ77MaxWindowsSize nil
	LZ77MaxWindowsSize = 32 * 1024 //滑动窗口的大小
	//LZ77CmpHeadSize nil
	//LZ7_MaxWindowsSize = 32 //滑动窗口的大小
	LZ77CmpHeadSize = LZ77MaxWindowsSize
	//LZ77CmpPrevSize nil
	LZ77CmpPrevSize = LZ77MaxWindowsSize
	//LZ77WindowsMask nil
	LZ77WindowsMask = LZ77MaxWindowsSize - 1
	//LZ77MinCmpSize nil
	LZ77MinCmpSize = 3
	//LZ77MaxCmpNum nil
	LZ77MaxCmpNum = 8 //最大比较次数
	//LZ77MaxCmpLength nil
	LZ77MaxCmpLength = 256 //最大比较长度
	//LZ77EndFlag nil
	LZ77EndFlag = 256 //結束
	//LZ77HeadInfo nil
	LZ77HeadInfo = "zls1129@gmail.com version 1.0.1"
	//LZ77ChunkSize nil
	LZ77ChunkSize = 1024 * 1024 * 5
)

const (
	//RLCLength nil
	RLCLength = 3 //游程编码最小匹配长度
	//RLCSpecial nil
	RLCSpecial = 17 //特殊字符
	//RLCMaxLength nil
	RLCMaxLength = 18 //最大长度
	//RLCZero nil
	RLCZero = 0 //特殊字符
	//RLCEndFlag nil
	RLCEndFlag = 255
)
