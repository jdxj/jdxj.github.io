---
title: "buffer和cache的区别"
date: 2023-09-29T09:56:30+08:00
---

man free

- Buffers 是内核缓冲区用到的内存，对应的是 /proc/meminfo 中的 Buffers 值。
- Cache 是内核页缓存和 Slab 用到的内存，对应的是 /proc/meminfo 中的 Cached 与 SReclaimable 之和。

```bash
$ cat /proc/meminfo | grep -E "SReclaimable|Cached|Buffers"
Buffers:          818084 kB
Cached:          7201820 kB
SwapCached:            0 kB
SReclaimable:     939332 kB
```

man proc

- Buffers 是对原始磁盘块的临时存储，也就是用来缓存磁盘的数据，通常不会特别大（20MB 左右）。这样，内核就可以把分散的写集中起来，统一优化磁
  盘的写入，比如可以把多次小的写合并成单次大的写等等。
- Cached 是从磁盘读取文件的页缓存，也就是用来缓存从文件读取的数据。这样，下次访问这些文件数据时，就可以直接从内存中快速获取，而不需要再次访
  问缓慢的磁盘。SReclaimable 是 Slab 的一部分。Slab 包括两部分，其中的可回收部分，用 SReclaimable 记录；而不可回收部分，用
  SUnreclaim 记录。

```text
Buffers %lu
    Relatively temporary storage for raw disk blocks that shouldn't get tremendously large (20MB or so).

Cached %lu
   In-memory cache for files read from the disk (the page cache).  Doesn't include SwapCached.
...
SReclaimable %lu (since Linux 2.6.19)
    Part of Slab, that might be reclaimed, such as caches.
    
SUnreclaim %lu (since Linux 2.6.19)
    Part of Slab, that cannot be reclaimed on memory pressure.
```

写文件时会用到 Cache 缓存数据，而写磁盘则会用到 Buffer 来缓存数据。

读文件时数据会缓存到 Cache 中，而读磁盘时数据会缓存到 Buffer 中。

# 结论

Buffer 是对磁盘数据的缓存，而 Cache 是文件数据的缓存，它们既会用在读请求中，也会用在写请求中。