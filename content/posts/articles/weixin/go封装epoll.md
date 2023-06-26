---
title: "在 golang 中是如何对 epoll 进行封装的？"
date: 2023-06-25T20:41:59+08:00
tags:
  - go
  - epoll
---

[原文](https://mp.weixin.qq.com/s?src=11&timestamp=1687695474&ver=4612&signature=vgBnMEPklFiisTC8lihEgTPiLAz42pRv7GCGt092qPsL8BXQBq6luq7PMN6QzUxLpEFVSW0aHQS7Flg9Xtk4eaRvKqKR8c0ynJj2lgorXdhUN6*DZqbk36e1GmJAAlLK&new=1)

简单来说就是封装非阻塞fd, 用户代码在调用非阻塞fd时, 由go来实现调度, 将本goroutine阻塞, 但是不阻塞该线程.
