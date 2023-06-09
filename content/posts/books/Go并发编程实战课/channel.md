---
title: "Channel"
date: 2023-06-09T11:49:52+08:00
---

# Channel 的应用场景

- 数据交流：当作并发的 buffer 或者 queue，解决生产者 - 消费者问题。多个 goroutine 可以并发当作生产者（Producer）和消费者（Consumer）。
- 数据传递：一个 goroutine 将数据交给另一个 goroutine，相当于把数据的拥有权 (引用) 托付出去。
- 信号通知：一个 goroutine 可以将信号 (closing、closed、data ready 等) 传递给另一个或者另一组 goroutine 。
- 任务编排：可以让一组 goroutine 按照一定的顺序并发或者串行的执行，这就是编排的功能。
- 锁：利用 Channel 也可以实现互斥锁的机制。

# Channel 的实现原理

![](https://static001.geekbang.org/resource/image/81/dd/81304c1f1845d21c66195798b6ba48dd.jpg?wh=2334*2250)

