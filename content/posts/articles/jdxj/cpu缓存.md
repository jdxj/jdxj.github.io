---
title: "CPU缓存"
date: 2023-06-20T19:54:09+08:00
tags:
  - cpu
  - cache
---

三级缓存要比一、二级缓存大许多倍，这是因为当下的 CPU 都是多核心的，每个核心都有自己的一、二级缓存，但三级缓存却是一颗 CPU 上所有核心共享的。

![](https://static001.geekbang.org/resource/image/92/0c/9277d79155cd7f925c27f9c37e0b240c.jpg?wh=3749*2433)

CPU 会区别对待指令与数据, 要分开来看二者的缓存命中率

CPU Cache Line 定义了缓存一次载入数据的大小

```bash
# 查看缓存大小
$ cat /sys/devices/system/cpu/cpu0/cache/index0/size
# 查看cache line大小
$ cat /sys/devices/system/cpu/cpu0/cache/index1/coherency_line_size
```

提升数据缓存命中率

- 按照内存布局顺序访问将会带来很大的性能提升。
- 哈希表里桶的大小如 server_names_hash_bucket_size，它默认就等于 CPU Cache Line 的值, 可以尽量减少访问内存的次数

提升指令缓存的命中率

- 在循环中连续走同一分支

提升多核 CPU 下的缓存命中率

- 操作系统提供了将进程或者线程绑定到某一颗 CPU 上运行的能力 (避免切换到其他核心时导致的缓存不中问题)

# 结论

CPU 缓存分为数据缓存与指令缓存，对于数据缓存，我们应在循环体中尽量操作同一块内存上的数据，由于缓存是根据 CPU Cache Line 批量操作数据的，所以顺
序地操作连续内存数据时也有性能提升。

对于指令缓存，有规律的条件分支能够让 CPU 的分支预测发挥作用，进一步提升执行效率。对于多核系统，如果进程的缓存命中率非常高，则可以考虑绑定 CPU 来
提升缓存命中率。
