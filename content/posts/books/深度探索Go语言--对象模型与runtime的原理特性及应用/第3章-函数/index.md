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

{{< embedcode go "code/3_1/main.go" >}}

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

{{< embedcode go "code/3_3/main.go" >}}

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

{{< embedcode go "code/3_5/main.go" >}}

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P53_5449.jpg)

- f1()函数有一个返回值和一个参数，而且都是int8类型，如果返回值和参数作为同一个struct进行内存对齐，则a和b应该是紧邻的，中间不会插入padding。
- 可以看到参数a和返回值b并没有紧邻，而是分别按照8字节的边界进行对齐的，也就说明返回值和参数是分别对齐的，不是合并在一起作为单个struct。

局部变量的对齐

{{< embedcode go "code/3_6/main.go" >}}

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

{{< embedcode go "code/3_10/main.go" >}}

如果局部变量a仍分配在栈中, 那么返回的地址会变成一个[悬挂指针]({{< ref "posts/books/深度探索Go语言--对象模型与runtime的原理特性及应用/第2章-指针/index.md#dereference" >}})

反编译newInt()函数

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P60_5622.jpg)

- 重点关注上述汇编代码中runtime.newobject()函数调用，该函数是Go语言内置函数new()的具体实现，用来在运行阶段分配单个对象。
- CALL指令之后的两条MOVQ指令通过AX寄存器中转，把runtime.newobject()函数的返回值复制给了newInt()函数的返回值，这个返回值就是动态分配的
  int型变量的地址。

### 3.2.2 不逃逸分析

验证new()函数与堆分配是否有必然关系

```go
// 第3章 code_3_11.go
//go:noinline
func New() int {
	p := new(int)
	return *p
}
```

反编译New()函数

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P61_5650.jpg)

MOVQ指令直接把返回值赋值为0，其他的逻辑全都被优化掉了，所以即便是代码中使用了new()函数，只要变量的生命周期没有超过当前函数栈帧的生命周期，
编译器就不会进行堆分配。

### 3.2.3 不逃逸判断

如果把局部变量的地址赋值给包级别的指针变量，应该也会造成变量逃逸

```go
// 第3章 code_3_12.go
var pt *int

//go:noinline
func setNew() {
	var a int
	pt = &a
}
```

反编译setNew()函数

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P62_5672.jpg)

验证逃逸分析的依赖传递性

```go
var pp **int

//go:noinline
func dep() {
	var a int
	var p *int
	p = &a
	pp = &p
}
```

反编译dep()函数

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P62_5689.jpg)

跨包测试

```go
// 第3章 code_3_14.go
package inner

//go:noinline
func RetAry(p *int) *int {
	return p
}

// 第3章 code_3_15.go
package main

//go:noinline
func arg() int {
    var a int
	return *inner.RetAry(&a)
}
```

反编译

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P63_5724.jpg)

阻止编译器参考函数实现的测试

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P64_5741.jpg)

反编译arg()函数

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P64_5749.jpg)

- 变量a依旧是栈分配，变量b已经逃逸了。
- 在上述代码中的retArg()函数只是个函数声明，没有给出具体实现，通过linkname机制让链接器在链接阶段链接到inner.RetArg()函数。
- retArg()函数只有声明没有实现，而且编译器不会跟踪linkname，所以无法根据代码逻辑判定变量b到底有没有逃逸。

## 3.3 Function Value

### 3.3.1 函数指针

函数指针存储的也是地址, 该地址指向代码段中某个函数的第一条指令

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P65_5768.jpg)

### 3.3.2 Function Value分析

{{< embedcode go "code/3_18/main.go" >}}

反编译

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P67_5870.jpg)

