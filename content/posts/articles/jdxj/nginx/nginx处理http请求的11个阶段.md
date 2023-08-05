---
title: "Nginx处理http请求的11个阶段"
date: 2023-08-05T15:24:33+08:00
tags:
  - nginx
---

# 阶段

| 编号 | 阶段             | 说明                  | 官方模块                                            |
|----|----------------|---------------------|-------------------------------------------------|
| 1  | post_read      | 读完header之后          | realip                                          |
| 2  | server_rewrite | server级别的重写         | rewrite                                         |
| 3  | find_config    | 使用重写后的url匹配location |                                                 |
| 4  | rewrite        | location级别的重写       | rewrite                                         |
| 5  | post_rewrite   | 判断是否需要阶段跳转          |                                                 |
| 6  | preaccess      | 进行访问控制之前进行一些操作      | limit_req, limt_conn                            |
| 7  | access         | 访问控制                | access, auth_basic , auth_request               |
| 8  | post_access    | 访问控制后续处理            |                                                 |
| 9  | precontent     | 生成响应前检查指定文件是否存在     | try_files, mirrors                              |
| 10 | content        | 生成响应阶段              | concat, random_index, index, auto_index, static |
| 11 | log            | 写入日志                | access_log                                      |

# 模块执行顺序

模块执行顺序不是依次执行, 可能会跳过某些模块, 或者跳转到前面阶段的某些模块再次执行

- 模块间的相对顺序参考objs/ngx_modules.c

# 参考

- [Nginx 的 11 个执行阶段详解](https://xie.infoq.cn/article/bc7a344d84c9fedfc6ca871fd)
