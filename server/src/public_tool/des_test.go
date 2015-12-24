package public_tool

import (
	"testing"
)

func Test_des(t *testing.T) {
	x := new(Des)
	x.Init([]byte("12345678"))
	buf, err := x.Encrypt([]byte("Hello wrold"))
	if err != nil {
		t.Log(err)
		return
	}
	buf, err = x.Decrypt(buf)
}

func Benchmark_des(t *testing.B) {
	x := new(Des)
	x.Init([]byte("12345678"))
	for i := 0; i < t.N; i++ {
		buf, err := x.Encrypt([]byte("Hello wrold"))
		if err != nil {
			t.Log(err)
			return
		}
		buf, err = x.Decrypt(buf)
	}
}
