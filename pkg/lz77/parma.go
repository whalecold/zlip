package lz77

const (
	// MaxWindowsSize max slide window size
	MaxWindowsSize = 32 * 1024
	// CmpHeadSize nil
	CmpHeadSize = MaxWindowsSize
	// CmpPrevSize nil
	CmpPrevSize = MaxWindowsSize
	// WindowsMask nil
	WindowsMask = MaxWindowsSize - 1
	// MinCmpSize nil
	MinCmpSize = 3
	// MaxCmpNum nil
	MaxCmpNum = 8 //最大比较次数
	// MaxCmpLength max comparison length
	MaxCmpLength = 256 //最大比较长度
	// EndFlag indicates the end.
	EndFlag = 256
	// HeadInfo head info
	HeadInfo = "zls1129@gmail.com version 1.0.1"
	// ChunkSize separates the whole file into several chunks and every chunk's length is ChunkSize.
	// Avoid the file size is too large to occupy to many memory and we can do parallel processing with several chunks.
	ChunkSize = 1024 * 1024 * 5
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
