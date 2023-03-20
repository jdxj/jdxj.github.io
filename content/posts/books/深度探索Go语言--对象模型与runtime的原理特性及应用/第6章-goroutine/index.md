---
title: "第6章 Goroutine"
date: 2023-02-19T13:24:14+08:00
summary: 6.2 IO多路复用写的不是特别详细
---

## 6.1 进程、线程与协程

### 6.1.1 进程

现代操作系统利用硬件提供的页表机制，通过为不同进程分配独立的页表，实现进程间地址空间的隔离。

图6-1 进程间地址空间的隔离

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P194_9294.jpg)

Linux通过clone系统调用来创建新的进程。

### 6.1.2 线程

为什么要有一个用户栈和一个内核栈呢？

- 因为我们的线程在执行过程中经常需要在用户态和内核态之间切换，通过系统调用进入内核态使用系统资源。
- 对于内核来讲，任何的用户代码都被视为不安全的，可能有Bug或者带有恶意的代码，所以操作系统不允许用户态的代码访问内核数据。

图6-2 线程的用户栈和内核栈

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P195_9306.jpg)

调度系统切换线程时，如果两个线程属于同一个进程，开销要比属于不同进程时小得多

- 因为不需要切换页表，相应地，TLB缓存也就不会失效。
- 同一个进程中的多个线程，因为共享同一个虚拟地址空间，所以线程间数据共享变得十分简单高效

图6-3 同进程间线程切换

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P195_9310.jpg)

在线程切换频繁时，调度本身的开销会占用大量CPU资源，造成系统吞吐量严重下降。

### 6.1.3 协程

从具体实现来看，纤程就是一个由入口函数地址、参数和独立的用户栈组成的任务，相当于让线程可以有多个用户栈

图6-4 纤程概念示意图

## 6.2 IO多路复用

### 6.2.1 3种网络IO模型

把一个常见的TCP socket的recv请求分成两个阶段：一是等待数据阶段，等待网络数据就绪；二是数据复制阶段，把数据从内核空间复制到用户空间。

对于阻塞式IO来讲，整个IO过程是一直阻塞的，直至这两个阶段都完成。

图6-5 经典的阻塞式网络IO模型

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P198_9335.jpg)

如果想要启用非阻塞式IO，需要在代码中使用fcntl()函数将对应socket的描述符设置成O_NONBLOCK模式。

- 在非阻塞模式下，线程等待数据的时候不会阻塞，从编程角度来看就是recv()函数会立即返回，并返回错误代码EWOULDBLOCK（某些平台的SDK也可能是
  EAGAIN），表明此时数据尚未就绪，可以先去执行别的任务。
- 程序一般会以合适的频率重复调用recv()函数，也就是进行轮询操作。在数据就绪之前，recv()函数会一直返回错误代码EWOULDBLOCK。
- 等到数据就绪后，再进入复制数据阶段，从内核空间到用户空间。
- **因为非阻塞模式下的数据复制也是同步进行的，所以可以认为第二阶段也是阻塞的**。

图6-6 非阻塞式网络IO模型

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P199_9341.jpg)

有了非阻塞式IO是不是就万事大吉了呢？

- 虽然第一阶段不会阻塞，但是需要频繁地进行轮询。一次轮询就是一次系统调用，如果轮询的频率过高就会空耗CPU，造成大量的额外开销
- 如果轮询频率过低，就会造成数据处理不及时，进而使任务的整体耗时增加。

IO多路复用技术就是为解决上述问题而诞生的

- 与非阻塞式IO相似，从socket读写数据不会造成线程挂起。
- 在此基础之上把针对单个socket的轮询改造成了批量的poll操作，可以通过设置超时时间选择是否阻塞等待。
- 只要批量socket中有一个就绪了，阻塞挂起的线程就会被唤醒，进而去执行后续的数据复制操作。

图6-7 IO多路复用

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P199_9344.jpg)

### 6.2.2 示例对比

## 6.3 巧妙结合

把每个网络请求放到一个单独的协程中去处理，底层的IO事件循环在处理不同的socket时直接切换到与之关联的协程栈

图6-10 协程与IO多路复用的结合

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P204_9398.jpg)

这样一来，就把IO事件循环隐藏到了runtime内部，开发者可以像阻塞式IO那样平铺直叙地书写代码逻辑，尽情地把数据存放在栈帧上的局部变量中，代码执行
网络IO时直接触发协程切换，切换到下一个网络数据已经就绪的协程。当底层的IO事件循环完成本轮所有协程的处理后，再次执行netpoll，如此循环往复，开
发者不会有任何感知，程序却得以高效执行。

## 6.4 GMP模型

### 6.4.1 基本概念

