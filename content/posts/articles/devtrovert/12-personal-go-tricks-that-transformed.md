---
title: "12 Personal Go Tricks That Transformed My Productivity"
date: 2023-11-29T20:28:34+08:00
tags:
  - go
---

[原文](https://blog.devtrovert.com/p/12-personal-go-tricks-that-transformed)

感觉有用的

slice转array

```go
// go 1.20
func main() {
   a := []int{0, 1, 2, 3, 4, 5}
   b := [3]int(a[0:3])

  fmt.Println(b) // [0 1 2]
}
```

编译时检查接口实现

```go
var _ Buffer = (*StringBuffer)(nil)

type Buffer interface {
  Write(p []byte) (n int, err error)
}

type StringBuffer struct{}

func (s *StringBuffer) Writeee(p []byte) (n int, err error) {
  return 0, nil
}
```

类三元操作符

```go
// our utility
func Ter[T any](cond bool, a, b T) T {
  if cond {
    return a
  }
  
  return b
}

func main() {
  fmt.Println(Ter(true, 1, 2)) // 1 
  fmt.Println(Ter(false, 1, 2)) // 2
}
```

验证interface真为nil

```go
func IsNil(x interface{}) bool {
  if x == nil {
    return true
  }

  return reflect.ValueOf(x).IsNil()
}
```
