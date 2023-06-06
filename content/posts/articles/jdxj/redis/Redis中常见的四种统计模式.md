---
title: "Redis中常见的四种统计模式"
date: 2023-06-04T23:23:45+08:00
tags:
  - redis
---

# 聚合统计

集合的交, 差, 并集

- SUNIONSTORE
- SDIFFSTORE
- SINTERSTORE

**这些命令复杂度较高**

# 排序统计

可以使用sorted set保证分页时, 元素不会串

- ZRANGEBYSCORE

# 二值状态统计

bitmap

- SETBIT
- BITOP
- BITCOUNT

# 基数统计

指统计一个集合中不重复的元素个数

set

- SADD
- SCARD

hash

- HSET
- HLEN

hyperLogLog

- 用于统计基数的数据集合类型，它的最大优势就在于，当集合元素数量非常多时，它计算基数所需的空间总是固定的，而且还很小。
- 有一定误差
- PFADD
- PFCOUNT

![](https://static001.geekbang.org/resource/image/c0/6e/c0bb35d0d91a62ef4ca1bd939a9b136e.jpg?wh=2866*1739)
