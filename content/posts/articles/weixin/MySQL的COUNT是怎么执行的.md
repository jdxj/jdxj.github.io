---
title: "MySQL的COUNT是怎么执行的"
date: 2023-02-20T17:28:18+08:00
tags:
  - mysql
---

COUNT函数用于统计在符合搜索条件的记录中，指定的表达式expr不为NULL的行数有多少。

对于`COUNT(*)`、`COUNT(常数)`、`COUNT(主键)`形式的COUNT函数来说，优化器可以选择最小的索引执行查询，从而提升效率，它们的执行过程是一样的，只不
过在判断表达式是否为NULL时选择不同的判断方式，这个判断为NULL的过程的代价可以忽略不计，所以我们可以认为`COUNT(*)`、`COUNT(常数)`、`COUNT(主键)`
所需要的代价是相同的。

而对于`COUNT(非主键列)`来说，server层必须要从InnoDB中读到包含非主键列的记录，所以优化器并不能随心所欲的选择最小的索引去执行。

[原文](https://mp.weixin.qq.com/s/_z7BnFlm4gEAOGTVWTcsWA)
