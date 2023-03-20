---
title: "第8章 堆"
date: 2023-02-22T21:08:46+08:00
---

## 8.1 内存分配

在Go的runtime中，有一系列函数被用来分配内存。

- 例如与new语义相对应的有newobject()函数和newarray()函数，分别负责单个对象的分配和数组的分配。
- 与make语义相对应的有makeslice()函数、makemap()函数及makechan()函数及一些变种，分别负责分配和初始化切片、map和channel。

无论是new系列还是make系列，这些函数的内部无一例外都会调用runtime.mallocgc()函数，它就是Go语言堆分配的关键函数。

### 8.1.1 sizeclasses

Go的堆分配采用了与tcmalloc内存分配器类似的算法，tcmalloc是谷歌公司开发的一款针对C/C++的内存分配器，在对抗内存碎片化和多核性能方面非常优秀

参考tcmalloc实现的内存分配器，内部针对小块内存的分配进行了优化。这类分配器会按照一组预置的大小规格把内存页划分成块，然后把不同规格的内存块放
入对应的空闲链表中

图8-1 tcmalloc内存分配器预置不同规格的链表

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P323_12879.jpg)

在Go源代码runtime包的sizeclasses.go文件中，给出了一组预置的大小规格。

表8-1 sizeclasses预置的大小规格

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-T323_16111.jpg)

续表

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-T325_16115.jpg)

[sizeclasses.go](https://github.com/golang/go/blob/master/src/runtime/sizeclasses.go)

- 第一列是所谓的sizeclass，实际上就是所有规格按空间大小升序排列的序号。
- 第二列是规格的空间大小，单位是字节。
- 第三列表示需要申请多少字节的连续内存，目的是保证划分成目标大小的内存块以后，尾端因不能整除而剩余的空间要小于12.5%。Go使用8192字节作为页面
  大小，底层内存分配的时候都是以整页面为单位的，所以第三列都是8192的整数倍。
- 第四列是第三列与第二列做整数除法得到的商
- 第五列则是余数，分别表示申请的连续内存能划分成多少个目标大小的内存块，以及尾端因不能整除而剩余的空间，也就是在内存块划分的过程中浪费掉的空
  间。
- 最后一列表示的是最大浪费百分比，结合了内存块划分时造成的尾端浪费和内存分配时向上对齐到最接近的块大小造成的块内浪费。

[sizeclasses.go](https://github.com/golang/go/blob/master/src/runtime/sizeclasses.go)文件是被程序生成出来的，源码就在
[mksizeclasses.go](sizeclasses.go文件是被程序生成出来的，源码就在mksizeclasses.go文件中)文件中

### 8.1.2 heapArena

Go语言的runtime将堆地址空间划分成多个arena，在amd64架构的Linux环境下，每个arena的大小是64MB，起始地址也是对齐到64MB的。每个arena都有一
个与之对应的heapArena结构，用来存储arena的元数据

图8-2 area与heapArena的关系

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P326_13891.jpg)

heapArena是在Go的堆之外分配和管理的

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P326_13899.jpg)

bitmap字段是个位图

- 它用两个二进制位来对应arena中一个指针大小的内存单元，所以对于64MB大小的arena来讲，heapArenaBitmapBytes的值是
  64MB/8/8×2＝2MB(64MB/8B=8M, 8M*2b/8=2MB)，这个位图在GC扫描阶段会被用到。
- bitmap第一字节中的8个二进制位，对应的就是arena起始地址往后32字节的内存空间。
- 用来描述一个内存单元的两个二进制位当中，低位用来区分内存单元中存储的是指针还是标量，1表示指针，0表示标量，所以也被称为**指针／标量位**。
- 高位用来表示当前分配的这块内存空间的后续单元中是否包含指针，例如在堆上分配了一个结构体，可以知道后续字段中是否包含指针，如果没有指针就不需
  要继续扫描了，所以也被称为**扫描／终止位**。
- 为了便于操作，一个位图字节中的指针／标量位和扫描／终止位被分开存储，**高4位存储4个扫描／终止位**，**低4位存储4个指针／标量位**。

图8-3 arena起始处分配一个slice对应的bitmap标记

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P327_13905.jpg)

spans数组用来把当前arena中的页面映射到对应的mspan，暂时先认为一个mspan管理一组连续的内存页面

pagesPerArena表示arena中共有多少个页面，用arena大小(64MB)除以页面大小(8KB)得到的结果是8192

图8-4 arena中的页面到mspan的映射

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P328_13912.jpg)

pageInUse是个长度为1024的uint8数组，实际上被用作一个8192位的位图

- 通过它和spans可以快速地找到那些处于mSpanInUse状态的mspan。
- 虽然pageInUse位图为arena中的每个页面都提供了一个二进制位，但是对于那些包含多个页面的mspan，只有第1个页面对应的二进制位会被用到，标记的
  是整个span。

图8-5 pageInUse位图标记使用中的span

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P328_13915.jpg)

pageMarks表示哪些span中存在被标记的对象

- 与pageInUse一样用与起始页面对应的一个二进制位来标记整个span。
- 在GC的标记阶段会原子性地修改这个位图，标记结束之后就不会再进行改动了。
- 清扫阶段如果发现某个span中不存在任何被标记的对象，就可以释放整个span了。

> 不是被标记的才释放吗?

pageSpecials又是一个与pageInUse类似的位图，只不过标记的是哪些span包含特殊设置，目前主要指的是包含finalizers，或者runtime内部用来存储
heap profile数据的bucket。

checkmarks是一个大小为1MB的位图，其中每个二进制位对应arena中一个指针大小的内存单元。当开启调试debug.gccheckmark的时候，checkmarks位图
用来存储GC标记的数据。该调试模式会在STW的状态下遍历对象图，用来校验**并发回收器**能够正确地标记所有存活的对象。

zeroedBase记录的是当前arena中下个还未被使用的页面的位置，相对于arena起始地址的偏移量。页面分配器会按照地址顺序分配页面，所以zeroedBase之
后的页面都还没有被用到，因此还都保持着清零的状态。通过它可以快速判断分配的内存是否还需要进行清零。

