---
title: "stress"
date: 2023-04-20T09:55:35+08:00
---

stress 是一个 Linux 系统压力测试工具

```bash
$ apt install stress
```

使用

```bash
$ stress --cpu 1 --timeout 600
$ stress -i 1 --timeout 600
$ stress -c 8 --timeout 600
```