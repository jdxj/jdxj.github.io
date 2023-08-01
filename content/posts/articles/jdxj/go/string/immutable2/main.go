package main

import (
	"fmt"
	"unsafe"
)

// chapter3/sources/string_immutable2.go

func main() {
	// 原始string
	var s string = "hello"
	fmt.Println("original string:", s)

	// 试图通过unsafe指针改变原始string
	modifyString(&s)
	fmt.Println(s)
}

func modifyString(s *string) {
	// 取出第一个8字节的值
	// p是stringStruct的地址
	p := (*uintptr)(unsafe.Pointer(s))

	// 获取底层数组的地址
	var array *[5]byte = (*[5]byte)(unsafe.Pointer(*p))

	var len *int = (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(s)) + unsafe.Sizeof((*uintptr)(nil))))

	for i := 0; i < (*len); i++ {
		fmt.Printf("%p => %c\n", &((*array)[i]), (*array)[i])
		p1 := &((*array)[i])
		v := (*p1)
		(*p1) = v + 1 // try to change the character
	}
}

/* 输出
original string: hello
0x49b7d3 => h
unexpected fault address 0x49b7d3
fatal error: fault
[signal SIGSEGV: segmentation violation code=0x2 addr=0x49b7d3 pc=0x483417]
*/
