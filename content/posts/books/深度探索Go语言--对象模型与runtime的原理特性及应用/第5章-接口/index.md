---
title: "第5章 接口"
date: 2023-02-15T21:29:46+08:00
draft: true
---

## 5.1 空接口

是指不包含任何方法的接口interface{}

### 5.1.1 一个更好的void∗

如果用unsafe.Sizeof()函数获取一个interface{}类型变量的大小，在64位平台上是16字节，在32位平台上是8字节。interface{}类型本质上是个
struct，由两个指针类型的成员组成，在runtime中可以找到对应的struct定义

```go
type eface struct {
	_type *_type
	data unsafe.Pointer
}
```

还有一个专门的类型转换函数efaceOf()，该函数接受的参数是一个interface{}类型的指针，返回值是一个eface类型的指针，内部实际只进行了一下指针类
型的转换，也就说明interface{}类型在内存布局层面与eface类型完全等价。

```go
func efaceOf(ep *interface{}) *eface {
	return (*eface)(unsafe.Pointer(ep))
}
```

- data字段是一个unsafe.Pointer类型的指针，用来存储实际数据的地址。
- unsafe.Pointer在含义上和C语言中的void∗有些类似，只用来表明这是一个指针，并不限定指向的目标数据的类型，可以接受任意类型的地址。
- _type字段用来描述data的类型元数据

```go
// 第5章 code_5_1.go
var n int
var e interface{} = &n
```

图5-1 空接口变量e与赋值变量n的关系

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P124_7145.jpg)

就变量n本身而言，它的类型信息只会被编译器使用，编译阶段参考这种类型信息来分配存储空间、生成机器指令，但是并不会把这种类型信息写入最终生成的
可执行文件中。从内存布局的角度来讲，变量n在64位和32位平台分别占用8字节和4字节，占用的这些空间全部用来存放整型的值，没有任何空间被用来存放整
型类型信息。

把变量n的地址赋值给interface{}类型的变量e的这个操作，意味着编译器要把∗int的类型元数据生成出来，并把其地址赋给变量e的_type字段，这些类型元
数据会被写入最终的可执行文件

```go
// 第5章 code_5_2.go
func p2e(p *int) (e interface{}) {
	e = p
	return
}
```

反编译p2e()

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P125_7163.jpg)

等价伪代码

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P125_7174.jpg)

对于interface{}类型的变量e，它的声明类型是interface{}, _type会随着变量e装载不同类型的数据而发生改变，所以后文中将它称为变量e的动态类型，
并相应地把变量e的声明类型称为静态类型。

### 5.1.2 类型元数据

在C语言中类型信息主要存在于编译阶段，编译器从源码中得到具体的类型定义，并记录到相应的内存数据结构中，然后根据这些类型信息进行语法检查、生成
机器指令等。例如x86整数加法和浮点数加法采用完全不同的指令集，编译器根据数据的类型来选择。这些类型信息并不会被写入可执行文件，即使作为符号数
据被写入，也是为了方便调试工具，并不会被语言本身所使用。

Go与C语言不同的是，在设计之初就支持面向对象编程，还有其他一些动态语言特征，这些都要求运行阶段能够获得类型信息，所以语言的设计者就把类型信息
用统一的数据结构来描述，并写入可执行文件中供运行阶段使用，这就是所谓的类型元数据。

Go 1.15版本的runtime源码

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P126_7187.jpg)

表5-1 _type各字段的含义及主要用途

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-T126_16339.jpg)

_type提供了适用于所有类型的最基本的描述，对于一些更复杂的类型，例如复合类型slice和map等，runtime中分别定义了maptype、slicetype等对应的
结构。

```go
type slicetype struct {
	typ _type
	elem *_type
}
```

Go语言允许为自定义类型实现方法，这些方法的相关信息也会被记录到自定义类型的元数据中，一般称为类型的方法集信息。

```go
type Integer int
```

图5-2 自定义类型Integer的类型元数据结构

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P127_7285.jpg)

_type结构的tflag字段是几个标志位，当tflagUncommon这一位为1时，表示类型为自定义类型。从runtime的源码可以发现，_type类型有一个
uncommon()方法，对于自定义类型可以通过此方法得到一个指向uncommontype结构的指针

uncommontype结构的定义代码如下：

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P127_7293.jpg)

- 通过pkgpath可以知道定义该类型的包名称
- mcount表示该类型共有多少个方法
- xcount表示有多少个方法被导出
- moff是个偏移值，那里就是方法集的元数据，也就是一组method结构构成的数组。

例如，若为自定义类型Integer定义两个方法，它的类型元数据及其method数组的内存布局如图5-3所示。

图5-3 Integer类型元数据及其method数组的内存布局

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P128_7299.jpg)
![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P128_7306.jpg)

> 可以将type method struct理解为tcp封包中的header.

图5-4 method数组排序

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P129_7314.jpg)

- 通过name偏移能够找到方法的名称字符串
- mtyp偏移处是方法的类型元数据，进一步可以找到参数和返回值相关的类型元数据。
- ifn是供接口调用的方法地址
- tfn是正常的方法地址，这两个方法地址有什么不同呢？ifn的接收者类型一定是指针，而tfn的接收者类型跟源代码中的实现一致

以上这些类型元数据都是在编译阶段生成的，经过链接器的处理后被写入可执行文件中，runtime中的类型断言、反射和内存管理等都依赖于这些元数据.

### 5.1.3 逃逸与装箱
