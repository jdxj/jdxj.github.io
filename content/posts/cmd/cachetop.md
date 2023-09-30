---
title: "cachetop"
date: 2023-09-29T10:46:52+08:00
summary: 提供了每个进程的缓存命中情况
tags:
  - optimize
---

/usr/share/bcc/tools

```bash
$ cachetop
11:58:50 Buffers MB: 258 / Cached MB: 347 / Sort: HITS / Order: ascending
PID      UID      CMD              HITS     MISSES   DIRTIES  READ_HIT%  WRITE_HIT%
   13029 root     python                  1        0        0     100.0%       0.0%
```

READ_HIT 和 WRITE_HIT ，分别表示读和写的缓存命中率

cachetop 工具并不把直接 I/O 算进来