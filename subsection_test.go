package main

import (
	"fmt"
	"sort"
	"testing"
)

func TestSubsectionSlice_Len(t *testing.T) {
	s := SubsectionSlice{}
	s = append(s, &Subsection{2, []byte("123")})
	s = append(s, &Subsection{1, []byte("456")})
	s = append(s, &Subsection{0, []byte("789")})
	sort.Sort(s)
	for _, v := range s {
		fmt.Printf("s : %v\n", string(v.Content))

	}
}
