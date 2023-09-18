package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"io"
	"sort"
)

func largestN(r io.ReadSeeker, batchSize int, n int) []ElemIdx {
	batchBuf := make([]BatchElem, batchSize)
	sortBuf := make([]BatchElem, batchSize)

	elemHeap := make(BatchElemHeap, 0, n)
	heap.Init(&elemHeap)

	var cursour BatchElem

	scanner := bufio.NewScanner(r)
	scanner.Split(splitFn)

	for scanner.Scan() {
		var val int64
		// since scanner.Text allocates memory for new string each time
		// it's better to use scanner.Bytes and manually translate ASCII digits to int64
		for _, byt := range scanner.Bytes() {
			val = val*10 + int64(byt-'0')
		}

		cursour.Val = val
		batchBuf[cursour.RowNumber] = cursour

		if cursour.RowNumber == len(batchBuf)-1 {
			consumeBatch(&elemHeap, sortBuf, batchBuf, n)
			cursour.RowNumber = -1
			cursour.BatchNumber++
		}
		cursour.RowNumber++
	}
	consumeBatch(&elemHeap, sortBuf, batchBuf[:cursour.RowNumber], n)
	elemsInHeap := elemHeap.Len()

	// preparing final result
	res := make([]ElemIdx, elemsInHeap)
	for i := 0; i < elemsInHeap; i++ {
		res[i] = heap.Pop(&elemHeap).(BatchElem).ElemIdx
	}

	sort.Slice(res, func(i, j int) bool {
		if res[i].BatchNumber == res[j].BatchNumber {
			return res[i].RowNumber < res[j].RowNumber
		}
		return res[i].BatchNumber < res[j].BatchNumber
	})
	return res
}

// these are to save on allocations
var count = make([]int, 10)
var zeroes = make([]int, 10)

// consumeBatch takes a batch of elems and pushes n largest into heap
func consumeBatch(h heap.Interface, sortBuf, elems []BatchElem, n int) {
	var max BatchElem
	// finding max elem first will allow us to skip whole batch if max is smaller than last elem in heap
	for _, e := range elems {
		if e.Val > max.Val {
			max = e
		}
	}

	if h.Len() == n {
		minStored := heap.Pop(h).(BatchElem)
		if minStored.Val > max.Val {
			// skipping batch if smallest elem is greater than elems's max
			heap.Push(h, minStored)
			return
		}
	}
	heap.Push(h, max)

	// since we know that Val is integer we can do radix sort here
	for exp := int64(1); max.Val/exp > 0; exp *= 10 {
		for i := 0; i < len(elems); i++ {
			count[(elems[i].Val/exp)%10]++
		}

		for i := 1; i < len(count); i++ {
			count[i] += count[i-1]
		}

		for i := len(elems) - 1; i >= 0; i-- {
			digit := (elems[i].Val / exp) % 10
			sortBuf[count[digit]-1] = elems[i]
			count[digit]--
		}

		copy(elems, sortBuf)
		copy(count, zeroes)
	}

	// take n largest elems from batch and populate the heap
	for i := 1; i < n; i++ {
		idx := len(elems) - 1 - i
		if h.Len() == n {
			minStored := heap.Pop(h).(BatchElem)
			if minStored.Val > elems[idx].Val {
				// once we reach an elem in sorted batch that is smaller than the smallest stored we quit
				heap.Push(h, minStored)
				return
			}
		}
		heap.Push(h, elems[idx])
	}
}

// ElemIdx represents the position of row inside the input file
type ElemIdx struct {
	BatchNumber int64
	RowNumber   int
}

// BatchElem represents the row long value + position inside the input file
type BatchElem struct {
	ElemIdx
	Val int64
}

type BatchElemHeap []BatchElem

func (h BatchElemHeap) Len() int           { return len(h) }
func (h BatchElemHeap) Less(i, j int) bool { return h[i].Val < h[j].Val }
func (h BatchElemHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *BatchElemHeap) Push(x any) {
	*h = append(*h, x.(BatchElem))
}

func (h *BatchElemHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1 : 10]
	return x
}

var spaceByte = []byte{' '}

func splitFn(data []byte, atEOF bool) (advance int, token []byte, err error) {
	advance, token, err = bufio.ScanLines(data, atEOF)
	return advance, token[bytes.Index(token, spaceByte)+1:], err
}
