---
title: "netstat"
date: 2023-04-24T10:51:20+08:00
---

```bash
# 列出listen/非listen套接字
$ netstat -a
# 列出tcp套接字
$ netstat -at
# 展示端口号
$ netstat -ltn
# 展示进程
$ netstat -ltnp

# 显示 8080 端口所有处于 ESTABLISHED 状态的连接
$ netstat -atnp | grep ":8080" | grep ESTABLISHED
```