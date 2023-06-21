---
title: "DMA零拷贝技术"
date: 2023-06-21T00:18:03+08:00
tags:
  - zero-copy
---

[原文](https://zhuanlan.zhihu.com/p/377237946)

你可以在你的 Linux 系统通过下面这个命令，查看网卡是否支持 scatter-gather 特性：
```bash
$ ethtool -k eth0 | grep scatter-gather scatter-gather: on
```