---
title: "聊一个 string 和 []byte 转换问题"
date: 2022-12-20T14:18:18+08:00
tags:
  - go
---

错误的使用uintptr转换结果可能被 GC.

[原文](https://huoding.com/2021/10/14/964)
