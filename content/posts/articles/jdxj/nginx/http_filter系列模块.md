---
title: "http_*_filter系列模块"
date: 2023-08-07T16:38:46+08:00
summary: content阶段之前或之后
tags:
  - nginx
---

用于加工响应内容

# copy_filter

复制包体内容

# postpone_filter

处理子请求

# header_filter

构造响应头部

# write_filter

发送响应

# sub_filter

默认不启用, 启用--with-http_sub_module

替换响应中的字符串

# addition

在body前或之后增加内容

默认不启用, 启用--with-http_addition_module
