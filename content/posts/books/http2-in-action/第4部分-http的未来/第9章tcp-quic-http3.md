---
title: "第9章TCP, QUIC和HTTP/3"
date: 2023-03-29T20:56:14+08:00
---

## 9.1 TCP的低效率因素，以及HTTP

TCP运行的基本方式会导致5个主要的问题，它们至少影响到了HTTP：

- 有一个连接创建的延迟。要在连接开始时协商发送方和接收方可以使用的序列号。
- TCP慢启动算法限制了TCP的性能，它小心翼翼地处理发送的数据量，以尽可能防止重传。
- 不充分使用连接会导致限流阈值降低。如果连接未被充分使用，TCP会将拥塞窗口的大小减小，因为它不确定在上个最优的拥塞窗口之后网络参数有没有发生变化。
- 丢包也会导致TCP的限流阈值降低。TCP认为所有的丢包都是由窗口拥堵造成的，但其实并不是。
- 数据包可能被排队。乱序接收到的数据包会被排队，以保证数据是有序的。

### 9.1.1 创建HTTP连接的延迟

图9.1　HTTPS连接所需要的TCP和HTTPS设置

![](https://res.weread.qq.com/wrepub/epub_32517945_324)

### 9.1.2 TCP拥塞控制对性能的影响

拥塞控制算法和概念增加了稳定性，但也带来了低效的问题

**TCP慢启动**

表9.1　常见的TCP慢启动增长

![](https://res.weread.qq.com/wrepub/epub_32517945_327)

图9.4　TCP慢启动到最佳容量的过程

![](https://res.weread.qq.com/wrepub/epub_32517945_328)

当达到最大容量之后，如果没有发生丢包，TCP拥塞控制就进入拥塞避免阶段，随后拥塞窗口还会持续增长，但会变成慢得多的线性增长, 直到它开始看到丢包，认为
到了最大容量

TCP慢启动很慢吗

- 因为TCP慢启动的指数增长特性，所以按大多数定义来说它并不慢。

尽量把东西放到前14KB中

开始的10个TCP数据包（至少）会被如下消息使用：

- 两个HTTPS响应（Server Hello和Change Spec）
- 两个HTTP/2 SETTINGS帧（服务器发送一个，另外一个用来确认客户端的SETTINGS帧）
- 一个HEADERS帧，响应第一个请求

表9.2　通常使用6个连接的TCP慢启动增长

![](https://res.weread.qq.com/wrepub/epub_32517945_329)

**连接闲置降低性能**

在连接刚启动时和连接闲置时，TCP慢启动算法会导致延迟。TCP比较小心谨慎，在闲置一段时间后，网络情况可能发生变化，所以TCP将拥塞窗口大小降低，重新进行
慢启动流程，以再次找到最佳的拥塞窗口大小。

**丢包降低TCP性能**

tcp在遇到丢包时, 直接将拥塞窗口大小减半

图9.6　TCP拥塞窗口大小受丢包影响

![](https://res.weread.qq.com/wrepub/epub_32517945_331)

丢包带来的影响在HTTP/2中尤其严重，因为它只使用单个连接。在HTTP/2的世界中，一次丢包会导致所有的资源下载速度变慢。而HTTP/1.1可能有6个独立的连接，
一次丢包只会减慢其中一个连接，但是另外5个不受影响。

**丢包会导致数据排队**

图9.7　同时传输多个响应

![](https://res.weread.qq.com/wrepub/epub_32517945_333)

图9.8　TCP重传一个HTTP/2帧的一部分

![](https://res.weread.qq.com/wrepub/epub_32517945_334)

如果没有发生其他的丢包，流7和流9会在重传的数据到来之前被完整接收。但这些响应必须排队，因为TCP要保证顺序，所以尽管已经完整下载，script.js和
image.jpg还不能被使用。

图9.9　HTTP/1.1下TCP重传只影响需要重传的连接

![](https://res.weread.qq.com/wrepub/epub_32517945_335)

HTTP/2在HTTP层解决了队头阻塞（HOL）的问题，因为有多路复用，单个响应的延迟不会影响其他资源使用当前的HTTP连接。但是，在TCP层队头阻塞依然存在。
一个流的丢包会直接影响到其他所有的流，尽管它们可能不需要排队。

### 9.1.3 TCP低效率因素对HTTP/2的影响

### 9.1.4 优化TCP

在Linux中，大多数TCP设置在如下路径中

```
/proc/sys/net/ipv4
```

查看

```
cat /proc/sys/net/ipv4/tcp_slow_start_after_idle
```

设置

```
sysctl -w net.ipv4.tcp_slow_start_after_idle=0
```

提高初始拥塞窗口大小

- 这个设置值通常被写死到内核代码中，所以除非升级操作系统，否则不建议修改。

支持窗口缩放

在传统情况下，TCP所允许的最大拥塞窗口大小是65 535字节，但新版本中添加了缩放因子，理论上允许拥塞窗口最大到1GB。

```
cat /proc/sys/net/ipv4/tcp_window_scaling
```

使用SACK

- Selective Acknowledgment

```
/proc/sys/net/ipv4/tcp_sack
```

禁止重启慢启动

```
cat /proc/sys/net/ipv4/tcp_slow_start_after_idle
// 禁用
sysctl -w net.ipv4.tcp_slow_start_after_idle=0
```

使用TFO

- TCP Fast Open
- 出于安全的原因，这个数据包只能在TCP重连时使用，而不能在初次连接时使用，它同时需要客户端和服务端的支持。

图9.17　未使用TFO和使用了TFO的TCP和HTTPS重连握手

![](https://res.weread.qq.com/wrepub/epub_32517945_352)

```
cat /proc/sys/net/ipv4/tcp_fastopen
```

表9.5　TFO设置项的取值

![](https://res.weread.qq.com/wrepub/epub_32517945_354)

```
echo "3" > /proc/sys/net/ipv4/tcp_fastopen
```

**使用拥塞控制算法，PRR和BBR**

- BBR是Google发明的，从Linux内核4.9版开始可以启用它，它不需要客户端支持。

### 9.1.5 TCP和HTTP的未来

## 9.2 QUIC

创建QUIC时考虑到了以下特性

- 大量减少连接创建时间。
- 改善拥塞控制 。
- 多路复用，但不要带来队头阻塞。
- 前向纠错。
- 连接迁移。

FEC（Forward Error Correction，前向纠错）试图通过在邻近的数据包中添加一个QUIC数据包的部分数据来减少数据包重传的需求。这个想法是，如果只丢了
一个数据包，那应该可以从成功传送的数据包中重新组合出该数据包。

连接迁移旨在减少连接创建的开销，它通过支持连接在网络之间迁移来实现。

### 9.2.1 QUIC的性能优势

9.2.2 QUIC和网络技术栈

图9.18　QUIC在HTTP技术栈中所处的位置

![](https://res.weread.qq.com/wrepub/epub_32517945_359)

**QUIC不会替代HTTP/2，但它会接管传输层的一些工作，在上层运行较轻的HTTP/2实现。**

### 9.2.3 什么是UDP，为什么QUIC基于它

为什么不改进TCP

- 主要缺点是此类改进的实施速度慢。

为什么不使用SCTP

- SCTP的采用率很低，这主要是因为到目前为止TCP已经足够好。因此，迁移到SCTP可能与升级TCP需要一样长的时间。

**为什么不直接使用IP**

- 直接使用IP与直接使用SCTP具有相同的问题。该协议必须在操作系统级别实现，因为很少有应用程序可以直接访问IP数据包。

UDP的优点

UDP是一种基础协议，也在内核中实现。在它之上的任何东西都需要在应用层中构建，也就是所说的用户空间。在内核之外构建，可以通过部署应用程序来实现快速创
新，无论是在服务端还是在客户端。

图9.19　查看www.google.com上部署的QUIC版本

![](https://res.weread.qq.com/wrepub/epub_32517945_360)

QUIC会一直使用UDP吗

### 9.2.4 标准化QUIC

两个版本的QUIC：gQUIC和iQUIC

- gQUIC（Google QUIC）和iQUIC（IETF QUIC）

gQUIC和iQUIC的区别

- 随着两个版本协议的发展，它们之间的区别会越来越大
- Google使用自定义加密设计，而iQUIC使用TLSv1.3

QUIC标准

- QUIC Invariants —— QUIC中恒定不变的部分
- QUIC Transport —— 核心传输协议
- QUIC Recovery —— 丢包检测和拥塞控制
- QUIC TLS —— QUIC中如何使用TLS加密
- HTTP/3 —— 主要基于HTTP/2，但有一些不同
- QUIC QPACK —— 使用QUIC的HTTP协议的首部压缩

需要注意的一点是，QUIC旨在成为一种通用的协议，HTTP只是它的一种用途。虽然HTTP目前是QUIC的主要用例，也是工作组目前正在关注的焦点，但该协议的设计
考虑了潜在的其他应用场景。

### 9.2.5 HTTP/2和QUIC的不同

QUIC和HTTPS

HTTPS内置于QUIC中，与HTTP/2不同，QUIC不能用于未加密的HTTP连接。做出这个选择的原因与HTTP/2相同，无论从实际使用上，还是人们的意愿上，只能通过
HTTPS进行Web浏览

创建一个QUIC连接

因为QUIC是基于UDP的，连接到Web服务器的浏览器必须先使用TCP连接，然后再升级到QUIC。这个过程就需要依赖基于TCP的HTTP，这就抵消了QUIC带来的一
个关键好处（大量减少连接创建时间）。有一些变通方法，比如同时尝试TCP和UDP，或者就接受第一次的性能损耗，并记住下次服务器使用QUIC。

QPACK

QUIC旨在消除连接层顺序传输数据包的要求，以允许流独立处理。HPACK仍然需要这种保证（至少对于HEADERS帧），因此它重新引入了队头阻塞，而这正是它试图
解决的问题。

因此，HTTP/3需要有一种HPACK的变体，也就是QPACK（原因显而易见）。这个变体很复杂

其他区别

一些传输层协议的帧从HTTP/3层中被移除了（例如PING和WINDOW_UPDATE帧），移动到了核心QUIC-Transport层，这不是针对HTTP的（这是合理的，因为这些
帧很可能会用于基于QUIC的非HTTP协议）。

### 9.2.6 QUIC的工具

### 9.2.7 QUIC实现

- Caddy

### 9.2.8 你应该使用QUIC吗

与SPDY不同，gQUIC并未被更广泛的社区所接受，而且iQUIC似乎也不太可能现在被标准化。所以，除非你使用了Google Cloud Platform，否则不推荐你使用
QUIC。

## 总结

- 在TCP和HTTPS层中，当前的HTTP网络栈存在若干低效率因素。
- 由于TCP的连接建立延迟和谨慎的拥塞控制算法，TCP连接达到最大容量需要时间，HTTPS握手会增加更多时间。
- 有一些创新可以解决这些低效问题，但它们的推出速度很慢，特别是TCP中的创新。
- QUIC是一种基于UDP的新协议。
- 通过使用UDP，QUIC谋求比TCP创新的速度更快。
- QUIC的创建基于HTTP/2，其使用了许多相同的概念，它是在原来的基础上再创新。
- QUIC不仅适用于HTTP，它未来也可以用于其他协议。
- 基于QUIC的HTTP将被称为HTTP/3。
- QUIC有两个版本：Google QUIC（gQUIC），当前有少量应用，但没有被标准化；IETF QUIC（iQUIC），正在标准化过程中。
- 在iQUIC被批准成为正式标准时，gQUIC将被取代，就像HTTP/2取代了SPDY。
