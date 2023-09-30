---
title: "free"
date: 2023-05-14T10:03:18+08:00
tags:
  - optimize
---

```bash
# 注意不同版本的free输出可能会有所不同
$ free
              total        used        free      shared  buff/cache   available
Mem:        8169348      263524     6875352         668     1030472     7611064
Swap:             0           0           0
```

- total 是总内存大小
- used 是已使用内存的大小，包含了共享内存
- free 是未使用内存的大小
- shared 是共享内存的大小
- buff/cache 是缓存和缓冲区的大小
- available 是新进程可用内存的大小

available 不仅包含未使用内存，还包括了可回收的缓存，所以一般会比未使用内存更大。不过，并不是所有缓存都可以回收，因为有些缓存可能正在使用中。

