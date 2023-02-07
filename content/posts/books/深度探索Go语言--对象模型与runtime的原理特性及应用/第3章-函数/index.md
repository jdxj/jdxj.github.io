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