1. 4～7行和最后两行用于栈增长，暂不需要关心。
2. 第8～10行分配栈帧并赋值caller’s BP，RET之前的两行还原BP寄存器并释放栈帧。
3. CALL后面的两行用来复制返回值。
4. CALL连同之前的6条MOVQ指令，实现了Function Value的传参和过程调用。
   1. MOVQ 0x30(SP)，AX和MOVQ AX，0(SP)用于把helper()函数的第2个参数a的值复制给fn()函数的第1个参数。
   2. MOVQ 0x38(SP)，AX和MOVQ AX，0x8(SP)同理，把helper()函数第3个参数b的值复制给fn()函数的第2个参数。
   3. MOVQ 0x28(SP)，DX把helper()函数第1个参数fn的值复制到DX寄存器，MOVQ 0(DX)，AX把DX用作基址，加上位移0，也就是从DX存储的地址处读
      取出一个64位的值，存入了AX寄存器中。
   4. CALL AX说明，上一步中AX寄存器最终存储的是实际函数的地址。

栈分析

```
40(SP) return value -|
38(SP) b             | stack of main
30(SP) a             |
28(SP) fn           -|
20(SP) return addr
18(SP) bp           -|
10(SP) return value  | stack of helper
 8(SP) b             |
 0(SP) a            -|
```

### 3.3.3 闭包

```go
// 第3章 code_3_19.go
func mc(n int) func() int {
	return func() int {
		return n
    }
}
```

闭包的状态保存在哪里呢？

1. 闭包对象

反编译

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P69_5896.jpg)

栈分析

```
...                             |
28(SP) main arg (mc-func()int)  | stack of main
20(SP) main arg (mc-n)         -|
18(SP) return address of mc
10(SP) bp                      -|
 8(SP) newobject ret            | stack of mc
 0(SP) newobject arg           -|
```

推测newobject所创建的对象的结构

```go
// 闭包对象
struct {
	// 闭包函数
    F uintptr
	// 捕获列表
    n int
}
```

2. 看到闭包

newobject的原型

```go
func newobject(typ *_type) unsafe.Pointer
```

使用自定义的newobject实现来查看_type的布局

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P71_5941.jpg)

运行结果

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P72_5963.jpg)

因为`start++`导致start变量逃逸, 所以调用了两次newobject

- `int`
- `struct { F uintptr; start *int }`

图3-12 Function Value和闭包对象

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P73_5970.jpg)

3. 调用闭包

闭包函数在被调用的时候，必须得到当前闭包对象的地址才能访问其中的捕获列表，这个地址是如何传递的呢？

{{< embedcode go "code/3_22/main.go" >}}

反编译

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P74_5996.jpg)

- 将DX寄存器用作基址，再加上位移8，把该地址处的值复制到AX寄存器中。
- 把AX寄存器的值复制给闭包函数的返回值。
- 闭包函数返回。

> 书中说把AX的值给闭包函数的返回值, 不太理解为啥0x8(SP)是返回值地址.

4. 闭包与变量逃逸

```go
// 第3章 code_3_23.go
func sc(n int) int {
	f := func() int {
		return n
    }
	return f()
}
```

禁用内联优化

```shell
$ go build -gcflags='-l'
```

反编译

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P75_6022.jpg)

return f()之前的6行汇编代码

- XORPS和MOVUPS这两行利用128位的寄存器X0，把栈帧上从位移8字节开始的16字节清零，这段区间就是sc()函数的局部变量区，正好符合捕获了一个int变
  量的闭包对象大小。
- LEAQ和MOVQ把闭包函数的地址复制到栈帧上位移8字节处，正是闭包对象中的函数指针。
- 接下来的两个MOVQ把sc()函数的参数n的值复制到栈帧上位移16字节处，也就是闭包捕获列表中的int变量。

图3-13 sc()函数中构造的闭包对象f

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P76_6028.jpg)

return之后的5行汇编代码

- MOVQ把闭包函数的地址复制到AX寄存器中，LEAQ把闭包对象的地址存储到DX寄存器中。
- CALL指令调用闭包函数，接下来的两条MOVQ把闭包函数的返回值复制到sc()函数的返回值。

图3-14 调用闭包函数f()

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P76_6032.jpg)

闭包对象的捕获列表，捕获的是变量的值还是地址？

