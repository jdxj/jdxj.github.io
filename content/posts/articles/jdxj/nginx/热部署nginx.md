---
title: "热部署nginx"
date: 2023-08-05T11:57:17+08:00
tags:
  - nginx
---

1. 备份原有nginx可执行文件, 并复制新可执行文件
2. 通知使用新文件

```bash
$ kill -USR2 旧pid
```

3. 通知优雅关闭旧worker

```bash
$ kill -WINCH 旧pid
```
