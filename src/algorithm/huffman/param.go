package huffman

const (
	HUFFMAN_LiteralLimit = 256		//literal 和 length的分界线
	HUFFMAN_EndFlag 	 = 256 		//结束标志
	//CCL树的长度
	HUFFMAN_CCLLen 		 = 18
	HUFFMAN_CCLEndFlag 	 = HUFFMAN_CCLLen + 1 		//结束标志
)