- 只有在变量的值不会再改变的前提下，才可以复制变量的值，否则就会出现不一致错误。

示例, 需要禁用内联优化

```go
// 第3章 code_3_24.go
// 捕获地址
func sc(n int) int {
	f := func() int {
        n++
        return n
    }
    return f()
}

// 第3章 code_3_25.go
// 捕获值
func sc(n int) int {
	n++
	f := func() int {
		return n
    }
	return f()
}

// 第3章 code_3_26.go
// 捕获地址
func sc(n int) int {
    f := func() int {
        return n
    }
    n++
    return f()
}
```

## 3.4 defer

### 3.4.1 最初的链表

使用go1.12

{{< embedcode go "code/3_28/main.go" >}}

反编译df()

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P79_6112.jpg)

#### deferproc

- Go语言中，每个goroutine都有自己的一个defer链表，而runtime.deferproc()函数做的事情就是把defer函数及其参数添加到链表中。
- 编译器还会在当前函数结尾处插入调用runtime.deferreturn()函数的代码，该函数会按照FILO的顺序调用当前函数注册的所有defer函数。
- 如果当前goroutine发生了panic（宕机），或者调用了runtime.Goexit()函数，runtime的panic处理逻辑会按照FILO的顺序遍历当前goroutine的整
  个defer链表，并逐一调用defer函数，直到某个defer函数执行了recover，或者所有defer函数执行完毕后程序结束运行。

runtime.deferproc()函数原型

```go
func deferproc(size int32, fun *funcval)
```

- Go语言用两级指针结构统一了函数指针和闭包，这个funcval结构就是用来支持两级指针的。
- funcval结构中只定义了uintptr

图3-15 funcval对Function Value两级指针的支持

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P81_6139.jpg)

- 参数siz表示defer函数的参数占用空间的大小，这部分参数也是通过栈传递的，虽然没有出现在deferproc()函数的参数列表里，但实际上会被编译器追加
  到fn的后面
- 注意defer函数的参数在栈上的fn后面，而不是在funcval结构的后面。这点不符合正常的Go语言函数调用约定，属于编译器的特殊处理。

图3-16 df()函数调用deferproc时的栈帧

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P81_6142.jpg)

基于第3章/code_3_28.go反编译得到的汇编代码，整理出等价的伪代码如下：

```go
func df(n int) (v int) {
	r := runtime.deferproc(8, df.func1, &n)
	if r > 0 {
		goto ret
	}
	v = n
	runtime.deferreturn()
	return
ret:
	runtime.deferreturn()
	return
}

func df.func1(i *int) {
	*i *= 2
}
```

deferproc()函数的返回值为0或非0时代表不同的含义

- 0代表正常流程，也就是已经把需要延迟执行的函数注册到了链表中，这种情况下程序可正常执行后续逻辑。
- 返回值为1则表示发生了panic，并且当前defer函数执行了recover，这种情况会跳过当前函数后续的代码，直接执行返回逻辑。

deferproc()函数的具体实现, 摘抄自runtime包的panic.go

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P82_6165.jpg)

通过getcallersp()函数获取调用者的SP，也就是调用deferproc()函数之前SP寄存器的值。这个值有两个用途

- 一是在deferreturn()函数执行defer函数时用来判断该defer是不是被当前函数注册的
- 二是在执行recover的时候用来还原栈指针。

基于unsafe指针运算得到编译器追加在fn之后的参数列表的起始地址，存储在argp中。

通过getcallerpc()函数获取调用者指令指针的位置，在amd64上实际就是deferproc()函数的返回地址，从调用者df()函数的视角来看就是CALL
runtime.deferproc后面的那条指令的地址。这个地址主要用来在执行recover的时候还原指令指针。

调用newdefer()函数分配一个runtime._defer结构，newdefer()函数内部使用了两级缓冲池来避免频繁的堆分配，并且会自动把新分配的_defer结构添加
到链表的头部。

runtime._defer的定义

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P83_6179.jpg)

