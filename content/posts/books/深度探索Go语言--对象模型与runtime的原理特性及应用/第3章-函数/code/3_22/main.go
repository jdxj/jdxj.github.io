// 第3章 code_3_22.go
package main

func main() {
	a := mc(2)
	a()
}

func mc(n int) func() int {
	return func() int {
		return n
	}
}
