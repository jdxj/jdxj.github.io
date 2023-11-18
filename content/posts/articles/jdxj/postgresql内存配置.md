---
title: "PostgreSQL内存配置"
date: 2023-11-17T11:18:45+08:00
tags:
  - postgresql
---

```sql
ALTER SYSTEM SET
 max_connections = '20';
ALTER SYSTEM SET
 shared_buffers = '64MB';
ALTER SYSTEM SET
 effective_cache_size = '768MB';
ALTER SYSTEM SET
 maintenance_work_mem = '32MB';
ALTER SYSTEM SET
 random_page_cost = '1.1';
ALTER SYSTEM SET
 work_mem = '6553kB';
ALTER SYSTEM SET
 huge_pages = 'off';
```

# 参考

- [必看！PostgreSQL参数优化](https://www.modb.pro/db/48129)
- [PGTune](https://pgtune.leopard.in.ua/)
