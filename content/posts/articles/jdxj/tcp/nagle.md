---
title: "Nagle"
date: 2023-04-11T20:07:52+08:00
tags:
  - tcp
---

减少发送端频繁的发送小包给对方。

算法思路

```
if there is new data to send
  if the window size >= MSS and available data is >= MSS
    send complete MSS segment now
  else
    if there is unconfirmed data still in the pipe
      enqueue data in the buffer until an acknowledge is received
    else
      send data immediately
    end if
  end if
end if
```

![](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2019/4/23/16a49eab67e29995~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)
