---
title: "sysbench"
date: 2023-04-20T11:31:53+08:00
summary: 一个多线程的基准测试工具，一般用来评估不同系统参数下的数据库负载情况
tags:
  - optimize
---

```bash
# 以10个线程运行5分钟的基准测试，模拟多线程切换的问题
$ sysbench --threads=10 --max-time=300 threads run
```

