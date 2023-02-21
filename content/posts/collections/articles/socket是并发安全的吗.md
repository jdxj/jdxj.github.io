---
title: "Socket是并发安全的吗"
date: 2023-02-17T14:51:20+08:00
draft: false
tags:
  - tcp
  - udp
---

多线程并发读/写同一个TCP socket是线程安全的，因为TCP socket的读/写操作都上锁了。虽然线程安全，但依然不建议你这么做，因为TCP本身是基于数据流的协
议，一份完整的消息数据可能会分开多次去写/读，内核的锁只保证单次读/写socket是线程安全，锁的粒度并不覆盖整个完整消息。因此建议用一个线程去
读/写TCP socket。

多线程并发读/写同一个UDP socket也是线程安全的，因为UDP socket的读/写操作也都上锁了。UDP写数据报的行为是"原子"的，不存在发一半包或收一半包的问题，
要么整个包成功，要么整个包失败。因此多个线程同时读写，也就不会有TCP的问题。虽然如此，但还是建议用一个线程去读/写UDP socket。

[原文](https://mp.weixin.qq.com/s/rNfBHtpFLxwY7-CiBvkQ5A)
