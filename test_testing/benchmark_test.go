package test_testing

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

// 利用+=连接
func TestConcatStringByAdd(t *testing.T) {
	assert := assert.New(t)
	elems := []string{"1", "2", "3", "4", "5"}
	ret := ""
	for _, elem := range elems {
		ret += elem
	}
	assert.Equal("12345", ret)
}

// 利用buffer连接
func TestConcatStringBytesBuffer(t *testing.T) {
	assert := assert.New(t)
	var buf bytes.Buffer
	elems := []string{"1", "2", "3", "4", "5"}
	for _, elem := range elems {
		buf.WriteString(elem)
	}
	assert.Equal("12345", buf.String())
}

func BenchmarkConcatStringByAdd(b *testing.B) {
	elems := []string{"1", "2", "3", "4", "5"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ret := ""
		for _, elem := range elems {
			ret += elem
		}
	}
	b.StopTimer()
}

func BenchmarkConcatStringBytesBuffer(b *testing.B) {
	elems := []string{"1", "2", "3", "4", "5"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		for _, elem := range elems {
			buf.WriteString(elem)
		}
	}
}

//go test -bench=".*BenchmarkConcat.*" -benchmem
//go test -v -run TestA select_test.go   这里指定 TestA 进行测试

//-benchtime t
//-count n
//-cpu n
//-benchmem

//https://www.cnblogs.com/yahuian/p/14461152.html

//https://geektutu.com/post/hpg-benchmark.html

//http://c.biancheng.net/view/124.html
