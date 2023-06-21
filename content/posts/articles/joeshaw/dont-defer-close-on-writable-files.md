---
title: "Don't defer Close() on writable files"
date: 2023-06-19T13:11:20+08:00
tags:
  - go
---

[原文](https://www.joeshaw.org/dont-defer-close-on-writable-files/)

Write()为了性能使用异步方式, 操作系统同步数据到磁盘的一个时机在Close(), 如果不捕获可能造成数据丢失.

安全的Close

```go
func helloNotes() error {
    f, err := os.Create("/home/joeshaw/notes.txt")
    if err != nil {
        return err
    }
    defer f.Close()

    if err = io.WriteString(f, "hello world"); err != nil {
        return err
    }

    return f.Sync()
}
```
