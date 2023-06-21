---
title: "TCP分手优化"
date: 2023-06-18T16:17:20+08:00
tags:
  - tcp
---

net.ipv4.tcp_tw_reuse = 1

- 开启后, 作为客户端时新连接可以使用仍然处于TIME-WAIT状态的端口
- 由于timestamp的存在, 操作系统可以拒绝迟到的报文, net.ipv4.tcp_timestamps = 1

net.ipv4.tcp_tw_recycle = 0

- 开启后, 同时作为客户端和服务器都可以使用TIME-WAIT状态的端口
- 不安全, 无法避免报文延迟, 重复等给新连接造成混乱

net.ipv4.tcp_max_tw_buckets = 262144

- time_wait状态连接的最大数量
- 超出后直接关闭连接

close 函数会让连接变为孤儿连接，shutdown 函数则允许在半关闭的连接上长时间传输数据。

# 四次挥手的流程

为什么建立连接是三次握手，而关闭连接需要四次挥手呢？

- TCP 不允许连接处于半打开状态时就单向传输数据，所以在三次握手建立连接时，服务器会把 ACK 和 SYN 放在一起发给客户端
- 其中，ACK 用来打开客户端的发送通道，SYN 用来打开服务器的发送通道。这样，原本的四次握手就降为三次握手了。

