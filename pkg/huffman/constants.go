package huffman

const (
	// LiteralBoundary
	LiteralBoundary uint16 = 256 //literal 和 length的分界线
	// EndFlag nil
	EndFlag = 256 //结束标志
	// CCLLen nil
	// CCL树的长度
	// zip中会进行剪枝 不会超过15的
	CCLLen = 32
	// CCLEndFlag indicates the  of ccl end flag.
	CCLEndFlag = CCLLen + 1

	// ElementNum nil
	ElementNum = 2 ^ 16
)
