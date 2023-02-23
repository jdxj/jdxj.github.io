---
title: "第2章 指针"
date: 2023-02-05T13:55:47+08:00
draft: false
---

## 2.1 指针构成

```go
var p *int
```

无论指针的元素类型是什么，指针变量本身的格式都是一致的，即一个无符号整型，变量大小能够容纳当前平台的地址。例如在386架构上是一个32位无符号整
型，在amd64架构上是一个64位无符号整型。

有着不同元素类型的指针被视为不同类型，这是语言设计层面强加的一层安全限制，因为不同的元素类型会使编译器对同一个内存地址进行不同的解释。

### 2.1.1 地址

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P35_4878.jpg)
![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P35_4890.jpg)

### 2.1.2 元素类型

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P36_4901.jpg)
![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P36_4909.jpg)

后两条指令由MOVQ变为MOVL

## 2.2 相关操作

### 2.2.1 取地址

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P37_4923.jpg)
![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P37_4934.jpg)

LEAQ指令的作用就是取得main.n的地址并装入AX寄存器中。后面的MOVQ指令则把AX的值复制到返回值p。

- 这里获取的是一个包级别变量n的地址，等价于C语言的全局变量，变量n的地址是在编译阶段静态分配的，所以LEAQ指令通过位移寻址的方式得到了main.n
  的地址。
- LEAQ同样也支持基于基址和索引获取地址

Go语言通过逃逸分析机制避免返回局部变量地址所引发的问题, 实际上在堆上分配

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P37_4942.jpg)

### 2.2.2 解引用 {#dereference}

1. 空指针异常

- 所谓空指针，就是地址值为0的指针。按照操作系统的内存管理设计，进程地址空间中地址为0的内存页面不会被分配和映射，保留地址0在程序代码中用作无效
指针判断，所以对空指针进行解引用操作就会造成程序异常崩溃
- 遭遇空指针异常并非语言设计方面的缺陷，而是程序逻辑上的Bug。

2. 野指针问题

- 在C语言中, 未初始化的指针变量是随机值, 会绕过代码中的空指针判断逻辑，从而造成内存访问错误。
- Go语言中声明的变量默认都会初始化为对应类型的零值，指针类型变量都会初始化为nil

3. 悬挂指针问题

- 指程序过早地释放了内存，而后续代码又对已经释放的内存进行访问，从而造成程序出现错误或异常。
- Go语言实现了自动内存管理，由GC负责释放堆内存对象。GC基于标记清除算法进行对象的存活分析，只有明确不可达的对象才会被释放

### 2.2.3 强制类型转换

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P39_4958.jpg)
![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P39_4966.jpg)

### 2.2.4 指针运算

假如有一个元素类型为int的指针p，要把p移动到下一个int的位置，在C语言中可以通过指针的自增运算实现，代码如下：

```c
++p;
```

在Go语言中等价的代码如下：

```go
p = (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(p))+unsafe.Sizeof(*p)))
```



## 2.3 unsafe包

经典的类型转换

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P40_2499.jpg)

图2-1 String Header和Slice Header的结构

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P40_2503.jpg)

如果不经意修改了slice就可能会造成程序逻辑错误。

### 2.3.1 标准库与keyword

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P41_2515.jpg)

- ArbitraryType在这里只是用于文档目的，实际上并不属于unsafe包，它可以表示任意的Go表达式类型。
- Sizeof()函数用来返回任意类型的大小
- Offsetof()函数用来返回任意结构体类型的某个字段在结构体内的偏移
- Alignof()函数用来返回任意类型的对齐边界
- 最重要的是这3个函数的返回值都是常量。

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P41_2523.jpg)
![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P41_2531.jpg)

- 这条MOVQ指令直接向返回值o中写入了立即数8，也就说明Sizeof()函数在编译阶段就被转换成了立即数。
- 上述测试方法同样适用于Offsetof()函数和Alignof()函数。

### 2.3.2 关于uintptr

很多人都认为uintptr是个指针，其实不然。不要对这个名字感到疑惑，它只不过是个**uint**，大小与当前平台的指针宽度一致。因为unsafe.Pointer可
以跟uintptr互相转换，所以Go语言中可以把指针转换为uintptr进行数值运算，然后转换回原类型，以此来模拟C语言中的指针运算。

需要注意的是，不要用uintptr来存储堆上对象的地址。具体原因和GC有关，GC在标记对象的时候会跟踪指针类型，而
**uintptr不属于指针，所以会被GC忽略, 造成堆上的对象被认为不可达，进而被释放**。用unsafe.Pointer就不会存在这个问题了，unsafe.Pointer类
似于C语言中的void∗，虽然未指定元素类型，但是本身类型就是个指针。

> 参考[聊一个string和[]byte转换问题]({{< ref "posts/collections/huoding/聊一个string和[]byte转换问题.md" >}})

### 2.3.3 内存对齐

- 硬件的实现一般会将内存的读写对齐到数据总线的宽度，这样既可以降低硬件实现的复杂度，又可以提升传输的效率。
- Go语言的内存对齐规则参考了两方面因素：一是数据类型自身的大小，复合类型会参考最大成员大小；二是硬件平台机器字长。

机器字长是指计算机进行一次整数运算所能处理的二进制数据的位数，在x86平台可以理解成数据总线的宽度。当数据类型自身大小小于机器字长时，会被对齐
到自身大小的整数倍；当自身大小大于机器字长时，会被对齐到机器字长的整数倍。

通过unsafe.Sizeof()函数和unsafe.Alignof()函数可以得到目标数据类型的大小和对齐边界

> 地址应该是Alignof的倍数.

表2-1 常见内置类型的大小和对齐边界

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-T42_16075.jpg)

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P43_5189.jpg)
![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P44_5195.jpg)

通过调整结构体成员的位置，尽量避免编译器添加padding

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P44_5203.jpg)
![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P44_5207.jpg)

## 2.4 本章小结
