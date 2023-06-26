---
title: "Golang中的热重启"
date: 2023-06-25T21:35:58+08:00
tags:
  - go
---

[原文](https://mp.weixin.qq.com/s?src=11&timestamp=1687697673&ver=4612&signature=qlI5-v11MvpO4HQaMeyRXmZm69zqrnWaKWVnT*QKsGmg6r-i17j5zCxZbkAaK5PXy8V2*lGvnpFBLSscDCDFPZnDHXslrLtXP0T*PoON8hGvuzCb2tW2dLMNMNwy8T5Z&new=1)

# 热重启的原理

1. 监听重启信号；
2. 收到重启信号时fork子进程，同时需要将服务监听的socket文件描述符传递给子进程；
3. 子进程接收并监听父进程传递的socket；
4. 等待子进程启动成功之后，停止父进程对新连接的接收；
5. 父进程退出，重启完成
