---
title: "MSL"
date: 2023-04-08T11:06:58+08:00
tags:
  - tcp
---

# Max Segment Lifetime

MSL（报文最大生存时间）是 TCP 报文在网络中的最大生存时间。这个值与 IP 报文头的 TTL 字段有密切的关系。

- TTL:  IP 报文最大可经过的路由数

![](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2019/6/14/16b54c4b9038f7aa~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)
![](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2019/6/14/16b54c4b904314f8~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)

- Linux 的套接字实现假设 MSL 为 30 秒，因此在 Linux 机器上 TIME_WAIT 状态将持续 60秒。

# TIME_WAIT 存在的原因

第一个原因是：数据报文可能在发送途中延迟但最终会到达，因此要等老的“迷路”的重复报文段在网络中过期失效，这样可以避免用相同源端口和目标端口创建新连接
时收到旧连接姗姗来迟的数据包，造成数据错乱。

![](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2019/10/15/16dce163cb0bd1d8~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)

第二个原因是确保可靠实现 TCP 全双工终止连接。关闭连接的四次挥手中，最终的 ACK 由主动关闭方发出，如果这个 ACK 丢失，对端（被动关闭方）将重发
FIN，如果主动关闭方不维持 TIME_WAIT 直接进入 CLOSED 状态，则无法重传 ACK，被动关闭方因此不能及时可靠释放。

![](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2019/6/14/16b54c4bb50e0f93~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)

> 可以想象为ack(4)马上就要到了, 但是发生了丢包.
