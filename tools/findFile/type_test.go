package findFile

import (
	"fmt"
	"testing"
)

func Test_match(t *testing.T) {
	m := NewMatch(0, "*.json")
	fmt.Println(m.Walk("./"))
}

func Benchmark_match(t *testing.B) {
	m := NewMatch(0, "*.go")
	for i := 0; i < t.N; i++ {
		m.Walk("./")
	}
}
