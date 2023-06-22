---
title: "随机数生成"
date: 2023-06-22T14:10:38+08:00
tags:
  - go
  - crypto
---

Go密码学包crypto/rand提供了密码学级别的随机数生成器实现rand.Reader，在不同平台上rand.Reader使用的数据源有所不同。在类Unix操作系统上，它使
用的是该平台上密码学应用的首选随机数源/dev/urandom

```go
// chapter9/sources/go-crypto/rand_generate.go
package main

import (
    "crypto/rand"
    "fmt"
)

func main() {
    c := 32
    b := make([]byte, c)
    _, err := rand.Read(b)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%x\n", b)
}
```
