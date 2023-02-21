// 第7章
package main

import (
	"math/rand"
	"time"
)

func main() {
	var data int
	var ok bool

	someValue := rand.Int()

	// 第7章 code_7_1
	go func() {
		for {
			if !ok {
				data = someValue
				ok = true
			}
		}
	}()

	var sum int

	// 第7章 code_7_2
	go func() {
		for {
			if ok {
				sum += data
				ok = false
			}
		}
	}()

	time.Sleep(time.Second * 10)
}