- siz表示defer参数占用的空间大小，与deferproc()函数的第1个参数一样。
- started表示有个panic或者runtime.Goexit()函数已经开始执行该defer函数。
- _panic的值是在当前goroutine发生panic后，runtime在执行defer函数时，将该指针指向当前的_panic结构。
- link指针用来指向下一个_defer结构，从而形成链表。

_defer中没有发现用来存储defer函数参数的空间，参数应该被存储到哪里？

实际上runtime.newdefer()函数用了和编译器一样的手段，在分配_defer结构的时候，后面额外追加了siz大小的空间，如图3-17所示，所以deferproc()
函数接下来会将fn、callerpc、sp都复制到_defer结构中相应的字段，然后根据siz大小来复制参数，最后通过return0()函数来把返回值0写入AX寄存器中。

图3-17 deferproc执行中为_defer赋值

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P84_6185.jpg)

通过deferproc()函数注册完一个defer函数后，deferproc()函数的返回值是0。后面如果发生了panic，又通过该defer函数成功recover，那么指令指针
和栈指针就会恢复到这里设置的pc、sp处，看起来就像刚从runtime.deferproc()函数返回，只不过返回值为1，编译器插入的if语句继而会跳过函数体，仅
执行末尾的deferreturn()函数。

#### deferreturn

在正常情况下，注册过的defer函数是由runtime.deferreturn()函数负责执行的，正常情况指的就是没有panic或runtime.Goexit()函数，即当前函数完
成执行并正常返回时。

deferreturn()函数的代码如下：

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P84_6192.jpg)

值得注意的是参数arg0的值没有任何含义，实际上编译器并不会传递这个参数，deferreturn()函数内部通过它获取调用者栈帧上args to callee区间的起
始地址，从而可以将defer函数所需参数复制到该区间。defer函数的参数个数要比编译器传给deferproc()函数的参数还少两个，所以调用者的
args to callee区间大小肯定足够，不必担心复制参数会覆盖掉栈帧上的其他数据。

deferreturn()函数的主要逻辑如下：

1. 若defer链表为空，则直接返回，否则获得第1个_defer的指针d，但并不从链表中移除。
2. 判断d.sp是否等于调用者的SP，即判断d是否由当前函数注册，如果不是，则直接返回。
3. 如果defer函数有参数，d.siz会大于0，就将参数复制到栈上&arg0处。
4. 将d从defer链表移除，链表头指向d.link，通过runtime.freedefer()函数释放d。和runtime.newdefer()函数对应，runtime.freedefer()函数
   会把d放回缓冲池中，缓冲池内部按照defer函数参数占用空间的多少分成了5个列表，对于参数太多且占用空间太大的d，超出了缓冲池的处理范围则不会被
   缓存，后续会被GC回收。
5. 通过runtime.jmpdefer()函数跳转到defer函数去执行。

runtime.jmpdefer()函数是用汇编语言实现的，amd64平台下的实现代码如下：

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P85_6210.jpg)

第2行把fn赋值给DX寄存器，3.3节中已经讲过Function Value调用时用DX寄存器传递闭包对象地址。接下来的3行代码通过设置SP和BP来还原
deferreturn()函数的栈帧，结合最后一条指令是跳转到defer函数而不是通过CALL指令来调用，这样从调用栈来看就像是deferreturn()函数的调用者直接
调用了defer函数。

jmpdefer()函数会调整返回地址，在amd64平台下会将返回地址减5，即一条CALL指令的大小，然后才会跳转到defer函数去执行。这样一来，等到defer函数
执行完毕返回的时候，刚好会返回编译器插入的runtime.deferreturn()函数调用之前，从而实现无循环、无递归地重复调用deferreturn()函数。直到当
前函数的所有defer都执行完毕，deferreturn()函数会在第1、第2步判断时返回，不经过jmpdefer()函数调整栈帧和返回地址，从而结束重复调用。

