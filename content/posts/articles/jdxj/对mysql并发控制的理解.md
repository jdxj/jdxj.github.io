---
title: "对MySQL并发控制的理解"
date: 2023-11-04T17:18:18+08:00
tags:
  - mysql
---

可以先从编程语言中的锁思考, 比如Go. 在Go中访问并发资源可能会使用互斥Mutex, 如果想进一步提升性能,
那么可以使用RWMutex. 如果极致一点可以利用无锁编程技术CAS.

在MySQL中也有类似的想法. 首先对于更新操作并发访问肯定是不允许的, 所以要用锁.

对于读的情况

- 如果有多个线程对相同数据行只读, 那么无需使用锁, 只需利用事务隔离级别(MVCC/快照读)来保证读取前后一致性
- 如果存在先读数据行, 之后再更新数据行的逻辑不能用快照读, 应当使用当前读(for update)

# 相关阅读

- [事务篇](https://xiaolincoding.com/mysql/transaction/mvcc.html)
- [锁篇](https://xiaolincoding.com/mysql/lock/mysql_lock.html#%E5%85%A8%E5%B1%80%E9%94%81)
