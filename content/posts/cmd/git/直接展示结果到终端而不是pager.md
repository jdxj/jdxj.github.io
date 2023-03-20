---
title: "直接展示结果到终端而不是pager"
date: 2023-03-20T16:18:51+08:00
summary: 要不还得按`q`退出
tags:
  - git config
---

```shell
$ git config --global pager.branch false
```

执行`git branch`的结果将直接输出, 而不是输出到类似 more 命令的界面.

可以利用自动补全查看`pager.*`下其他可配置东西.
