---
title: "Chapter15 The InnoDB Storage Engine"
date: 2022-12-08T15:50:12+08:00
draft: true
tocopen: false
---

## 15.1 Introduction to InnoDB

### 15.1.2 Best Practices for InnoDB Tables

- 使用索引
- 关闭自动提交 (写速度限制)

### 15.1.3 Verifying that InnoDB is the Default Storage Engine

DEFAULT 字样
```mysql
SHOW ENGINES;
SELECT * FROM INFORMATION_SCHEMA.ENGINES;
```



## 15.2 InnoDB and the ACID Model

- A: atomicity.
- C: consistency.
- I: isolation.
- D: durability.

