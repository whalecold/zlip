package huffman

type HuffmanAlg struct {
	tree *DeflateTree
	common DeflateCommon
}

func (h *HuffmanAlg)InitDis() {
	h.tree = &DeflateTree{}
	h.tree.Init(DistanceZone)
	h.common = &Distance{}
}

func (h *HuffmanAlg)InitLiteral() {
	h.tree = &DeflateTree{}
	h.tree.Init(LengthZone)
	h.common = &Literal{}
}

func (h *HuffmanAlg)AddElement(element uint16, length bool) {
	h.tree.AddElement(element, h.common, length)
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
	return h.tree.EnCodeElement(ele, bytes, offset, dataSet, h.common, length)
}

func (h *HuffmanAlg)SerializeBitsStream() ([]byte, uint32) {
	h.tree.SerializeBitsStream(h.common)
	return h.tree.GetBits(), h.tree.BitesLen()
}

func (h *HuffmanAlg)UnSerializeAndBuild(bits  []byte) {
	h.tree.UnSerializeBitsStream(bits, h.common)
	h.tree.BuildTreeByMap()
}

func (h *HuffmanAlg)DecodeEle(bytes []byte,
								offset uint32)(uint16, uint32, uint32, bool) {
	return h.tree.DecodeEle(bytes, offset, h.common)
}