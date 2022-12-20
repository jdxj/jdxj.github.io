---
title: "聊一个 string 和 []byte 转换问题"
date: 2022-12-20T14:18:18+08:00
draft: false
tags:
  - go
---

[原文](https://huoding.com/2021/10/14/964)

错误的使用 uintptr, 转换结果可能被 GC.
