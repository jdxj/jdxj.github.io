---
title: "MTU"
date: 2023-04-02T11:38:53+08:00
tags:
  - tcp
---

# MTU

数据链路层传输的帧大小是有限制的，不能把一个太大的包直接塞给链路层，这个限制被称为「最大传输单元（Maximum Transmission Unit, MTU）」

- MTU是指整个IP数据报的大小

以太网帧格式

![](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2020/2/3/1700a73e260cd0cd~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)


# IP分段

当一个 IP 数据包大于 MTU 时，IP 会把数据报文进行切割为多个小的片段(小于 MTU），使得这些小的报文可以通过链路层进行传输

IP 头部中有一个表示分片偏移量的字段，用来表示该分段在原始数据报文中的位置

![](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2020/2/3/1700a73e185162dc~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)

# 路径 MTU

一个包从发送端传输到接收端，中间要跨越很多个网络，每条链路的 MTU 都可能不一样，这个通信过程中最小的 MTU 称为「路径 MTU（Path MTU）」。
