---
title: "TCP分手优化"
date: 2023-06-18T16:17:20+08:00
tags:
  - tcp
---

net.ipv4.tcp_tw_reuse = 1

- 开启后, 作为客户端时新连接可以使用仍然处于TIME-WAIT状态的端口
- 由于timestamp的存在, 操作系统可以拒绝迟到的报文, net.ipv4.tcp_timestamps = 1

net.ipv4.tcp_tw_recycle = 0

- 开启后, 同时作为客户端和服务器都可以使用TIME-WAIT状态的端口
- 不安全, 无法避免报文延迟, 重复等给新连接造成混乱

net.ipv4.tcp_max_tw_buckets = 262144

- time_wait状态连接的最大数量
- 超出后直接关闭连接