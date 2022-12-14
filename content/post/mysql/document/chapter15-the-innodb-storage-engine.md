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



## 15.3 InnoDB Multi-Versioning

InnoDB 在每行记录上添加三个字段

- DB_TRX_ID: 最新进行操作(insert, update, delete)的事务 id
- DB_ROLL_PTR: 指向 undo log 的指针
- DB_ROW_ID: row id, InnoDB 自动生成的聚集索引才有

insert undo log 在事务提交后删除, update undo log 在没有依赖后删除.



## 15.7 InnoDB Locking and Transaction Model

### 15.7.1 InnoDB Locking

#### Shared and Exclusive Locks

就是读写锁

#### Intention Locks

表级锁, 表明一个事务稍后要对该表某行请求的锁
- intention shared lock (IS)
- intention exclusive lock (IX)

`SELECT ... FOR SHARE` sets an IS lock, and `SELECT ... FOR UPDATE` sets an IX lock.

#### Record Locks

`SELECT c1 FROM t WHERE c1 = 10 FOR UPDATE`, 阻止其他事务 insert, update, delete `t.c1 = 10` 的行

#### Gap Locks

锁定索引记录间的间隙, 或者第一条记录之前的间隙, 或者最后一条记录后的间隙

`SELECT c1 FROM t WHERE c1 BETWEEN 10 and 20 FOR UPDATE`, `t.c1 = 15` 不能被插入

#### Next-Key Locks

record lock 和索引记录之前的 gap lock 组合

#### Insert Intention Locks

#### AUTO-INC Locks

AUTO_INCREMENT 列

#### Predicate Locks for Spatial Indexes

### 15.7.5 Deadlocks in InnoDB

隔离级别不影响死锁的可能性, 因为隔离级别负责读操作, 然而死锁发生在写操作