使用deferproc()函数实现defer的好处是通用性比较强，能够适应各种不同的代码逻辑。

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P86_6220.jpg)

因为defer函数的注册是运行阶段才进行的，可以跟代码逻辑很好地整合在一起，所以像if这种条件分支不用完成额外工作就能支持。由于每个
runtime._defer结构都是基于缓冲池和堆动态分配的，所以即使不定次数的循环也不用额外处理，多次注册互不干扰。

但是链表与堆分配组合的最大缺点就是慢，即使用了两级缓冲池来优化runtime._defer结构的分配，性能方面依然不太乐观，所以在后续的版本中就开始了对
defer的优化之旅。

### 3.4.2 栈上分配

在1.13版本中对defer做了一点小的优化，即把runtime._defer结构分配到当前函数的栈帧上。很明显这不适用于循环中的defer，循环中的defer仍然需要
通过deferproc()函数实现，这种优化只适用于只会执行一次的defer。

编译器通过runtime.deferprocStack()函数来执行这类defer的注册，相比于runtime.deferproc()函数，少了通过缓冲池或堆分配_defer结构的步骤，
性能方面还是稍有提升的。

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P87_16307.jpg)

runtime._defer结构中新增了一个bool型的字段heap来表示是否为堆上分配，对于这种栈上分配的_defer结构，deferreturn()函数就不会用
freedefer()函数进行释放了。因为编译器在栈帧上已经把_defer结构的某些字段包括后面追加的fn的参数都准备好了，所以deferprocStack()函数这里只
需为剩余的几个字段赋值，与deferproc()函数的逻辑基本一致。最后几行中通过unsafe.Pointer做类型转换再赋值，源码注释中解释为避免写屏障，暂时理
解成为提升性能就行了

同样使用第3章/code_3_28.go，经过Go 1.13编译器转换后的伪代码如下：

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P87_6255.jpg)

图3-18 df()函数调用deferprocStack()时的栈帧

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P88_6273.jpg)

栈上分配_defer这种优化只是节省了_defer结构的分配、释放时间，仍然需要将defer函数添加到链表中，在调用的时候也还要复制栈上的参数，整体提升比
较有限。

### 3.4.3 高效的open coded defer

在Go 1.14版本中又进行了一次优化，这次优化也是针对那些只会执行一次的defer。编译器不再基于链表实现这类defer，而是将这类defer直接展开为代码中
的函数调用，按照倒序放在函数返回前去执行，这就是所谓的open coded defer。

使用第3章/code_3_28.go，在1.14版本中经编译器转换后的伪代码如下：

```go
func df(n int) (v int) {
	v = n
	func(i *int) {
		*i *= 2
    }(&n)
	return
}
```

两个问题：

- 如何支持嵌套在if语句块中的defer？
- 当发生panic时，如何保证这些defer得以执行呢？

第1个问题其实并不难解决，可以**在栈帧上分配一个变量**，用每个二进制位来记录一个对应的defer函数是否需要被调用。Go语言实际上用了一字节作为标
志，可以最多支持8个defer，为什么不支持更多呢？笔者是这样理解的，open coded defer本来就是为了提高性能而设计的，一个函数中写太多defer，应该
是不太在意这种层面上的性能了。

还需要考虑的一个问题是，deferproc()函数在注册的时候会存储defer函数的参数副本，defer函数的参数经常是当前函数的局部变量，即使它们后来被修改
了，deferproc()函数存储的副本也是不会变的，副本是注册那一时刻的状态，所以在open coded defer中编译器需要在当前函数栈帧上分配额外的空间来存
储defer函数的参数。

示例

{{< embedcode go "code/3_30/main.go" >}}

经编译器转换后的等价代码如下：

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P90_6312.jpg)

其中局部变量f就是专门用来支持if这类条件逻辑的标志位，局部变量i用作n在defer注册那一刻的副本，函数返回前根据标志位判断是否调用defer函数。

图3-19 fn()函数通过open coded defer的方式调用defer函数

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P90_6329.jpg)

## 3.5 panic
