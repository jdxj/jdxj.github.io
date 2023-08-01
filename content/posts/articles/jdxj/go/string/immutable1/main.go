package main

import "fmt"

// chapter3/sources/string_immutable1.go
func main() {
	// 原始字符串
	var s string = "hello"
	fmt.Println("original string:", s)

	// 切片化后试图改变原字符串
	sl := []byte(s)
	sl[0] = 't'
	fmt.Println("slice:", string(sl))
	fmt.Println("after reslice, the original string is:", string(s))
}

/* 输出
original string: hello
slice: tello
after reslice, the original string is: hello
*/
