package main

import "fmt"

func main() {
	var a, b, c int8 = 127, 1, 0
	c = a + b // -128
	fmt.Println(c)

	var d, e, f uint8 = 255, 1, 0
	f = d + e // 0
	fmt.Println(f)
}
