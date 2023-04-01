---
title: "TCP Flags"
date: 2023-04-01T21:07:11+08:00
tags:
  - tcp
---

![](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/59260fc36dfa468693085a6ac4600448~tplv-k3u1fbpfcp-zoom-in-crop-mark:3024:0:0:0.awebp)

这些标记可以组合使用，比如 SYN+ACK，FIN+ACK 等

- SYN（Synchronize）：用于发起连接数据包同步双方的初始序列号
- ACK（Acknowledge）：确认数据包
- RST（Reset）：这个标记用来强制断开连接，通常是之前建立的连接已经不在了、包不合法、或者实在无能为力处理
- FIN（Finish）：通知对方我发完了所有数据，准备断开连接，后面我不会再发数据包给你了。
- PSH（Push）：告知对方这些数据包收到以后应该马上交给上层应用，不能缓存起来
