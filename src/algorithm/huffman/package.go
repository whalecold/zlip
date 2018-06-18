package huffman

type HuffmanAlg struct {
	tree *DeflateTree
	//common DeflateCommon
}

func (h *HuffmanAlg)InitDis() {
	h.tree = &DeflateTree{condition:&Distance{extraCode:DistanceZone}}
	h.tree.Init()
	//h.common = &Distance{extraCode:DistanceZone}
}

func (h *HuffmanAlg)InitCCL() {
	h.tree = &DeflateTree{condition:&CCl{}}
	h.tree.Init()
	//h.common = &CCl{}
}

func (h *HuffmanAlg)InitLiteral() {
	h.tree = &DeflateTree{condition:&Literal{extraCode:LengthZone}}
	h.tree.Init()
	//h.common = &Literal{extraCode:LengthZone}
}

func (h *HuffmanAlg)AddElement(element uint16, length bool) {
	h.tree.AddElement(element, length)
}

func (h *HuffmanAlg)BuildHuffmanMap() {
	h.tree.BuildTree()
	h.tree.BuildMap()
}

func (h *HuffmanAlg)EnCodeElement(ele uint16,
									bytes *[]byte,
									offset uint32,
									dataSet *uint64,
									length bool ) (uint32) {
	return h.tree.EnCodeElement(ele, bytes, offset, dataSet, length)
}

func (h *HuffmanAlg)SerializeBitsStream() ([]byte, uint32) {
	h.tree.SerializeBitsStream()
	return h.tree.GetBits(), h.tree.BitesLen()
}

func (h *HuffmanAlg)UnSerializeAndBuild(bits  []byte) {
	h.tree.UnSerializeBitsStream(bits)
	h.tree.BuildTreeByMap()
}

func (h *HuffmanAlg)DecodeEle(bytes []byte,
								offset uint32)(uint16, uint32, uint32, bool) {
	return h.tree.DecodeEle(bytes, offset)
}