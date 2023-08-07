---
title: "secure_link模块"
date: 2023-08-07T18:39:38+08:00
tags:
  - nginx
---

防盗链

默认不启用, 启用-with-http_secure_link_module

原理, 客户端只能拿到哈希过的url, url需要包含的信息

- 资源位置
- 用户信息
- 时间戳
- 密钥

# [secure_link指令](https://nginx.org/en/docs/http/ngx_http_secure_link_module.html#secure_link)
