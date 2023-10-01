---
title: "ss"
date: 2023-10-01T14:35:06+08:00
summary: 查看套接字、网络栈、网络接口以及路由表的信息
tags:
  - optimize
---

```bash
# -l 表示只显示监听套接字
# -t 表示只显示 TCP 套接字
# -n 表示显示数字地址和端口(而不是名字)
# -p 表示显示进程信息
$ ss -ltnp | head -n 3
State    Recv-Q    Send-Q        Local Address:Port        Peer Address:Port
LISTEN   0         128           127.0.0.53%lo:53               0.0.0.0:*        users:(("systemd-resolve",pid=840,fd=13))
LISTEN   0         128                 0.0.0.0:22               0.0.0.0:*        users:(("sshd",pid=1459,fd=3))
```

接收队列（Recv-Q）和发送队列（Send-Q），它们通常应该是 0。当你发现它们不是 0 时，说明有网络包的堆积发生。

当套接字处于连接状态（Established）时

- Recv-Q 表示套接字缓冲还没有被应用程序取走的字节数（即接收队列长度）。
- Send-Q 表示还没有被远端主机确认的字节数（即发送队列长度）。

当套接字处于监听状态（Listening）时

- Recv-Q 表示全连接队列的长度。
- Send-Q 表示全连接队列的最大长度。

所谓全连接，是指服务器收到了客户端的 ACK，完成了 TCP 三次握手，然后就会把这个连接挪到全连接队列中。这些全连接中的套接字，还需要被 accept()
系统调用取走，服务器才可以开始真正处理客户端的请求。

所谓半连接是指还没有完成 TCP 三次握手的连接，连接只进行了一半。服务器收到了客户端的 SYN 包后，就会把这个连接放到半连接队列中，然后再向客户
端发送 SYN+ACK 包。

协议栈统计信息

```bash
$ ss -s
Total: 186 (kernel 1446)
TCP:   4 (estab 1, closed 0, orphaned 0, synrecv 0, timewait 0/0), ports 0

Transport Total     IP        IPv6
*    1446      -         -
RAW    2         1         1
UDP    2         2         0
TCP    4         3         1
...
```

