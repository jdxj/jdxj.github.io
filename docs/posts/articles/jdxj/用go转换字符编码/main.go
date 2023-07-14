package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func main() {
	str := "中国人"
	r := strings.NewReader(str)
	// r: utf8
	t := transform.NewReader(r, simplifiedchinese.GB18030.NewEncoder())

	b, err := io.ReadAll(t)
	if err != nil {
		panic(err)
	}
	fmt.Print("gb18030编码: ")
	for _, v := range b {
		fmt.Printf("%X ", v)
	}
	fmt.Println()

	r2 := bytes.NewReader(b)
	// r2: gb18030
	t = transform.NewReader(r2, simplifiedchinese.GB18030.NewDecoder())
	b, err = io.ReadAll(t)
	if err != nil {
		panic(err)
	}
	fmt.Print("gb18030解码(utf8编码):")
	for _, v := range b {
		fmt.Printf("%X ", v)
	}
}
