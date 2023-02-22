---
title: "第8章 堆"
date: 2023-02-22T21:08:46+08:00
draft: true
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
