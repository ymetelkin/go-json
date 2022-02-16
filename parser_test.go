package v2

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"
)

var (
	data = make([][]byte, 3)
	keys = []string{"hits", "messi", "search"}
)

func benchmarkParsing(b *testing.B, key int) {
	buf := data[key]
	if len(buf) == 0 {
		path := fmt.Sprintf("test_data/%s.json", keys[key])
		f, _ := os.Open(path)
		buf, _ = io.ReadAll(f)
		f.Close()
		data[key] = buf
	}

	b.ResetTimer()

	for x := 0; x < b.N; x++ {
		ParseObject(buf)
	}
}

func benchmarkParsingGo(b *testing.B, key int) {
	buf := data[key]
	if len(buf) == 0 {
		path := fmt.Sprintf("test_data/%s.json", keys[key])
		f, _ := os.Open(path)
		buf, _ = io.ReadAll(f)
		f.Close()
		data[key] = buf
	}

	b.ResetTimer()

	for x := 0; x < b.N; x++ {
		jo := make(map[string]interface{})
		json.Unmarshal(buf, &jo)
	}
}

func BenchmarkHits(b *testing.B) {
	benchmarkParsing(b, 0)
}

func BenchmarkHitsGo(b *testing.B) {
	benchmarkParsingGo(b, 0)
}

func BenchmarkMessi(b *testing.B) {
	benchmarkParsing(b, 1)
}

func BenchmarkMessiGo(b *testing.B) {
	benchmarkParsingGo(b, 1)
}

func BenchmarkSearch(b *testing.B) {
	benchmarkParsing(b, 2)
}

func BenchmarkSearchGo(b *testing.B) {
	benchmarkParsingGo(b, 2)
}
