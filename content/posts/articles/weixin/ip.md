---
title: '如何获取客户端真实 IP？从 Gin 的一个 "Bug" 说起'
date: 2023-03-10T13:56:11+08:00
draft: false
---

边缘节点应该用以下方法来设置

```
proxy_set_header X-Forwarded-For $remote_addr;
```

[原文](https://mp.weixin.qq.com/s/C-Xf6haLrOWkmBm2lRTdEQ)
