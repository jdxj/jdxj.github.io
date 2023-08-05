---
title: "rewrite模块"
date: 2023-08-05T18:28:33+08:00
summary: rewrite阶段
tags:
  - nginx
---

# [return指令](https://nginx.org/en/docs/http/ngx_http_rewrite_module.html#return)

nginx自定义状态码

- 444: 关闭连接, 不向用户返回内容

http

- 301永久重定向
- 302临时重定向, 禁止被缓存
- 303临时重定向, 允许改变方法, 禁止被缓存
- 307临时重定向, 不允许改变方法, 禁止被缓存
- 308永久重定向, 不允许改变方法

# [error_page指令](https://nginx.org/en/docs/http/ngx_http_core_module.html#error_page)

# [rewrite指令](https://nginx.org/en/docs/http/ngx_http_rewrite_module.html#rewrite)

# [if指令](https://nginx.org/en/docs/http/ngx_http_rewrite_module.html#if)
