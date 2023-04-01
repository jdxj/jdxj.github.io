package main

import "fmt"

func main() {
	var (
		seq1 uint8 = 255
		seq2 uint8 = 1
	)
	fmt.Println("case1(未回绕):")
	fmt.Printf("befort: %t\n\n", before(seq1, seq2))

	seq1 = 255
	seq2 = 128
	fmt.Println("case2(已回绕):")
	fmt.Printf("befort: %t\n", before(seq1, seq2))
}

func before(seq1, seq2 uint8) bool {
	fmt.Printf("seq1: %d, signed: %d\n", seq1, int8(seq1))
	fmt.Printf("seq2: %d, signed: %d\n", seq2, int8(seq2))
	return int8(seq1-seq2) < 0
}
