---
title: "第3章 函数"
date: 2023-02-06T21:40:56+08:00
draft: true
---

**图3-1 函数调用发生前**

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P45_5231.jpg)

1. CALL指令会先把下一条指令的地址(返回地址)压入栈中, IP寄存器存储f1的地址

**图3-2 CALL指令执行后**

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P46_5242.jpg)

2. 执行f1()
3. f1()最后有条RET指令, 弹出栈顶的返回地址(应该弹到IP中), 跳到返回地址处继续执行

**图3-3 RET指令执行后**

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P46_5245.jpg)

## 3.1 栈帧

### 3.1.1 栈帧布局

函数栈帧是由编译器管理的。

图3-4 Go语言函数栈帧布局示意图

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P47_5259.jpg)

> 猜测代码逻辑在其他地方, 这里只保存函数状态.

- return address：函数返回地址，占用一个指针大小的空间。实际上是在函数被调用时由CALL指令自动压栈的，并非由被调用函数分配。
- caller’s BP：调用者的栈帧基址，占用一个指针大小的空间。用来将调用路径上所有的栈帧连成一个链表，方便栈回溯之类的操作，
  **只在部分平台架构上存在**。**函数通过将栈指针SP直接向下移动指定大小，一次性分配caller’s BP、locals和args to callee所占用的空间**，
  在x86架构上就是使用SUB指令将SP减去指定大小的。
- locals：局部变量区间，占用若干机器字。用来存放函数的局部变量，根据函数的局部变量占用空间大小来分配，没有局部变量的函数不分配。
- args to callee：调用传参区域，占用若干机器字。这一区域所占空间大小，会按照当前函数调用的所有函数中**返回值**加上**参数**所占用的最大空
  间来分配。当没有调用任何函数时，不需要分配该区间。callee视角的args from caller区间包含在caller视角的args to callee区间内，占用空间
  大小是小于或等于的关系。

{{< embedcode go "code_3_1/main.go" >}}

实际上，代码中的println()函数会被编译器转换为多次调用runtime包中的printlock()、printunlock()、printpointer()、printsp()、
printnl()等函数。前两个函数用来进行并发同步，后3个函数用来打印指针、空格和换行。这5个函数均无返回值，
**只有printpointer()函数有一个参数，会在调用者的args to callee区间占用一个机器字**。

输出结果

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P48_5308.jpg)

表3-1 3个函数栈帧上各区间的大小

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-T49_16079.jpg)

- (1+4+4)*8 = 72B = 0x48B
- 依次类推

**图3-5 main调用f1()函数和f2()函数的栈帧布局图**

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P49_5366.jpg)

调用f2()函数时的栈，在a1和v4之间空了3个机器字。这是因为Go语言的函数是固定栈帧大小的，args to callee是按照所需的最大空间来分配的。

### 3.1.2 寻址方式

**图3-6 SUB指令分配整个栈帧**

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P50_5372.jpg)

如果把图3-6中整个函数栈帧视为一个struct，SP存储着这个struct的起始地址，然后就可以通过基址＋位移的方式来寻址struct的各个字段，也就是栈帧上
的局部变量、参数和返回值。

{{< embedcode go "code_3_3/main.go" >}}

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P51_5394.jpg)

**图3-7 函数fa的栈帧布局**

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P51_5397.jpg)

1. 4～7行和最后两行汇编代码主要用来检测和执行动态栈增长
2. 倒数第4行的RET指令用于在函数执行完成后跳转回返回地址。
3. 第8行的SUBQ指令向下移动栈指针SP，完成当前函数栈帧的分配。倒数第5行的ADDQ指令在函数返回前向上移动栈指针SP，释放当前函数的栈帧。释放与分
   配时的大小一致，均为0x18，即24字节，其中BP of main占用了8字节，args to fb占用了16字节。
4. 第9行代码把BP寄存器的值存到栈帧上的BP of main中，第10行把当前栈帧上BP of main的地址存入BP寄存器中。倒数第6行指令在当前栈帧释放前用
   BP of main的值还原BP寄存器。
5. 第12行和第13行代码，通过AX寄存器中转，把参数n的值从args to fa区间复制到args to fb区间，也就是在fa中把main()函数传递过来的参数n，复
   制到调用fb()函数的参数区间。
6. 第14行代码通过CALL指令调用fb()函数。

Go语言中函数的返回值可以是匿名的，也可以是命名的。对于匿名返回值而言，只能通过return语句为返回值赋值。对于命名返回值，可以在代码中通过其名称
直接操作，与参数和局部变量类似。**无论返回值命名与否，都不会影响函数的栈帧布局**。

