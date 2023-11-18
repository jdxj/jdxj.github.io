---
title: "Configuring MySQL to use minimal memory"
date: 2023-11-16T11:44:54+08:00
tags:
  - mysql
---

[原文](http://www.tocker.ca/2014/03/10/configuring-mysql-to-use-minimal-memory.html)

```
# /etc/my.cnf:
innodb_buffer_pool_size=5M
innodb_log_buffer_size=256K
max_connections=50
key_buffer_size=8
thread_cache_size=1
host_cache_size=0
innodb_ft_cache_size=1600000
innodb_ft_total_cache_size=32000000

# per thread or per operation settings
thread_stack=131072
sort_buffer_size=32K
read_buffer_size=8200
read_rnd_buffer_size=8200
max_heap_table_size=16K
tmp_table_size=1K
bulk_insert_buffer_size=0
join_buffer_size=128
net_buffer_length=1K
innodb_sort_buffer_size=64K

#settings that relate to the binary log (if enabled)
binlog_cache_size=4K
binlog_stmt_cache_size=4K
```

- [innodb_buffer_pool_size](https://cloud.tencent.com/developer/article/1834199#:~:text=nnodb%E5%8F%82%E6%95%B0-,innodb_buffer_pool_size,-%E8%BF%99%E4%B8%AA%E6%98%AFInnodb)
- [innodb_log_buffer_size](https://cloud.tencent.com/developer/article/1834199#:~:text=innodb_log_file_size%3D256M-,innodb_log_buffer_size,-%E4%BA%8B%E5%8A%A1%E5%9C%A8%E5%86%85%E5%AD%98)
- [max_connections](https://cloud.tencent.com/developer/article/2076584)
- [key_buffer_size](https://segmentfault.com/a/1190000016509398?u_atoken=c677c59a-2f57-4aa4-8598-ca6bfbe27ffd&u_asession=01KgT0K9QGdF5ow_DHlwf5kNQeOmeCwFFzRPvZDMaIJJ-WXEcFc9UUSlHw2_OH0XVv_WPf92SPdniySFNSEkWa3Nsq8AL43dpOnCClYrgFm6o&u_asig=05Hx7n3q-etCphSLOzvgXLwCwnmJD5-yI_E7Ff8i50yRvYx0n8F36jpVv7yur0QJOPxmja1nhXiS-c7TxW22yAQUrIqoC7FhZ0mdKPs03KSNwhvd1AN7rz4k2JTLKhvu3s8A7-G9C3jEqAGdkB4IL9Ogk4LgUZHqNa-HmisEo2NSFUbx8mcGS-bAdwAkHugkNJksmHjM0JOodanL5-M1Qs1YyiTfzj2Bnd6SkLJNCPY2-6LlJDQJEBkfXHjXZk54prIGFavrILYh1e6SvSAZTaHorbYkUbnxpk0IM8ajdNtdXY94r_LXIIil3Y3aVPRGAe&u_aref=Gzj7DnfDtljSF0azYOGIDdew4%2Fk%3D)
- [thread_cache_size](https://developer.aliyun.com/article/1173107)
- [host_cache_size](https://cloud.tencent.com/developer/article/2008560)
- [innodb_ft_cache_size](https://www.cnblogs.com/datamining-bio/p/17082860.html)
- [innodb_ft_total_cache_size](https://www.cnblogs.com/datamining-bio/p/17082860.html#:~:text=innodb_ft_total_cache_size)
- [thread_stack](https://www.modb.pro/db/43243#:~:text=variables%20like%20%22thread_cache_size%22%3B-,Thread_stack,-%EF%BC%9A%E6%AF%8F%E4%B8%AA%E8%BF%9E%E6%8E%A5)
- [sort_buffer_size](https://blog.csdn.net/sdyu_peter/article/details/116212312)
- [read_buffer_size](https://www.modb.pro/db/55592#:~:text=%E5%8C%BA%E4%BD%BF%E7%94%A8%E5%86%85%E5%AD%98(-,read_buffer_size,-)%EF%BC%9A%0A%E8%BF%99%E9%83%A8%E5%88%86)
- [read_rnd_buffer_size](https://www.modb.pro/db/55592#:~:text=%E5%8C%BA%E4%BD%BF%E7%94%A8%E5%86%85%E5%AD%98(-,read_rnd_buffer_size,-)%EF%BC%9A%0A%E5%92%8C%E9%A1%BA%E5%BA%8F)
- [max_heap_table_size](https://www.cnblogs.com/sunss/archive/2011/01/10/1932004.html)
- [tmp_table_size](https://www.cnblogs.com/sunss/archive/2011/01/10/1932004.html)
- [bulk_insert_buffer_size](https://www.cnblogs.com/ggjucheng/archive/2012/11/11/2765336.html#:~:text=%E6%8F%90%E9%AB%98%E6%9F%A5%E8%AF%A2%E6%95%88%E7%8E%87%E3%80%82-,bulk_insert_buffer_size,-(thread))
- [join_buffer_size](https://www.modb.pro/db/628790)
- [net_buffer_length](https://blog.csdn.net/cschmin/article/details/123328160#:~:text=MySQL%20%E7%B3%BB%E7%BB%9F%E5%8F%98%E9%87%8F-,net_buffer_length,-%EF%BC%8C%E4%BB%8E%E5%AE%83%E7%9A%84%E5%90%8D%E5%AD%97)
- [innodb_sort_buffer_size](https://juejin.cn/s/mysql%20innodb_sort_buffer_size)
- [binlog_cache_size](https://juejin.cn/s/mysql%20binlog_cache_size)
- [binlog_stmt_cache_size](https://juejin.cn/s/mysql%20binlog_stmt_cache_size)
