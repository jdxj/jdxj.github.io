---
title: "Redis FAQ"
date: 2022-12-30T10:22:21+08:00
---

# How is Redis different from other key-value stores?

- 拥有复杂的数据类型, 及其上的原子操作
- 是可持久化的内存型数据库, 在数据集比内存大的限制下可以权衡读写速度

# What's the Redis memory footprint?

- 空实例~3MB
- 1百万小 key 的 string 类型~85MB
- 1百万个有5个字段的 hash ~160MB

# What happens if Redis runs out of memory?

- 不接受写, 但是可以读
- 可以配置 [key 淘汰策略](https://redis.io/docs/reference/eviction/)

# How can Redis use multiple CPUs or cores?

- Redis 通常受内存和网络的限制
- 如果想使用多个 CPU, 那么应该今早使用 [Redis Cluster](https://redis.io/docs/management/scaling/)
- Redis 未来会越来越线程化

# What is the maximum number of keys a single Redis instance can hold? What is the maximum number of elements in a Hash, List, Set, and Sorted Set?

- 理论上2^32个 key, 经过测试可以保存2.5亿个 key
- hash, list, set, sorted set 每种类型可以保存2^32个 key

# Where does the name "Redis" come from?

**RE**mote **DI**ctionary **S**erver.

# How is Redis pronounced?

/ˈrɛd-ɪs/