G指的就是goroutine；M是Machine的缩写，指的是工作线程；P则是指处理器Processor，代表了一组资源，M要想执行G的代码，必须持有一个P才行。

### 6.4.2 从GM到GMP

在早期版本的Go实现中（1.1版本之前），是没有P的，只有G和M

图6-11 GM调度模型

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P205_9408.jpg)

GM调度模型有几个明显的问题

- 用一个全局的mutex保护着一个全局的runq（就绪队列），所有goroutine的创建、结束，以及调度等操作都要先获得锁，造成对锁的争用异常严重。
- G的每次执行都会被分发到随机的M上，造成在不同M之间频繁切换，破坏了程序的局部性
- 每个M都会关联一个内存分配缓存mcache，造成了大量的内存开销，进一步使数据的局部性变差。
- 在存在系统调用的情况下，工作线程经常被阻塞和解除阻塞，从而增加了很多开销。

为了解决上述这些问题，新的调度器被设计出来。

- 总体的优化思路就是将处理器P的概念引入runtime，并在P之上实现工作窃取调度程序。
- M仍旧是工作线程，P表示执行Go代码所需的资源。当一个M在执行Go代码时，它需要有一个关联的P，当M执行系统调用或者空闲时，则不需要P。

图6-12 GMP调度模型

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P205_9412.jpg)

1. 本地runq和全局runq

- 当一个G从等待状态变成就绪状态后，或者新创建了一个G的时候，这个G会被添加到当前P的本地runq。
- 当M执行完一个G后，它会先尝试从关联的P的本地runq中取下一个，如果本地runq为空，则到全局runq中去取
- 如果全局runq也为空，就会去其他的P那里窃取一半的G过来。

图6-13 本地runq为空到全局runq获取G

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P206_9419.jpg)

图6-14 全局runq也为空窃取其他P的G

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P206_9422.jpg)

2. M的自旋

- 当一个M进入系统调用时，它必须确保有其他的M来执行Go代码。新的调度器设计引入了一定程度的自旋，就不用再像之前那样过于频繁地挂起和恢复M了，这
  会多消耗一些CPU周期，但是对整体性能的影响是正向的。
- 自旋分两种：第一种是一个有关联P的M，自旋寻找可执行的G；第二种是一个没有P的M，自旋寻找可用的P。这两种自旋的M的个数之和不超过GOMAXPROCS，
  当存在第二种自旋的M时，第一种自旋的M不会被挂起(因为P不够)。
- 当一个新的G被创建出来或者M即将进行系统调用，或者M从空闲状态变成忙碌状态时，它会确保至少有一个处于自旋状态的M（除非所有的P都忙碌），这样保
  证了处于可执行状态的G都可以得到调度，同时还不会频繁地挂起、恢复M。

## 6.5 GMP主要数据结构

### 6.5.1 runtime.g

runtime.g部分字段

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P207_9434.jpg)

表6-1 runtime.g部分字段的用途

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-T208_16088.jpg)

stack是个结构体类型, 用来描述goroutine的栈空间的，对应的内存区间是一个左闭右开区间[lo，hi]。

```go
type stack struct {
	lo uintptr
	hi uintptr
}
```

sched用来存储goroutine执行上下文, 它与goroutine协程切换的底层实现直接相关，其对应的gobuf结构代码如下

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P208_9539.jpg)

- sp字段存储的是栈指针
- pc字段存储的是指令指针
- g用来反向关联到对应的G
- ctxt指向闭包对象，也就是说用go关键字创建协程的时候传递的是一个闭包，这里会存储闭包对象的地址
- ret用来存储返回值，实际上是利用AX寄存器实现类似C函数的返回值，目前只发现panic-recover机制用到了该字段。lr在arm等架构上用来存储返回地址，
  x86没有用到该字段
- bp用来存储栈帧基址。

atomicstatus描述了当前G的状态

表6-2 atomicstatus的取值及其含义

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-T209_16090.jpg)

waiting对应的sudog结构

### 6.5.2 runtime.m

runtime.m部分字段

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P210_9614.jpg)

表6-3 runtime.m部分字段的用途

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-T210_16409.jpg)

### 6.5.3 runtime.p

runtime.p部分字段

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P211_9724.jpg)

表6-4 runtime.p各个字段的主要用途

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-T211_16413.jpg)

status字段有5种不同的取值，分别表示P所处的不同状态

表6-5 P的不同状态

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-T212_16096.jpg)

### 6.5.4 schedt

Go 1.16版源代码中的schedt结构定义

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P212_9888.jpg)
![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P213_9897.jpg)

表6-6 schedt部分字段的主要用途

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-T214_16098.jpg)

