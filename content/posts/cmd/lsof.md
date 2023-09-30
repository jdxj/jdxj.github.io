---
title: "lsof"
date: 2023-09-30T14:04:57+08:00
summary: 用来查看进程打开文件列表
tags:
  - optimize
---

“文件”不只有普通文件，还包括了目录、块设备、动态库、网络套接字等。

```bash
$ lsof -p 18940 
COMMAND   PID USER   FD   TYPE DEVICE  SIZE/OFF    NODE NAME 
python  18940 root  cwd    DIR   0,50      4096 1549389 / 
python  18940 root  rtd    DIR   0,50      4096 1549389 / 
… 
python  18940 root    2u   CHR  136,0       0t0       3 /dev/pts/0 
python  18940 root    3w   REG    8,1 117944320     303 /tmp/logtest.txt 
```

- FD 表示文件描述符号，TYPE 表示文件类型，NAME 表示文件路径。
- 这个进程打开了文件 /tmp/logtest.txt，并且它的文件描述符是 3 号，而 3 后面的 w ，表示以写的方式打开。
