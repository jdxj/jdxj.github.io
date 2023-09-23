---
title: "cat"
date: 2023-04-24T14:56:23+08:00
tags:
  - optimize
---

```bash
# 查看中断情况
# -d 参数表示高亮显示变化的区域
$ watch -d cat /proc/interrupts
           CPU0       CPU1
...
RES:    2450431    5279697   Rescheduling interrupts
...
```