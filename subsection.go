package main

type Subsection struct {
	Sequence int64
	Content  []byte
}

type SubsectionSlice []*Subsection

func (s SubsectionSlice) Less(i, j int) bool {
	return s[i].Sequence < s[j].Sequence
}

func (s SubsectionSlice) Swap(i, j int) {
	s[i].Sequence, s[j].Sequence = s[j].Sequence, s[i].Sequence
	s[i].Content, s[j].Content = s[j].Content, s[i].Content
}

func (s SubsectionSlice) Len() int {
	return len(s)
}
