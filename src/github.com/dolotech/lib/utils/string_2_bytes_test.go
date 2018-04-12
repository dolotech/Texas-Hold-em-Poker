package utils

import (
	"testing"
)

/*func BenchmarkBytes2String(t *testing.B) {
	for i := 1; i < N; i++ {
		b := []byte("34asdfasdfasdfasdfasdfasdfasdfasdf12")
		t.Log(Bytes2String(b))
	}
}*/
func BenchmarkTestString2Bytes(t *testing.B) {
	for i := 1; i < N; i++ {
		b := "34asdfasdfasdfasdfasdfasdfasdfasdf12"
		t.Log(Bytes2String(String2Bytes(b)))
	}
}

const N = 30000
/*func BenchmarkBytes2String1(t *testing.B) {
	for i := 1; i < N; i++ {
		b := []byte("34asdfasdfasdfasdfasdfasdfasdfasdf12")
		t.Log(Bytes2String(b))
	}
}*/
func BenchmarkTestString2Bytes1(t *testing.B) {
	for i := 1; i < N; i++ {
		b := "34asdfasdfasdfasdfasdfasdfasdfasdf12"
		t.Log(string([]byte(b)))
	}
}
