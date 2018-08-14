package huffman

const (
	HUFFMAN_LiteralLimit = 256 //literal 和 length的分界线
	HUFFMAN_EndFlag      = 256 //结束标志
	//CCL树的长度
	//zip中会进行剪枝 不会超过15的
	HUFFMAN_CCLLen     = 32
	HUFFMAN_CCLEndFlag = HUFFMAN_CCLLen + 1 //结束标志

	HUFFMAN_ElementNum = 2 ^ 16
)
