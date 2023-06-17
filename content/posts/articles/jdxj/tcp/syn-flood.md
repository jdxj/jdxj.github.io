---
title: "SYN Flood 攻击"
date: 2023-04-05T11:06:54+08:00
tags:
  - tcp
---

客户端大量伪造 IP 发送 SYN 包，服务端回复的 ACK+SYN 去到了一个「未知」的 IP 地址，势必会造成服务端大量的连接处于 SYN_RCVD 状态，而服务器的
半连接队列大小也是有限的，如果半连接队列满，也会出现无法处理正常请求的情况。

![](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2019/6/29/16ba36e681b24ff3~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)

如何应对 SYN Flood 攻击

- 调大net.ipv4.tcp_max_syn_backlog的值，不过这只是一个心理安慰，真有攻击的时候，这个再大也不够用。
- 重试次数由 /proc/sys/net/ipv4/tcp_synack_retries控制，默认情况下是 5 次，当收到SYN+ACK故意不回 ACK 或者回复的很慢的时候，调小这个值
  很有必要。
- [tcp_syncookies]({{< ref "posts/articles/jdxj/tcp/syn-cookie.md" >}})
- net.core.netdev_max_backlog 接收自网卡, 但未被内核协议栈处理的报文队列长度
- [net.ipv4.tcp_abort_on_overflow]({{< ref "posts/articles/jdxj/tcp/连接队列.md#tcp-abort-on-overflow" >}})
  可以返回RST
