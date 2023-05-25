---
title: "LOD"
date: 2023-05-24T14:51:55+08:00
---

Law of Demeter, 迪米特法则

另外一个更加达意的名字: 最小知识原则, The Least Knowledge Principle

- 每个模块（unit）只应该了解那些与它关系密切的模块（units: only units “closely” related to the current unit）的有限知识（knowledge）。
  或者说，每个模块只和自己的朋友“说话”（talk），不和陌生人“说话”（talk）。

> 用接口隔离
