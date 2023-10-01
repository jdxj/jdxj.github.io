---
title: "iperf"
date: 2023-10-01T15:14:49+08:00
summary: 网络性能测试工具
tags:
  - optimize
---

目标机器上启动 iperf 服务端

```bash
# -s表示启动服务端，-i表示汇报间隔，-p表示监听端口
$ iperf3 -s -i 1 -p 10000
```

运行 iperf 客户端

```bash
# -c表示启动客户端，192.168.0.30为目标服务器的IP
# -b表示目标带宽(单位是bits/s)
# -t表示测试时间
# -P表示并发数，-p表示目标服务器监听端口
$ iperf3 -c 192.168.0.30 -b 1G -t 15 -P 2 -p 10000
```

iperf 的报告

```bash
[ ID] Interval           Transfer     Bandwidth
...
[SUM]   0.00-15.04  sec  0.00 Bytes  0.00 bits/sec                  sender
[SUM]   0.00-15.04  sec  1.51 GBytes   860 Mbits/sec                  receiver
```