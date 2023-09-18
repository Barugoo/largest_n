package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLargestN(t *testing.T) {
	input :=
		`https://url.url.url/bbb 2222222
	https://url.url.url/aaa 1111111
	https://url.url.url/ccc 3333333
	`
	t.Run("normal", func(t *testing.T) {
		expectedIndexies := []ElemIdx{
			{
				BatchNumber: 0,
				RowNumber:   0,
			},
			{
				BatchNumber: 2,
				RowNumber:   0,
			},
		}
		indexies, err := largestN(strings.NewReader(input), 1, 2)
		assert.NoError(t, err, "largestN should not return err")
		assert.Equal(t, expectedIndexies, indexies, "expected and actual indexies should match")
	})

	t.Run("zero batch size", func(t *testing.T) {
		_, err := largestN(strings.NewReader(input), 0, 2)
		assert.Error(t, err, "largestN should return err")
	})

}

func TestLargestNBigSize(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	batchSize := 1000
	n := 10
	var totalAmount int64 = 1000000

	r, w := io.Pipe()
	go func() {
		for i := int64(0); i < totalAmount; i++ {
			row := fmt.Sprintf("https://url.url.url/%d %d\n", i, i)
			w.Write([]byte(row))
		}
		w.Close()
	}()

	var expectedIndexies []ElemIdx
	for i := 0; i < n; i++ {
		expectedIndexies = append(expectedIndexies, ElemIdx{
			BatchNumber: totalAmount/int64(batchSize) - 1,
			RowNumber:   batchSize - n + i,
		})
	}
	indexies, err := largestN(r, batchSize, n)
	assert.NoError(t, err, "largestN should not return err")
	assert.Equal(t, expectedIndexies, indexies, "expected and actual indexies should match")
}

func BenchmarkLargest10(b *testing.B) {
	f, err := os.Open("file.txt")
	if err != nil {
		b.FailNow()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		largestN(f, 10000, 10)
		f.Seek(0, 0)
	}

}
