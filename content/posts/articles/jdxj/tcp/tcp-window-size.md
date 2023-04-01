---
title: "TCP窗口大小"
date: 2023-04-01T21:10:57+08:00
tags:
  - tcp
---

![](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/f8301dda93e1401599dc68ef1d64af97~tplv-k3u1fbpfcp-zoom-in-crop-mark:3024:0:0:0.awebp)

- window size只有16位, 起初表示最大窗口为65535B
- 后来不够用就引入了TCP窗口缩放选项, 范围为0~14
  - 0: 不缩放
  - !=0: 窗口大小为 windowSize * 2^n
- 窗口缩放在握手时指定
