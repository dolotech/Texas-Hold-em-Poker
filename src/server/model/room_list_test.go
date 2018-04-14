package model

import "testing"

func BenchmarkCreateNumber(b *testing.B) {
	for i:=0;i<b.N;i++{

		b.Log(createNumber())
	}

}