---
title: "execsnoop"
date: 2023-09-24T10:56:53+08:00
summary: 是一个专为短时进程设计的工具。它通过 ftrace 实时监控进程的 exec() 行为，并输出短时进程的基本信息，包括进程 PID、父进程 PID、命令行参数以及执行的结果。
tags:
  - optimize
---

```bash
$ execsnoop
```