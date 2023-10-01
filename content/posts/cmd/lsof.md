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

```bash
$ lsof -p 27458
COMMAND  PID USER   FD   TYPE DEVICE SIZE/OFF NODE NAME
...
mysqld  27458      999   38u   REG    8,1 512440000 2601895 /var/lib/mysql/test/products.MYD
```

- 38 后面的 u 表示， mysqld 以读写的方式访问文件。

```bash
$ lsof -p 9085
redis-ser 9085 systemd-network    3r     FIFO   0,12      0t0 15447970 pipe
redis-ser 9085 systemd-network    4w     FIFO   0,12      0t0 15447970 pipe
redis-ser 9085 systemd-network    5u  a_inode   0,13        0    10179 [eventpoll]
redis-ser 9085 systemd-network    6u     sock    0,9      0t0 15447972 protocol: TCP
redis-ser 9085 systemd-network    7w      REG    8,1  8830146  2838532 /data/appendonly.aof
redis-ser 9085 systemd-network    8u     sock    0,9      0t0 15448709 protocol: TCP
```

描述符编号为 3 的是一个 pipe 管道，5 号是 eventpoll，7 号是一个普通文件，而 8 号是一个 TCP socket。

```bash
# -i表示显示网络套接字信息
$ lsof -i
COMMAND    PID            USER   FD   TYPE   DEVICE SIZE/OFF NODE NAME
redis-ser 9085 systemd-network    6u  IPv4 15447972      0t0  TCP localhost:6379 (LISTEN)
redis-ser 9085 systemd-network    8u  IPv4 15448709      0t0  TCP localhost:6379->localhost:32996 (ESTABLISHED)
python    9181            root    3u  IPv4 15448677      0t0  TCP *:http (LISTEN)
python    9181            root    5u  IPv4 15449632      0t0  TCP localhost:32996->localhost:6379 (ESTABLISHED)
```
