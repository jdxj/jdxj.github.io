---
title: "用Go 请求接口不执行body.Close()会内存溢出吗？这次告诉你真相！"
date: 2023-11-17T16:16:11+08:00
tags:
  - go
---

[原文](https://mp.weixin.qq.com/s/MwnArLI04gy-_XzbSzAzsw)

- 既不执行 ioutil.ReadAll(resp.Body) 也不执行resp.Body.Close()，并且不设置http.Client内timeout的时候，就会导致协程泄露。
- SetDeadline是指 tcp 连接存活时长
