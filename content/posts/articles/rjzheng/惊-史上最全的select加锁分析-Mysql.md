---
title: "惊! 史上最全的select加锁分析(Mysql)"
date: 2023-02-06T16:49:26+08:00
summary: 该文章展示了事务隔离级别与锁的关系, 推荐阅读
tags:
  - mysql
  - lock
---

[原文](https://www.cnblogs.com/rjzheng/p/9950951.html)

- RC/RU+条件列非索引
- RC/RU+条件列是聚簇索引
- RC/RU+条件列是非聚簇索引
- RR/Serializable+条件列非索引
- RR/Serializable+条件列是聚簇索引
- RR/Serializable+条件列是非聚簇索引
