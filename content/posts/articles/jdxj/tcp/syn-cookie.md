---
title: "SYN Cookie"
date: 2023-04-05T11:12:14+08:00
tags:
  - tcp
---

用来解决 SYN Flood 攻击的，现在服务器上的 tcp_syncookies 都是默认等于 1，表示连接队列满时启用，等于 0 表示禁用，等于 2 表示始终启用。由
/proc/sys/net/ipv4/tcp_syncookies控制。

SYN Cookie 机制其实原理比较简单，就是在三次握手的最后阶段才分配连接资源

![](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2019/6/29/16ba36e691d04901~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)

- cookie占用序列号空间, 导致tcp可选功能失效, e.g.扩充窗口, 时间戳

参考

- [深入浅出TCP中的SYN-Cookies](https://segmentfault.com/a/1190000019292140)
