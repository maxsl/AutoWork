package intTobyte

import (
	"math/rand"
	"testing"
)

func Test_des(t *testing.T) {
	buf := make([]byte, 4)
	for i := 0; i < 4; i++ {
		buf[i] = byte(rand.Intn(255))
	}
	ToByte(ToInt(buf))
}
func Benchmark_des(t *testing.B) {
	for i := 0; i < t.N; i++ {
		buf := make([]byte, 4)
		for i := 0; i < 4; i++ {
			buf[i] = byte(rand.Intn(255))
		}
		ToByte(ToInt(buf))
	}
}
