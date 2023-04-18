---
title: "Linux内核参数优化及原理"
date: 2023-04-08T10:06:49+08:00
tags:
  - tcp
---

[原文](https://code2life.top/2020/01/22/0036-linux-kernel-param/)

修改内核参数的方法

1. 编辑`/etc/sysctl.conf`添加配置
2. 执行`sysctl -p`立即生效