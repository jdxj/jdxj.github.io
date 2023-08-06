---
title: "log模块"
date: 2023-08-06T16:38:47+08:00
summary: log阶段
tags:
  - nginx
---

# [log_format指令](https://nginx.org/en/docs/http/ngx_http_log_module.html#log_format)

- combined格式

# [access_log指令](https://nginx.org/en/docs/http/ngx_http_log_module.html#access_log)

日志路径可以包含变量, 不打开cache时, 每记录一条日志都要打开, 关闭日志文件, 有性能问题

# [open_log_file_cache](https://nginx.org/en/docs/http/ngx_http_log_module.html#open_log_file_cache)
