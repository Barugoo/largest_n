package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

func main() {
	var filepath string
	fmt.Fscanln(os.Stdin, &filepath)

	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	batchSize := 10000
	indexies := largestN(f, batchSize, 10)

	f.Seek(0, 0)
	scanner := bufio.NewScanner(f)

	var cursour ElemIdx
	var foundCounter int
	for scanner.Scan() {

		if indexies[foundCounter] == cursour {
			url, _, _ := bytes.Cut(scanner.Bytes(), spaceByte)

			fmt.Printf("%s\n", url)
			foundCounter++
		}
		if cursour.RowNumber == batchSize-1 {
			cursour.RowNumber = -1
			cursour.BatchNumber++
		}
		cursour.RowNumber++
	}
}
