---
title: "数组与切片"
date: 2023-04-19T12:53:53+08:00
---

# 数组与切片有什么异同

数组就是一片连续的内存， slice 实际上是一个结构体，包含三个字段：长度、容量、底层数组。

```go
// runtime/slice.go
type slice struct {
	array unsafe.Pointer // 元素指针
	len   int // 长度 
	cap   int // 容量
}
```

[3]int, [4]int是不同的类型

# 切片的容量是怎样增长的

append函数返回值是一个新的slice，Go编译器不允许调用了 append 函数后不使用返回值。

测试代码

{{< embedcode go "slice-grow_test.go" >}}

# 切片作为函数参数

Go 语言的函数参数传递，只有值传递，没有引用传递。
