package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
)

const defaultN = 10
const defaultBatchSize = 10000

func main() {
	n := flag.Int("n", defaultN, "sets the amount of largest results")
	batchSize := flag.Int("b", defaultBatchSize, "sets the size of batch")

	flag.Parse()

	var filepath string
	fmt.Fscanln(os.Stdin, &filepath)

	f, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("unable to open file: %v", err)
	}

	indexies, err := largestN(f, *batchSize, *n)
	if err != nil {
		log.Fatalf("unable to get largest N results: %v", err)
	}

	// finding and outputting the rows from the initial file
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
		if cursour.RowNumber == *batchSize-1 {
			cursour.RowNumber = -1
			cursour.BatchNumber++
		}
		cursour.RowNumber++
	}
}