![](https://static001.geekbang.org/resource/image/74/51/74ac4e70ef719f19270c08201fb53a51.png?wh=943*613)

当连接处于半关闭状态时，TCP 是允许单向传输数据的。

![](https://static001.geekbang.org/resource/image/e2/b7/e2ef1347b3b4590da431dc236d9239b7.png?wh=1162*825)

# 主动方的优化

安全关闭连接的方式必须通过四次挥手，它由进程调用 close 或者 shutdown 函数发起，这二者都会向对方发送 FIN 报文（shutdown 参数须传入 SHUT_WR
或者 SHUT_RDWR 才会发送 FIN），区别在于 close 调用后，哪怕对方在半关闭状态下发送的数据到达主动方，进程也无法接收。(孤儿连接)

主动方发送 FIN 报文后，连接就处于 FIN_WAIT1 状态下，该状态通常应在数十毫秒内转为 FIN_WAIT2。只有迟迟收不到对方返回的 ACK 时，才能用
netstat 命令观察到 FIN_WAIT1 状态。此时，内核会定时重发 FIN 报文，其中重发次数由 tcp_orphan_retries 参数控制（注意，orphan 虽然是孤儿的
意思，该参数却不只对孤儿连接有效，事实上，它对所有 FIN_WAIT1 状态下的连接都有效），默认值是 0，特指 8 次：

```
net.ipv4.tcp_orphan_retries = 0
```

如果 FIN_WAIT1 状态连接有很多，你就需要考虑降低 tcp_orphan_retries 的值。当重试次数达到 tcp_orphan_retries 时，连接就会直接关闭掉。

对于正常情况来说，调低 tcp_orphan_retries 已经够用，但如果遇到恶意攻击，FIN 报文根本无法发送出去。这是由 TCP 的 2 个特性导致的。

- 首先，TCP 必须保证报文是有序发送的，FIN 报文也不例外，当发送缓冲区还有数据没发送时，FIN 报文也不能提前发送。
- 其次，TCP 有流控功能，当接收方将接收窗口设为 0 时，发送方就不能再发送数据。所以，当攻击者下载大文件时，就可以通过将接收窗口设为 0，导致 FIN
  报文无法发送，进而导致连接一直处于 FIN_WAIT1 状态。

tcp_max_orphans 定义了孤儿连接的最大数量。当进程调用 close 函数关闭连接后，无论该连接是在 FIN_WAIT1 状态，还是确实关闭了，这个连接都与该进
程无关了，它变成了孤儿连接。Linux 系统为防止孤儿连接过多，导致系统资源长期被占用，就提供了 tcp_max_orphans 参数。如果孤儿连接数量大于它，新增
的孤儿连接将不再走四次挥手，而是直接发送 RST 复位报文强制关闭。

当连接收到 ACK 进入 FIN_WAIT2 状态后，就表示主动方的发送通道已经关闭，接下来将等待对方发送 FIN 报文，关闭对方的发送通道。这时，如果连接是用
shutdown 函数关闭的，连接可以一直处于 FIN_WAIT2 状态。但对于 close 函数关闭的孤儿连接，这个状态不可以持续太久，而 tcp_fin_timeout 控制了
这个状态下连接的持续时长。

```
net.ipv4.tcp_fin_timeout = 60
```

保留 TIME_WAIT 状态，就可以应付重发的 FIN 报文，当然，其他数据报文也有可能重发，所以 TIME_WAIT 状态还能避免数据错乱。

为什么是 2 MSL 的时长呢？这其实是相当于至少允许报文丢失一次。比如，若 ACK 在一个 MSL 内丢失，这样被动方重发的 FIN 会在第 2 个 MSL 内到达，
TIME_WAIT 状态的连接可以应对。为什么不是 4 或者 8 MSL 的时长呢？你可以想象一个丢包率达到百分之一的糟糕网络，连续两次丢包的概率只有万分之一，这
个概率实在是太小了，忽略它比解决它更具性价比。

TIME_WAIT 和 FIN_WAIT2 状态的最大时长都是 2 MSL，由于在 Linux 系统中，MSL 的值固定为 30 秒，所以它们都是 60 秒。

Linux 提供了 tcp_max_tw_buckets 参数，当 TIME_WAIT 的连接数量超过该参数时，新关闭的连接就不再经历 TIME_WAIT 而直接关闭。

```
net.ipv4.tcp_max_tw_buckets = 5000
```

果服务器会主动向上游服务器发起连接的话，就可以把 tcp_tw_reuse 参数设置为 1，它允许作为客户端的新连接，在安全条件下使用 TIME_WAIT 状态下的端口。

```
net.ipv4.tcp_tw_reuse = 1
```

要想使 tcp_tw_reuse 生效，还得把 timestamps 参数设置为 1，它满足安全复用的先决条件（对方也要打开 tcp_timestamps ）：

```
net.ipv4.tcp_timestamps = 1
```

老版本的 Linux 还提供了 tcp_tw_recycle 参数，它并不要求 TIME_WAIT 状态存在 60 秒，很容易导致数据错乱，不建议设置为 1。

- 所以在 Linux 4.12 版本后，直接取消了这一参数。

```
net.ipv4.tcp_tw_recycle = 0
```

# 被动方的优化

内核没有权力替代进程去关闭连接，因为若主动方是通过 shutdown 关闭连接，那么它就是想在半关闭连接上接收数据。因此，Linux 并没有限制 CLOSE_WAIT
状态的持续时间。

由于 CLOSE_WAIT 状态下，连接已经处于半关闭状态，所以此时进程若要关闭连接，只能调用 close 函数（再调用 shutdown 关闭单向通道就没有意义了），
内核就会发出 FIN 报文关闭发送通道，同时连接进入 LAST_ACK 状态，等待主动方返回 ACK 来确认连接关闭。

如果被动方迅速调用 close 函数，那么被动方的 ACK 和 FIN 有可能在一个报文中发送，这样看起来，四次挥手会变成三次挥手，这只是一种特殊情况，不用在意。

# 连接双方同时关闭连接

![](https://static001.geekbang.org/resource/image/04/52/043752a3957d36f4e3c82cd83d472452.png?wh=1165*585)

双方在等待 ACK 报文的过程中，都等来了 FIN 报文。这是一种新情况，所以连接会进入一种叫做 CLOSING 的新状态，它替代了 FIN_WAIT2 状态。此时，内核
回复 ACK 确认对方发送通道的关闭，仅己方的 FIN 报文对应的 ACK 还没有收到。所以，CLOSING 状态与 LAST_ACK 状态下的连接很相似，它会在适时重发
FIN 报文的情况下最终关闭。
