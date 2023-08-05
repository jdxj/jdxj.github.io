---
title: "realip模块"
date: 2023-08-05T18:03:05+08:00
summary: postread阶段
tags:
  - nginx
---

拿到真实用户ip地址

- X-Forwarded-For用于传递ip
- X-Real-IP用于传递用户ip

realip提供的变量

- binary_remote_addr
- remote_addr

默认不启用realip模块, 启用--with-http_realip_module

拿到直接tcp的远程地址

- realip_remote_addr
- realip_remote_port

模块指令

- set_real_ip_from 从被信任的主机获取real_ip
- real_ip_header 从哪里获取real_ip, 顺序X-Real-IP, X-Forwarded-For, proxy_protocol
- real_ip_recursive 取与客户端不同的前一个地址