### 3.1.3 又见内存对齐

Go语言函数栈帧中返回值和参数的对齐方式与struct类似，对于有返回值和参数的函数，可以把所有返回值和所有参数等价成两个struct，一个返回值
struct和一个参数struct。因为内存对齐方式更加紧凑，所以在支持大量参数和返回值时能够做到较高的栈空间利用率。

验证函数参数和返回值的对齐方式与struct成员的对齐方式是一致的

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P52_5410.jpg)
![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P53_5431.jpg)

栈帧上的参数和返回值到底是分开后作为两个struct，还是按照一个struct来对齐的？

{{< embedcode go "code_3_5/main.go" >}}

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P53_5449.jpg)

- f1()函数有一个返回值和一个参数，而且都是int8类型，如果返回值和参数作为同一个struct进行内存对齐，则a和b应该是紧邻的，中间不会插入padding。
- 可以看到参数a和返回值b并没有紧邻，而是分别按照8字节的边界进行对齐的，也就说明返回值和参数是分别对齐的，不是合并在一起作为单个struct。

局部变量的对齐

{{< embedcode go "code_3_6/main.go" >}}

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P54_5467.jpg)

```go
struct {
    e int8
    a int8
    d int16
    c int32
    b int64
}
```

**局部变量的顺序被重排的, 布局更紧凑**

为什么编译器会对栈帧上局部变量的顺序进行调整以优化内存利用率，但是并不会调整参数和返回值呢？

- 因为函数本身就是对代码单元的封装，参数和返回值属于对外暴露的接口，编译器必须按照函数原型来呈现
- 局部变量属于封装在内部的数据，不会对外暴露，所以编译器按需调整局部变量布局不会对函数以外造成影响。

### 3.1.4 调用约定

对Go语言普通函数的调用约定进行如下总结：

- 返回值和参数都通过栈传递，对应的栈空间由调用者负责分配和释放。
- 返回值和参数在栈上的布局等价于两个struct，struct的起始地址按照平台机器字长对齐。

验证编译器能够参照函数声明来生成传参相关指令

```go
// 第3章 code_3_7.go
package main

import _ "unsafe"

func main() {
	Add(1, 2)
}

// 只有声明
func Add(a, b int)
```

编译

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P55_5503.jpg)

反编译

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P55_5514.jpg)

与Add()函数调用相关的几行汇编代码

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P56_5524.jpg)

- 可以看到两条MOVQ指令分别复制了参数1和2，证明编译阶段参照函数声明生成了正确的传参指令，也就是调用约定在发挥作用。
- CALL指令处，十六进制编码e800000000预留了32位的偏移量空间，在链接阶段会被链接器填写为实际的偏移值。

### 3.1.5 Go 1.17的变化

- 1.16版及以前的版本中都是通过栈来传递参数的，这样实现简单且能支持海量的参数传递，缺点就是与寄存器传参相比性能方面会差一些。
- 在1.17版本中就实现了基于寄存器的参数传递，当然只是在部分硬件架构上实现了。

结合Go自带的反编译工具，在汇编代码层面看一下1.17版本的函数调用是如何通过寄存器传递参数的。

**1. 函数入参的传递方式**

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P56_5534.jpg)
![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P56_5544.jpg)

- 第1～9个参数是依次用AX、BX、CX、DI、SI、R8、R9、R10和R11这9个通用寄存器来传递的
- 从第10个参数开始使用栈来传递 (注意`MOVW $0xb0a, 0(SP)`直接复制了两个数字10, 11)

图3-8 Go 1.17中in12()函数入参的传递方式

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P58_5574.jpg)

**2. 函数返回值的传递方式**

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P57_5561.jpg)

反编译out12()函数

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P58_5582.jpg)
![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P59_5588.jpg)

## 3.2 逃逸分析

### 3.2.1 什么是逃逸分析

{{< embedcode go "code_3_10/main.go" >}}

如果局部变量a仍分配在栈中, 那么返回的地址会变成一个[悬挂指针]({{< ref "posts/books/深度探索Go语言--对象模型与runtime的原理特性及应用/第2章-指针/index.md#dereference" >}})

反编译newInt()函数

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P60_5622.jpg)

- 重点关注上述汇编代码中runtime.newobject()函数调用，该函数是Go语言内置函数new()的具体实现，用来在运行阶段分配单个对象。
- CALL指令之后的两条MOVQ指令通过AX寄存器中转，把runtime.newobject()函数的返回值复制给了newInt()函数的返回值，这个返回值就是动态分配的
  int型变量的地址。

### 3.2.2 不逃逸分析

