package huffman

const (
	//HUFFMANLiteralLimit nil
	HUFFMANLiteralLimit = 256 //literal 和 length的分界线
	//HUFFMANEndFlag nil
	HUFFMANEndFlag = 256 //结束标志
	//HUFFMANCCLLen nil
	//CCL树的长度
	//zip中会进行剪枝 不会超过15的
	HUFFMANCCLLen = 32
	//HUFFMANCCLEndFlag nil
	HUFFMANCCLEndFlag = HUFFMANCCLLen + 1 //结束标志

	//HUFFMANElementNum nil
	HUFFMANElementNum = 2 ^ 16
)
