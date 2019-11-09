package huffman

//Alg alg
type Alg struct {
	tree *DeflateTree
	//common DeflateCommon
}

//InitDis init dis
func (h *Alg) InitDis() {
	h.tree = &DeflateTree{condition: &Distance{extraCode: DistanceZone}}
	h.tree.Init()
	//h.common = &Distance{extraCode:DistanceZone}
}

//InitCCL init ccl
func (h *Alg) InitCCL() {
	h.tree = &DeflateTree{condition: &CCl{}}
	h.tree.Init()
	//h.common = &CCl{}
}

//InitLiteral init
func (h *Alg) InitLiteral() {
	h.tree = &DeflateTree{condition: &Literal{extraCode: LengthZone}}
	h.tree.Init()
	//h.common = &Literal{extraCode:LengthZone}
}

//AddElement add element
//runtime.mapaccess2 占了20%左右
func (h *Alg) AddElement(element uint16, length bool) {
	h.tree.AddElement(element, length)
}

//BuildHuffmanMap build map
func (h *Alg) BuildHuffmanMap() {
	h.tree.BuildTree()
	h.tree.BuildMap()
}

//EnCodeElement encode element
func (h *Alg) EnCodeElement(ele uint16,
	bytes *[]byte,
	offset uint32,
	dataSet *uint64,
	length bool) uint32 {
	return h.tree.EnCodeElement(ele, bytes, offset, dataSet, length)
}

//SerializeBitsStream serialze bit stream
func (h *Alg) SerializeBitsStream() ([]byte, uint32) {
	h.tree.SerializeBitsStream()
	return h.tree.GetBits(), h.tree.BitesLen()
}

//UnSerializeAndBuild unserial
func (h *Alg) UnSerializeAndBuild(bits []byte) {
	h.tree.UnSerializeBitsStream(bits)
	h.tree.BuildTreeByMap()
}

//DecodeEle nil
func (h *Alg) DecodeEle(bytes []byte, offset uint32) (uint16, uint32, uint32, bool) {
	return h.tree.DecodeEle(bytes, offset)
}

//GetBits get bits
func (h *Alg) GetBits() []byte {
	return h.tree.GetBits()
}
