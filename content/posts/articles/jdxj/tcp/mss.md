---
title: "MSS"
date: 2023-04-02T12:09:48+08:00
tags:
  - tcp
---

# MSS

TCP 为了避免被发送方分片，会主动把数据分割成小段再交给网络层，最大的分段大小称之为 MSS（Max Segment Size）。

```
MSS = MTU - IP header头大小 - TCP 头大小
```

- MSS指TCP最大载荷的大小

![](https://p1-jj.byteimg.com/tos-cn-i-t2oaga2asx/gold-user-assets/2020/2/3/1700a73e8c79596f~tplv-t2oaga2asx-zoom-in-crop-mark:3024:0:0:0.awebp)
