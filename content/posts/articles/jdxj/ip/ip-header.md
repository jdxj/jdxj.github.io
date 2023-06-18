---
title: "IP Header"
date: 2023-06-17T11:10:16+08:00
tags:
  - ip
---

# IPv4 Header

![](https://nmap.org/book/images/hdr/MJB-IP-Header-800x576.png)

- IHL 头部长度, word(4byte)
- TL 总长度, byte
- Id 分片表示
- Flags 分片控制
  - DF = 1 不能分片
  - MF = 1 中间分片
- FO 分片内偏移, 8byte
- TTL 路由器跳数生存期
- Protocol 承载协议
- HC 校验和