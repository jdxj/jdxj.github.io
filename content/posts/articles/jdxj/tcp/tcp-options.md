---
title: "TCP Options"
date: 2023-04-01T21:16:50+08:00
tags:
  - tcp
---

![](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/22f5fa146fd3495ab40dc6f335f0d0b5~tplv-k3u1fbpfcp-zoom-in-crop-mark:3024:0:0:0.awebp)

可选项格式

![](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/a3e17de1391c4414b9405146778da880~tplv-k3u1fbpfcp-zoom-in-crop-mark:3024:0:0:0.awebp)

# 时间戳选项

TCP Timestamps Option，TSopt

由4部分组成

- 类别（kind）
- 长度（Length）
- 发送方时间戳（TS value）
- 回显时间戳（TS Echo Reply）

![](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2019/6/14/16b54c4be8611658~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)

是否使用时间戳选项是在三次握手里面的 SYN 报文里面确定的。

- 发送方发送数据时，将一个发送时间戳 1734581141 放在发送方时间戳TSval中
- 接收方收到数据包以后，将收到的时间戳 1734581141 原封不动的返回给发送方，放在TSecr字段中，同时把自己的时间戳 3303928779 放在TSval中
- 后面的包以此类推

![](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2019/6/14/16b54c4c5c7ae349~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)

Timestamps 选项的作用

- 两端往返时延测量（RTTM）
- 序列号回绕（PAWS）

## 测量 RTTM

![](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2020/3/22/17102ef66fd1e657~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)

启用Timestamps选项后可以避免重传包无法计算rtt的问题.

## PAWS

TCP 的窗口经过窗口缩放可以最高到 1GB（2^30)，在高速网络中，序列号在很短的时间内就会被重复使用。

![](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2020/3/22/17102ef66f71cbd6~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)

如果有 Timestamps 的存在，内核会维护一个为每个连接维护一个 ts_recent 值，记录最后一次通信的的 timestamps 值，在 t7 时间点收到迷途数据包 2
时，由于数据包 2 的 timestamps 值小于 ts_recent 值，就会丢弃掉这个数据包。等 t8 时间点真正的数据包 6 到达以后，由于数据包 6 的
timestamps 值大于 ts_recent，这个包可以被正常接收。

