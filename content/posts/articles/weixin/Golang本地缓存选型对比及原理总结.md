---
title: "Golang本地缓存选型对比及原理总结"
date: 2023-06-25T20:09:23+08:00
tags:
  - go
  - cache
---

[原文](https://mp.weixin.qq.com/s?src=11&timestamp=1687694260&ver=4612&signature=qlI5-v11MvpO4HQaMeyRXmZm69zqrnWaKWVnT*QKsGl4VXi1pfOogVBqmvLaNp7dAoeLsiuMG2DU7b61D*py35Z8HFAR2x2JowXqDnHl9*iJZh1-*Ygqjv5xgiO1ZLAt&new=1)

# 比较的库

- freecache
- bigcache
- fastcache
- offheap
- groupcache
- ristretto
- go-cache

# 实现零GC的方案

- 无GC：分配堆外内存(Mmap)
- 避免GC：map非指针优化(map[uint64]uint32)或者采用slice实现一套无指针的map。
- 避免GC：数据存入[]byte slice(可考虑底层采用环形队列封装循环使用空间)

# 实现高性能的关键

- 数据分片(降低锁的粒度)
