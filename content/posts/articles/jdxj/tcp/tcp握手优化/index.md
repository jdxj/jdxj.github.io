---
title: "TCP握手优化"
date: 2023-06-17T14:48:40+08:00
tags:
  - tcp
---

![](./tcp连接队列.drawio.svg)

# 延长syn, accept队列

- 应用层connect超时时间
- 操作系统内核
  - 服务器端syn_rcv状态
    - net.ipv4.tcp_max_syn_backlog: syn_rcvd状态连接的最大个数
    - net.ipv4.tcp_synack_retries: 被动建立连接时, 发syn/ack的重试次数
  - 客户端syn_ent状态
    - net.ipv4.tcp_syn_retries = 6 主动建立连接时, 发syn的重试次数
    - net.ipv4.ip_local_port_range = 32768 60999 建立连接时的本地端口可用范围
  - accept队列设置

# Fast Open降低时延

[TFO]({{< ref "posts/articles/jdxj/tcp/tfo.md" >}})

Linux上打开TFO

net.ipv4.tcp_fastopen

- 0: 关闭
- 1: 作为客户端时可以使用TFO
- 2: 作为服务器时可以使用TFO
- 3: 无论作为客户端还是服务器, 都可以使用TFO

# TCP_DEFER_ACCEPT

服务端收到客户端发来的数据后才唤醒阻塞在accept()的应用

# 客户端的优化

三次握手建立连接的首要目的是同步序列号。只有同步了序列号才有可靠的传输，TCP 协议的许多特性都是依赖序列号实现的，比如流量控制、消息丢失后的重发等等

![](https://static001.geekbang.org/resource/image/c5/aa/c51d9f1604690ab1b69e7c4feb2f31aa.jpg?wh=2052*620)

客户端在等待服务器回复的 ACK 报文。正常情况下，服务器会在几毫秒内返回 ACK，但如果客户端迟迟没有收到 ACK 会怎么样呢？客户端会重发 SYN，重试的次
数由 tcp_syn_retries 参数控制，默认是 6 次：

```
net.ipv4.tcp_syn_retries = 6
```

第 1 次重试发生在 1 秒钟后，接着会以翻倍的方式在第 2、4、8、16、32 秒共做 6 次重试，最后一次重试会等待 64 秒，如果仍然没有返回 ACK，才会终止
三次握手。所以，总耗时是 1+2+4+8+16+32+64=127 秒，超过 2 分钟。

如果这是一台有明确任务的服务器，你可以根据网络的稳定性和目标服务器的繁忙程度修改重试次数，调整客户端的三次握手时间上限。比如内网中通讯时，就可以适
当调低重试次数，尽快把错误暴露给应用程序。

![](https://static001.geekbang.org/resource/image/a3/8f/a3c5e77a228478da2a6e707054043c8f.png?wh=943*613)

# 服务器端的优化

当服务器收到 SYN 报文后，服务器会立刻回复 SYN+ACK 报文，既确认了客户端的序列号，也把自己的序列号发给了对方。此时，服务器端出现了新连接，状态是
SYN_RCV（RCV 是 received 的缩写）。这个状态下，服务器必须建立一个 SYN 半连接队列来维护未完成的握手信息，当这个队列溢出后，服务器将无法再建立
新连接。

![](https://static001.geekbang.org/resource/image/c3/82/c361e672526ee5bb87d5f6b7ad169982.png?wh=690*304)

获得由于队列已满而引发的失败次数

```bash
# 是一个累积值
$ netstat -s | grep "SYNs to LISTEN"
    1192450 SYNs to LISTEN sockets dropped
```

如果数值在持续增加，则应该调大 SYN 半连接队列

```
net.ipv4.tcp_max_syn_backlog = 1024
```

开启 syncookies 功能就可以在不使用 SYN 队列的情况下成功建立连接

![](https://static001.geekbang.org/resource/image/0d/c0/0d963557347c149a6270d8102d83e0c0.png?wh=690*319)

修改 tcp_syncookies

- 0表示关闭该功能
- 2表示无条件开启功能
- 1则表示仅当 SYN 半连接队列放不下时，再启用它。

由于 syncookie 仅用于应对 SYN 泛洪攻击（攻击者恶意构造大量的 SYN 报文发送给服务器，造成 SYN 半连接队列溢出，导致正常客户端的连接无法建立），
这种方式建立的连接，许多 TCP 特性都无法使用。所以，**应当把tcp_syncookies设置为1**，仅在队列满时再启用。

```
net.ipv4.tcp_syncookies = 1
```

如果服务器没有收到 ACK，就会一直重发 SYN+ACK 报文。当网络繁忙、不稳定时，报文丢失就会变严重，此时应该调大重发次数。反之则可以调小重发次数。修改
重发次数的方法是，调整 tcp_synack_retries 参数：

```
net.ipv4.tcp_synack_retries = 5
```

服务器收到 ACK 后连接建立成功，此时，内核会把连接从 SYN 半连接队列中移出，再移入 accept 队列，等待进程调用 accept 函数时把连接取出来。如果进
程不能及时地调用 accept 函数，就会造成**accept 队列溢出，最终导致建立好的 TCP 连接被丢弃**。

丢弃连接只是 Linux 的默认行为，我们还可以选择向客户端发送 RST 复位报文，告诉客户端连接已经建立失败。打开这一功能需要将
tcp_abort_on_overflow 参数设置为 1。

```
net.ipv4.tcp_abort_on_overflow = 0
```

- 通常情况下，应当把 tcp_abort_on_overflow 设置为 0，因为这样更有利于应对突发流量。
- 只有你非常肯定 accept 队列会长期溢出时，才能设置为 1 以尽快通知客户端。

listen 函数的 backlog 参数就可以设置 accept 队列的大小。事实上，backlog 参数还受限于 Linux 系统级的队列长度上限，当然这个上限阈值也可以通
过 somaxconn 参数修改。

```
net.core.somaxconn = 128
```

# TFO 技术如何绕过三次握手？

它把通讯分为两个阶段

- 第一阶段为首次建立连接，这时走正常的三次握手，但在客户端的 SYN 报文会明确地告诉服务器它想使用 TFO 功能，这样服务器会把客户端 IP 地址用只有自
  己知道的密钥加密（比如 AES 加密算法），作为 Cookie 携带在返回的 SYN+ACK 报文中，客户端收到后会将 Cookie 缓存在本地。
- 之后，如果客户端再次向服务器建立连接，就可以在第一个 SYN 报文中携带请求数据，同时还要附带缓存的 Cookie。很显然，这种通讯方式下不能再采用经典
  的“先 connect 再 write 请求”这种编程方法，而要改用 sendto 或者 sendmsg 函数才能实现。

服务器收到后，会用自己的密钥验证 Cookie 是否合法，验证通过后连接才算建立成功，再把请求交给进程处理，同时给客户端返回 SYN+ACK。虽然客户端收到后
还会返回 ACK，但服务器不等收到 ACK 就可以发送 HTTP 响应了，这就减少了握手带来的 1 个 RTT 的时间消耗。

![](https://static001.geekbang.org/resource/image/7a/c3/7ac29766ba8515eea5bb331fce6dc2c3.png?wh=961*806)

由于只有客户端和服务器同时支持时，TFO 功能才能使用，所以 tcp_fastopen 参数是按比特位控制的。其中，第 1 个比特位为 1 时，表示作为客户端时支持
TFO；第 2 个比特位为 1 时，表示作为服务器时支持 TFO，所以当 tcp_fastopen 的值为 3 时（比特为 0x11）就表示完全支持 TFO 功能。

```
net.ipv4.tcp_fastopen = 3
```
