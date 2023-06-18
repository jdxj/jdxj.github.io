---
title: "网络层"
date: 2023-06-18T18:34:08+08:00
cover:
  image: http://yupnet.org/popup.html?http://yupnet.org/zittrain/images/fig4-1.jpg
tags:
  - ip
---

![](https://www.caida.org/funding/nets-ipv6/images/Regional_Internet_Registries.png)

# 功能

- ip寻址
- 选路
- 封装打包
- 分片

# 如何传输ip报文

![](http://www.tcpipguide.com/free/diagrams/iphops.png)

- 直接传输
- 本地网络间接传输
  - 内部选路协议
    - RIP
    - OSPF
- 公网间接传输
  - 外部选路协议
    - BGP

路由表

![](http://www.tcpipguide.com/free/diagrams/iprouting.png)

