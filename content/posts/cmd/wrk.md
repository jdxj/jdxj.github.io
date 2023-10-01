---
title: "wrk"
date: 2023-10-01T15:25:50+08:00
summary: HTTP 性能测试工具
tags:
  - optimize
---

https://github.com/wg/wrk

```bash
# -c表示并发连接数1000，-t表示线程数为2
$ wrk -c 1000 -t 2 http://192.168.0.30/
Running 10s test @ http://192.168.0.30/
  2 threads and 1000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    65.83ms  174.06ms   1.99s    95.85%
    Req/Sec     4.87k   628.73     6.78k    69.00%
  96954 requests in 10.06s, 78.59MB read
  Socket errors: connect 0, read 0, write 0, timeout 179
Requests/sec:   9641.31
Transfer/sec:      7.82MB
```

# 参考

- [wrk——轻量级异步性能测试工具](https://sq.sf.163.com/blog/article/200008406328934400)
