---
title: "发生RST的情景"
date: 2023-04-09T09:58:47+08:00
tags:
  - tcp
---

# 端口未监听

![](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2019/6/19/16b6dd217748a3d1~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)

# 断电丢失连接

![](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2019/6/19/16b6dd2177aff16c~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)

# 设置SO_LINGER为true

如果设置 SO_LINGER 为 true，linger 设置为 0，当调用 socket.close() 时， close 函数会立即返回，同时丢弃缓冲区内所有数据并立即发送 RST 包
重置连接。

参考[SO_LINGER]({{< ref "./socket-options.md" >}})

# 丢失RST

如果客户端收到了这个 RST，就会自然进入CLOSED状态释放连接。如果 RST 依然丢失，客户端只是会单纯的数据丢包了，进入数据重传阶段。如果还一直收不到
RST，会在一定次数以后放弃。

![](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2019/6/19/16b6dd22e77ed16b~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)

# Connection reset by peer

其实就是收到了RST

# Broken pipe

在一个 RST 的套接字继续写数据

![](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2019/6/19/16b7073dca9493c8~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)
