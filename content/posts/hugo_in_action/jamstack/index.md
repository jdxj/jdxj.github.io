---
title: "Jamstack"
date: 2022-07-22T13:55:25+08:00
draft: false
---

# jam 含义

![dd](./figure1.4.png)

- m: Markup
- j: Javascript
- a: APIs

结合例图来看, Hugo 用 M 和 A (其实是 js, 这里表示应用逻辑) 生成静态文件, 通过 CDN 分发到 Client 后, 可以通过 js 生成动态内容甚至是访问
由云提供商开放的 API.

这个 stack 使得用户的维护工作变得很少.