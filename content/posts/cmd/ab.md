---
title: "ab"
date: 2023-04-21T18:33:33+08:00
summary: apache bench是一个常用的 HTTP 服务性能测试工具
tags:
  - optimize
---

```bash
# 并发10个请求测试Nginx性能，总共测试100个请求
$ ab -c 10 -n 100 http://192.168.0.10:10000/
This is ApacheBench, Version 2.3 <$Revision: 1706008 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, 
...
Requests per second:    11.63 [#/sec] (mean)
Time per request:       859.942 [ms] (mean)
...
```

```bash
# 测试的并发请求数改成 5，同时把请求时长设置为 10 分钟
$ ab -c 5 -t 600 http://192.168.0.10:10000/
```