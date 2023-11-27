---
title: "Go map 竟然也会发生内存泄漏？"
date: 2023-11-27T17:00:30+08:00
tags:
  - go
---

[原文](https://www.cnblogs.com/qcrao-2018/p/16885760.html)

map本身也占用内存, 即使用 delete 删除元素后也不释放, 且在 value <= 128B 时,
map使用原地存储会占用大量内存不释放, 造成内存泄漏.
