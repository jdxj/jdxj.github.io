---
title: "HTTP 代理原理及实现（一）"
date: 2023-04-17T17:44:43+08:00
tags:
  - http
---

[原文](https://imququ.com/post/web-proxy.html)

普通代理

- 代理服务器解析客户端req, 之后代理服务器向目标服务器发送该req

![](https://st.imququ.com/i/webp/static/uploads/2015/11/web_proxy.png.webp)

隧道代理

- 代理服务器监听CONNECT方法, 之后转发客户端发来的tcp流量到目标服务器

![](https://st.imququ.com/i/webp/static/uploads/2015/11/web_tunnel.png.webp)
