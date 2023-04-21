---
title: "Ab"
date: 2023-04-21T18:33:33+08:00
---

ab（apache bench）是一个常用的 HTTP 服务性能测试工具

安装

```bash
$ apt install apache2-utils
```

使用

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
