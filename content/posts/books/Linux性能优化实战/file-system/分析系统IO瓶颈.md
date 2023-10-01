---
title: "分析系统I/O瓶颈"
date: 2023-10-01T13:34:30+08:00
---

# I/O性能指标

文件系统 I/O 性能指标

- 存储空间的使用情况，包括容量、使用量以及剩余空间
- 索引节点的使用情况，它也包括容量、使用量以及剩余量等三个指标
- 缓存使用情况，包括页缓存、目录项缓存、索引节点缓存以及各个具体文件系统（如 ext4、XFS 等）的缓存。
- 文件 I/O, IOPS（包括 r/s 和 w/s）、响应时间（延迟）以及吞吐量（B/s）等。

磁盘 I/O 性能指标

- 使用率，是指磁盘忙处理 I/O 请求的百分比。过高的使用率（比如超过 60%）通常意味着磁盘 I/O 存在性能瓶颈。
- IOPS（Input/Output Per Second），是指每秒的 I/O 请求数。
- 吞吐量，是指每秒的 I/O 请求大小。
- 响应时间，是指从发出 I/O 请求到收到响应的间隔时间。
- free命令中的buffer分析

![](https://static001.geekbang.org/resource/image/b6/20/b6d67150e471e1340a6f3c3dc3ba0120.png?wh=2650*808)

# 根据指标找工具

![](https://static001.geekbang.org/resource/image/6f/98/6f26fa18a73458764fcda00212006698.png?wh=1705*1901)

# 根据工具查指标

![](https://static001.geekbang.org/resource/image/ee/f3/ee11664d015f034e4042b9fa4fyycff3.jpg?wh=910*1336)

# 迅速分析 I/O 的性能瓶颈

1. 先用 iostat 发现磁盘 I/O 性能瓶颈；
2. 再借助 pidstat ，定位出导致瓶颈的进程；
3. 随后分析进程的 I/O 行为；
4. 最后，结合应用程序的原理，分析这些 I/O 的来源

![](https://static001.geekbang.org/resource/image/18/8a/1802a35475ee2755fb45aec55ed2d98a.png?wh=3732*1886)
