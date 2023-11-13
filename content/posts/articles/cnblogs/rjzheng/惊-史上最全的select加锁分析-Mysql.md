---
title: "惊! 史上最全的select加锁分析(Mysql)"
date: 2023-02-06T16:49:26+08:00
summary: 该文章展示了事务隔离级别与锁的关系, 推荐阅读
tags:
  - mysql
  - lock
---

[原文](https://www.cnblogs.com/rjzheng/p/9950951.html)

# 锁的种类

- Record Locks
- Gap Locks, RR及以上级别才会加上
- Next-Key Locks

在RR, Serializable级别时, 在索引上的查询将锁表, 实现方式是Record+Gap Locks(Next-Key Locks)

# 加锁分析

- RC/RU+条件列非索引
- RC/RU+条件列是聚簇索引
- RC/RU+条件列是非聚簇索引
- RR/Serializable+条件列非索引
- RR/Serializable+条件列是聚簇索引
- RR/Serializable+条件列是非聚簇索引

# 总结

影响锁住的记录的范围因素

- 查询条件下的记录范围(等值, 范围)
- 索引类型(非索引, 聚簇索引, 非唯一索引)
- 隔离级别(Gap)
- 读/写锁(s, x)
