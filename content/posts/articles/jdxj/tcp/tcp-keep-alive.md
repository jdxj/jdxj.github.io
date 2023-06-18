---
title: "TCP Keep Alive"
date: 2023-06-18T16:27:43+08:00
tags:
  - tcp
---

- 发送心跳周期 net.ipv4.tcp_keepalive_time = 7200
  - 7200s 没数据交互时, 启动探测
- 探测包发送间隔 net.ipv4.tcp_keepalive_intvl = 75
- 探测包重试次数 net.ipv4.tcp_keepalive_probes = 9
