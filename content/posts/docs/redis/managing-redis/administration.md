---
title: "Redis administration"
date: 2022-12-30T14:22:21+08:00
---

# Redis setup tips

## Linux

设置内核参数 vm.overcommit_memory = 1
- 表示 Always overcommit
- 写到 /etc/sysctl.conf 配置中后重启
- 或者直接激活 sysctl vm.overcommit_memory=1
- [linux的vm.overcommit_memory的内存分配参数详解](https://www.cnblogs.com/ExMan/p/11586752.html)
- [理解LINUX的MEMORY OVERCOMMIT](http://linuxperf.com/?p=102)

关闭 Transparent Huge Pages
- echo never > /sys/kernel/mm/transparent_hugepage/enabled
- [避免碎片化访问 page](https://cloud.tencent.com/developer/article/1668633)

## Memory

- 启用交换区, 大小等于物理内存
- 配置 maxmemory 选项, 要到达内存限制时报错而不是失败