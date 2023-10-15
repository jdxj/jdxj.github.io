---
title: "pprof"
date: 2023-10-14T13:33:09+08:00
tags:
  - go
  - pprof
  - deadlock
---

- [介绍](https://golang2.eddycjy.com/posts/ch6/01-pprof-1/)
- [排查内存泄漏的示例](https://blog.csdn.net/pengpengzhou/article/details/107000659)
- [goroutine泄漏的示例](https://blog.csdn.net/pengpengzhou/article/details/106946013)
- [CPU占用情况的示例](https://blog.csdn.net/pengpengzhou/article/details/107023331)
- [Profiling Go Programs](https://go.dev/blog/pprof)
- [pprof](https://github.com/google/pprof)
- [go pprof 性能分析](https://juejin.cn/post/6844903588565630983#heading-8)
- [golang pprof 实战](https://blog.wolfogre.com/posts/go-ppof-practice/)
- [实战Go内存泄露](https://segmentfault.com/a/1190000019222661#item-4)
- [golang 死锁检测](https://juejin.cn/s/golang%20%E6%AD%BB%E9%94%81%E6%A3%80%E6%B5%8B)

一些命令

```bash
$ go tool pprof http://localhost:6060/debug/pprof/profile
# pprof 内部命令
(pprof) top
(pprof) list funcName
(pprof) web # 生成svg
```