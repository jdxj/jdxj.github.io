---
title: "Redis内部编码"
date: 2023-04-21T16:30:47+08:00
tags:
  - redis
---

![](redis内部编码.drawio.svg)

查看某key的内部编码

```bash
redis> RPUSH lst 1 3 5 10086 "hello" "world"
(integer)6
redis> OBJECT ENCODING lst
"ziplist"
```

- [raw]({{< ref "posts/books/Redis设计与实现/第2章-简单动态字符串.md" >}})
- [embstr]({{< ref "posts/books/Redis设计与实现/第8章-对象.md#string" >}})
- [hashtable]({{< ref "posts/books/Redis设计与实现/第4章-字典.md" >}})
- [linkedlist]({{< ref "posts/books/Redis设计与实现/第3章-链表.md" >}})
- [ziplist]({{< ref "posts/books/Redis设计与实现/第7章-压缩列表.md" >}})
- [intset]({{< ref "posts/books/Redis设计与实现/第6章-整数集合.md" >}})
- [skiplist]({{< ref "posts/books/Redis设计与实现/第5章-跳跃表.md" >}})