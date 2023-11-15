---
title: "Go测试编译为二进制执行"
date: 2023-11-15T15:47:15+08:00
tags:
  - go
---

[原文](http://www.hjwblog.com/archives/gotestrunregex)

```bash
$ go test -c pkgName
$ ./pkgName.test -test.run "TestXxx$" -test.v
```
