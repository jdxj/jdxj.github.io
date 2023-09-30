---
title: "cachestat"
date: 2023-09-29T10:43:59+08:00
summary: 提供了整个操作系统缓存的读写命中情况
tags:
  - optimize
---

```bash
# 以 1 秒的时间间隔，输出了 3 组缓存统计数据
$ cachestat 1 3
   TOTAL   MISSES     HITS  DIRTIES   BUFFERS_MB  CACHED_MB
       2        0        2        1           17        279
       2        0        2        1           17        279
       2        0        2        1           17        279 
```

- TOTAL ，表示总的 I/O 次数；
- MISSES ，表示缓存未命中的次数；
- HITS ，表示缓存命中的次数；
- DIRTIES， 表示新增到缓存中的脏页数；
- BUFFERS_MB 表示 Buffers 的大小，以 MB 为单位；
- CACHED_MB 表示 Cache 的大小，以 MB 为单位。
