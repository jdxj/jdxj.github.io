---
title: "stress"
date: 2023-04-20T09:55:35+08:00
summary: Linux系统压力测试工具
tags:
  - optimize
---

```bash
# 模拟一个 CPU 使用率 100% 的场景
$ stress --cpu 1 --timeout 600
# 模拟 I/O 压力
$ stress -i 1 --timeout 600
# 模拟的是 8 个进程
$ stress -c 8 --timeout 600
```