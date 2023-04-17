---
title: "Go语言高性能编程手册"
date: 2023-04-17T14:54:26+08:00
tags:
  - go
---

[原文](https://mp.weixin.qq.com/s?__biz=Mzg5Mjc3MjIyMA==&mid=2247544709&idx=1&sn=83390ee723b35dc4e5dda3cec91ffdab&source=41#wechat_redirect)

- 少使用反射
  - 优先使用 strconv 而不是 fmt
  - 少量的重复不比反射差
  - 慎用 binary.Read 和 binary.Write
- 避免重复的字符串到字节切片的转换
- 指定容器容量
- 行内拼接字符串推荐使用运算符`+`
- 非行内拼接字符串推荐使用 strings.Builder
  - strings.Builder.Grow()
- 遍历 []struct{} 使用下标而不是 range
- 使用空结构体节省内存
  - `type Set map[string]struct{}`
  - `ch := make(chan struct{})`
  - `type Door struct{}`
- struct 布局要考虑内存对齐
  - 将字段宽度从小到大由上到下排列，来减少内存的占用
  - 当 struct{} 或空 array 作为结构体最后一个字段时，需要内存对齐
- 减少逃逸，将变量限制在栈上
  - 小的拷贝好过引用
  - 一般情况下，对于需要修改原对象值，或占用内存比较大的结构体，选择返回指针。对于只读的占用内存较小的结构体，直接返回值能够获得更好的性能。
  - 如果变量类型不确定，那么将会逃逸到堆上
- sync.Pool 复用对象
- **并发情况采用无锁设计**
  - CAS
  - 串行无锁
  - 分片减少锁竞争
  - 优先使用共享锁而非互斥锁
- 限制协程数量
  - 开销
    - 内存开销
    - 调度开销
    - GC开销
  - 池化
- 使用 sync.Once 避免重复执行
- 使用 sync.Cond 通知协程
