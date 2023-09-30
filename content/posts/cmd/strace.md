---
title: "strace"
date: 2023-09-24T12:38:19+08:00
summary: 跟踪进程系统调用的工具
tags:
  - optimize
---

```bash
$ strace -p 6082
```

```bash
$ strace -p 18940 
strace: Process 18940 attached 
...
mmap(NULL, 314576896, PROT_READ|PROT_WRITE, MAP_PRIVATE|MAP_ANONYMOUS, -1, 0) = 0x7f0f7aee9000 
mmap(NULL, 314576896, PROT_READ|PROT_WRITE, MAP_PRIVATE|MAP_ANONYMOUS, -1, 0) = 0x7f0f682e8000 
write(3, "2018-12-05 15:23:01,709 - __main"..., 314572844 
) = 314572844 
munmap(0x7f0f682e8000, 314576896)       = 0 
write(3, "\n", 1)                       = 1 
munmap(0x7f0f7aee9000, 314576896)       = 0 
close(3)                                = 0 
stat("/tmp/logtest.txt.1", {st_mode=S_IFREG|0644, st_size=943718535, ...}) = 0 
```

- 从 write() 系统调用上，我们可以看到，进程向文件描述符编号为 3 的文件中，写入了 300MB 的数据。
- 后面的 stat() 调用，你可以看到，它正在获取 /tmp/logtest.txt.1 的状态。 这种“点 + 数字格式”的文件，在日志回滚中非常常见。

```bash
$ strace -p pid -f
```