---
title: "过多的使用String键会造成内存浪费"
date: 2023-06-04T19:45:35+08:00
tags:
  - redis
---

占用内存的元数据过多

1. 全局hash表一个项是dictEntry

![](https://static001.geekbang.org/resource/image/b6/e7/b6cbc5161388fdf4c9b49f3802ef53e7.jpg?wh=2219*1371)

2. 封装底层编码的RedisObject

![](https://static001.geekbang.org/resource/image/34/57/3409948e9d3e8aa5cd7cafb9b66c2857.jpg?wh=2214*1656)

如果想节约内存, 可以考虑使用底层编码为ziplist的集合
