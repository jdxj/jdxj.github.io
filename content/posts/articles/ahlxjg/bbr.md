---
title: "Debian11开启BBR"
date: 2023-04-08T09:11:45+08:00
tags:
  - tcp
  - bbr
---

Debian11应该默认开启了, 验证

```
lsmod | grep bbr
#或成功则会出现类似的内容 tcp_bbr  20480  1
```

参考

- [Debian11开启bbr](https://www.cnblogs.com/ahlxjg/p/16108241.html)
