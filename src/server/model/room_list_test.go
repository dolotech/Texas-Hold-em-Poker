package model

import (
	"testing"
	"time"
)

func BenchmarkCreateNumber(b *testing.B) {
	for i:=0;i<b.N;i++{

		b.Log(createNumber())
	}
}













func Benchmark_map(b *testing.B) {
	m:= make(map[int]int)
	for i:=0;i<b.N;i++{
		go func() {
			defer func() {
				if err:=recover();err!=nil{
					b.Log(err)
				}
			}()
			m[i] = i
		}()
	}


	<-time.After(time.Minute)
}