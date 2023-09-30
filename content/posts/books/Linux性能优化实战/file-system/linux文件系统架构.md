---
title: "Linux文件系统架构"
date: 2023-09-30T12:08:56+08:00
---

![](https://static001.geekbang.org/resource/image/72/12/728b7b39252a1e23a7a223cdf4aa1612.png?wh=530*590)

在 Linux 中一切皆文件。不仅普通的文件和目录，就连块设备、套接字、管道等，也都要通过统一的文件系统来管理。

- 索引节点，简称为 inode，用来记录文件的元数据，比如 inode 编号、文件大小、访问权限、修改日期、数据的位置等。索引节点和文件一一对应，它跟
  文件内容一样，都**会被持久化存储到磁盘中**。所以记住，索引节点同样占用磁盘空间。
- 目录项，简称为 dentry，用来记录文件的名字、索引节点指针以及与其他目录项的关联关系。多个关联的目录项，就构成了文件系统的目录结构。不过，
  不同于索引节点，目录项是由内核维护的一个**内存数据结构**，所以通常也被叫做目录项缓存。

目录项和索引节点的关系是多对一，你可以简单理解为，一个文件可以有多个别名。

磁盘读写的最小单位是扇区，然而扇区只有 512B 大小，如果每次都读写这么小的单位，效率一定很低。所以，文件系统又把连续的扇区组成了逻辑块，然后每
次都以逻辑块为最小单元，来管理数据。常见的逻辑块大小为 4KB，也就是由连续的 8 个扇区组成。

- 超级块，存储整个文件系统的状态。如索引节点和逻辑块的使用情况等。
- 索引节点区，用来存储索引节点。
- 数据块区，则用来存储文件数据。
- 逻辑块，是由连续磁盘扇区构成的最小读写单元，用来存储文件数据。

目录项是一个内存缓存；而超级块、索引节点和逻辑块，都是存储在磁盘中的持久化数据。

# 索引节点的容量

也就是Inode个数, 是在格式化磁盘时设定好的，一般由格式化工具自动生成。当你发现索引节点空间不足，但磁盘空间充足时，很可能就是过多小文件导致的。

```bash
# -i 参数，查看索引节点的使用情况
$ df -i /dev/sda1 
Filesystem      Inodes  IUsed   IFree IUse% Mounted on 
/dev/sda1      3870720 157460 3713260    5% / 
```

# 缓存

可以用 free 或 vmstat，来观察页缓存的大小。

```bash
$ cat /proc/meminfo | grep -E "SReclaimable|Cached" 
Cached:           748316 kB 
SwapCached:            0 kB 
SReclaimable:     179508 kB 
```

查看文件系统中的目录项和索引节点缓存

内核使用 Slab 机制，管理目录项和索引节点的缓存。/proc/meminfo 只给出了 Slab 的整体大小，具体到每一种 Slab 缓存，还要查看
/proc/slabinfo 这个文件。

```bash
$ cat /proc/slabinfo | grep -E '^#|dentry|inode' 
# name            <active_objs> <num_objs> <objsize> <objperslab> <pagesperslab> : tunables <limit> <batchcount> <sharedfactor> : slabdata <active_slabs> <num_slabs> <sharedavail> 
xfs_inode              0      0    960   17    4 : tunables    0    0    0 : slabdata      0      0      0 
... 
ext4_inode_cache   32104  34590   1088   15    4 : tunables    0    0    0 : slabdata   2306   2306      0hugetlbfs_inode_cache     13     13    624   13    2 : tunables    0    0    0 : slabdata      1      1      0 
sock_inode_cache    1190   1242    704   23    4 : tunables    0    0    0 : slabdata     54     54      0 
shmem_inode_cache   1622   2139    712   23    4 : tunables    0    0    0 : slabdata     93     93      0 
proc_inode_cache    3560   4080    680   12    2 : tunables    0    0    0 : slabdata    340    340      0 
inode_cache        25172  25818    608   13    2 : tunables    0    0    0 : slabdata   1986   1986      0 
dentry             76050 121296    192   21    1 : tunables    0    0    0 : slabdata   5776   5776      0 
```

- dentry 行表示目录项缓存
- inode_cache 行，表示 VFS 索引节点缓存
- 其余的则是各种文件系统的索引节点缓存。
