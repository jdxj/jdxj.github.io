---
title: "Socket Options"
date: 2023-04-09T10:05:58+08:00
tags:
  - tcp
---

# SO_LINGER

SO_LINGER 参数是一个 linger 结构体

```c
struct linger {
    int l_onoff;    /* linger active */
    int l_linger;   /* how many seconds to linger for */
};
```

- l_onoff 用来表示是否启用 linger 特性，非 0 为启用，0 为禁用 ，linux 内核默认为禁用。这种情况下 close 函数立即返回，操作系统负责把缓冲队
  列中的数据全部发送至对端
- l_linger 在 l_onoff 为非 0 （即启用特性）时才会生效
  - 如果 l_linger 的值为 0，那么调用 close，close 函数会立即返回，同时丢弃缓冲区内所有数据并立即发送 RST 包重置连接
  - 如果 l_linger 的值为非 0，那么此时 close 函数在阻塞直到 l_linger 时间超时或者数据发送完毕，发送队列在超时时间段内继续尝试发送，如果发送
    完成则皆大欢喜，超时则直接丢弃缓冲区内容 并 RST 掉连接。