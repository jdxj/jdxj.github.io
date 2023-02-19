---
title: "第5章 接口"
date: 2023-02-15T21:29:46+08:00
draft: false
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

data字段是个指针，那么它是如何接收来自一个值类型的赋值的呢？

```go
// 第5章 code_5_3.go
n := 10
var e interface{} = n
```

图5-5 interface{}类型的变量e的数据结构

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P129_7326.jpg)

把第5章/code_5_3.go放到一个函数中

```go
// 第5章 code_5_4.go
func v2e(n int) (e interface{}) {
	e = n
	return
}
```

反汇编

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P130_7344.jpg)

等价的伪代码

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P130_7352.jpg)

runtime.convT64()函数的源代码

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P131_7372.jpg)

staticuint64s是个长度为256的uint64数组，每个元素的值都跟下标一致，存储了0～255这256个值，主要用来避免常用数字频繁地进行堆分配。

图5-6 staticuint64s数组

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P131_7379.jpg)

图5-7 变量e的数据结构

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P132_7385.jpg)

值类型装箱就一定会进行堆分配吗？

```go
// 第5章 code_5_5.go
func fn(n int) bool {
	return notNil(n)
}

func notNil(a interface{}) bool {
	return a != nil
}
```

编译时需要禁止内联优化，编译器还能够通过notNil()函数的代码实现判定有没有发生逃逸，反编译fn()函数得到的汇编代码如下

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P132_7401.jpg)

转换为等价的伪代码

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P133_7421.jpg)

图5-8 fn()函数的调用栈

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P133_7414.jpg)

注意局部变量v，它实际上是被编译器采用隐式方式分配的，被用作变量n的值的副本，却并没有分配到堆上。

**interface{}在装载值的时候必须单独复制一份，而不能直接让data存储原始变量的地址，因为原始变量的值后续可能会发生改变，这就会造成逻辑错误。**

## 5.2 非空接口

### 5.2.1 动态派发

多态

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P134_7438.jpg)

1. 方法地址静态绑定

要进行方法（函数）调用，有两点需要确定：

- 一是方法的地址，也就是在代码段中的指令序列的起始地址；
- 二是参数及调用约定，也就是要传递什么参数及如何传递的问题（通过栈或者寄存器），返回值的读取也包含在调用约定范畴内。

不使用接口而直接通过自定义类型的对象实例调用其方法的例子

```go
//go:noinline
func ReadFile(f *os.File, b[]byte) (n int, err error) {
	return f.Read(b)
}
```

反编译

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P135_7462.jpg)

伪代码

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P136_7479.jpg)

从汇编语言的角度来看，上述方法的调用是通过CALL指令＋相对地址实现的，方法地址在可执行文件构建阶段就确定了，一般将这种情况称为
**方法地址的静态绑定**。

对于动态派发来讲，编译阶段能够确定的是要调用的方法的名字，以及方法的原型（参数与返回值列表）。

2. 动态查询类型元数据

让我门设计动态派发

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P137_7492.jpg)

例子代码

```go
var r io.Reader = f
n, err := r.Read(buf)
```

- 首先，可以通过变量r得到∗os.File的类型元数据
- 然后根据方法名称Read以二分法查找匹配的method结构
- 找到后再根据method.mtyp得到方法本身的类型元数据
- 最后对比方法原型是否一致（参数和返回值的类型、顺序是否一致）。
- 如果原型一致，就找到了目标方法，通过method.ifn字段得到方法的地址，然后就像调用普通函数一样调用就可以了。

图5-9 ∗os.File的类型元数据

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P137_7504.jpg)

单就动态派发而言，这种方式确实可以实现，但是有一个明显的问题，那就是效率低，或者说性能差。

- 跟地址静态绑定的方法调用比起来，原本一条CALL指令完成的事情，这里又多出了一次二分查找加方法原型匹配，增加的开销不容小觑，可能会造成动态派发
  的方法比静态绑定的方法多一倍开销甚至更多，所以必须进行优化。
- 不能在每次方法调用前都到元数据中去查找，尽量做到一次查找、多次使用，这里可以一定程度上参考C++的虚函数表实现。

3. C++虚函数机制

C++中的虚函数机制跟接口的思想很相似，编程语言允许父类指针指向子类对象，当通过父类的指针来调用虚函数时，就能实现动态派发。

具体实现原理就是

- 编译器为每个包含虚函数的类都生成一张虚函数表，实际上是个地址数组，按照虚函数声明的顺序存储了各个虚函数的地址。
- 此外还会在类对象的头部安插一个虚指针（GCC安插在头部，其他编译器或有不同），指向类型对应的虚函数表。
- 运行阶段通过类对象指针调用虚函数时，会先取得对象中的虚指针，进一步找到对象类型对应的虚函数表，然后基于虚函数声明的顺序，以数组下标的方式从
  表中取得对应函数的地址，这样整个动态派发过程就完成了。

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P138_7515.jpg)

测试代码

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P139_7527.jpg)

输出

```
A,8
B,16
```

图5-10 C++虚函数动态派发示例

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P139_7564.jpg)

参考C++的虚函数表思想，再回过头来看Go语言中接口的设计，如果把这种基于数组的函数地址表应用在接口的实现中，基本就能消除每次查询地址造成的性能
开销。显然这里需要对eface结构进行扩展，加入函数地址表相关字段，经过扩展的eface姑且称作efacex

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P139_7572.jpg)

> 怪不得叫`tab` (table)

图5-11 参照C++虚函数机制修改后的非空接口数据结构

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P140_7585.jpg)

什么时候为fun数组赋值呢？当然是在为整个efacex结构赋值的时候最合适

```go
// 第5章 code_5_9.go
f, _ := os.Open("gom.go")
var rw io.ReadWriter
rw = f
```

从f到rw这个看似简单的赋值，至少要展开成如下几步操作：

- ①根据rw接口中方法的个数动态分配tab结构，这里有两个方法，fun数组的长度是2。
- ②从∗os.File的方法集中找到Read()方法和Write()方法，把地址写入fun数组对应下标。
- ③把∗os.File的元数据地址赋值给tab._type。
- ④把f赋值给data，也就是数据指针。

图5-12 基于efacex设计的非空接口变量rw赋值后的数据结构

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P141_7599.jpg)

实际上，fun数组也不用每次都重新分配和初始化，从指定具体类型到指定接口类型变量的赋值，运行阶段无论发生多少次，每次生成的fun数组都是相同的。例
如从∗os.File到io.ReadWriter的赋值，每次都会生成一个长度为2的fun数组，数组的两个元素分别用于存储(∗os.File).Read和(∗os.File).Write的
地址。也就是说通过一个确定的接口类型和一个确定的具体类型，就能够唯一确定一个fun数组，因此可以通过一个全局的map将fun数组进行缓存，这样就能进
一步减少方法集的查询，从而优化性能。

### 5.2.2 具体实现

实际上在Go语言的runtime中与非空接口对应的结构类型是iface

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P141_7608.jpg)

因为也是通过数据指针data来装载数据的，所以也会有逃逸和装箱发生。其中的itab结构就包含了具体类型的元数据地址_type，以及等价于虚函数表的方法地
址数组fun，**除此之外还包含了接口本身的类型元数据地址inter**

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P141_7616.jpg)

1. 接口类型元数据

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P142_7627.jpg)

除去最基本的typ字段，pkgpath表示接口类型被定义在哪个包中，mhdr是接口声明的方法列表。

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P142_7635.jpg)

比自定义类型的method结构少了方法地址，只包含方法名和类型元数据的偏移。

- 这些**偏移**的实际类型为int32，与指针的作用一样，但是64位平台上比使用指针节省一半空间。
- 以ityp为起点，可以找到方法的参数（包括返回值）列表，以及每个参数的类型信息，也就是说这个ityp是方法的原型信息。

图5-13 io.ReadWriter类型的变量rw的数据结构

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P143_7657.jpg)

2. 如何获得itab

运行阶段可通过runtime.getitab函数来获得相应的itab，该函数被定义在runtime包中的iface.go文件中

```go
func getitab(inter *interfacetype, typ *_type, canfail bool) *itab
```

前两个参数inter和typ分别是接口类型和具体类型的元数据，canfail表示是否允许失败。如果typ没有实现inter要求的所有方法，则canfail为true时函
数返回nil，canfail为false时就会造成panic。对应到具体的语法就是comma ok风格的类型断言和普通的类型断言

```go
r, ok := a.(io.Reader) // comma ok
r := a.(io.Reader) //有可能造成panic
```

getitab()函数的代码摘抄自Go语言runtime源码

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P143_7665.jpg)

函数的主要逻辑如下：

- ①校验inter的方法列表长度不为0，为没有方法的接口生成itab是没有意义的。
- ②通过typ.tflag标志位来校验typ为自定义类型，因为只有自定义类型才能有方法集。
- ③在不加锁的前提下，以inter和typ作为key查找itab缓存itabTable，找到后就跳转到⑤。
- ④加锁后再次查找缓存，如果没有就通过persistentalloc()函数进行持久化分配，然后初始化itab并调用itabAdd添加到缓存中，最后解锁。
- ⑤通过itab的fun[0]是否为0来判断typ是否实现了inter接口，如果没实现，则根据canfail决定是否造成panic，若实现了，则返回itab地址。

判断itab.fun[0]是否为零，也就是判断第一个方法的地址是否有效，因为Go语言会把无效的itab也缓存起来，主要是为了避免缓存穿透。缓存中查不到对应
的itab，就会每次都查询元数据的方法列表，从而显著影响性能，所以Go语言会把有效、无效的itab都缓存起来，通过fun[0]加以区分。

> fun[0] 相当于标记.

图5-14 interfacetype和_type与itab的对应关系

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P145_7680.jpg)

3. itab缓存

itabTable就是runtime中itab的全局缓存，它本身是个itabTableType类型的指针

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P145_7688.jpg)

- entries是实际的缓存空间
- size字段表示缓存的容量，也就是entries数组的大小
- count表示实际已经缓存了多少个itab。

entries的初始大小是通过itabInitSize指定的，这个常量的值为512。当缓存存满以后，runtime会重新分配整个struct，entries数组是
itabTableType的最后一个字段，可以无限增大它的下标来使用超出容量大小的内存，只要在struct之后分配足够的空间就够了，这也是C语言里常用的手法。

itabTableType被实现成一个散列表。查找和插入操作使用的key是由接口类型元数据与动态类型元数据组合而成的，哈希值计算方式为接口类型元数据哈希值
inter.typ.hash与动态类型元数据哈希值typ.hash进行异或运算。

图5-15 itabTableType哈希表

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P146_7695.jpg)

方法find()和add()分别负责实现itabTableType的查找和插入操作，方法add()操作内部不会扩容存储空间，重新分配操作是在外层实现的，因此
**对于find()方法而言，已经插入的内容不会再被修改，所以查找时不需要加锁**。方法add()操作需要在加锁的前提下进行，getitab()函数是通过调用
itabAdd()函数来完成添加缓存的，itabAdd()函数内部会按需对缓存进行扩容，然后调用add()方法。因为缓存扩容需要重新分配itabTableType结构，为
了并发安全，使用原子操作更新itabTable指针。**加锁后立刻再次查询也是出于并发的考虑，避免其他协程已经将同样的itab添加至缓存**。

通过persistentalloc()函数分配的内存不会被回收

itab类型的init方法

- init()函数内部就是遍历接口的方法列表和具体类型的方法集，来寻找匹配的方法的地址。
- 虽然遍历操作使用了两层嵌套循环，但是方法列表和方法集都是有序的，两层循环实际上都只需执行一次。
- 匹配方法时还会考虑方法是否导出，以及接口和具体类型所在的包。如果是导出的方法则直接匹配成功，如果方法未导出，则接口和具体类型需要定义在同一
  个包中，方可匹配成功。
- 最后需要再次强调的是，对于匹配成功的方法，地址取的是method结构中的ifn字段

### 5.2.3 接收者类型

具体类型方法元数据中的ifn字段，该字段存储的是专门供接口使用的方法地址。所谓专门供接口使用的方法，实际上就是个接收者类型为指针的方法。

还记不记得第4章中分析OBJ文件时，发现**编译器总是会为每个值接收者方法包装一个指针接收者方法**？这也就说明，接口是不能直接使用值接收者方法的，
这是为什么呢？

5.2.2节已经看过了接口的数据结构iface，它包含一个itab指针和一个data指针，data指针存储的就是数据的地址。对于接口来讲，在调用指针接收者方法
时，传递地址是非常方便的，也不用关心数据的具体类型，地址的大小总是一致的。假如通过接口调用值接收者方法，就需要通过接口中的data指针把数据的值
复制到栈上，由于编译阶段不能确定接口背后的具体类型，所以编译器不能生成相关的指令来完成复制，进而无法调用值接收者方法。

如果基于reflectcall()函数，能不能实现通过接口调用值接收者方法呢？

- 肯定是可以实现的，接口的itab中有具体类型的元数据，确实能够应用reflectcall()函数
- 但是有个明显的问题，那就是性能太差。跟几条用于传参的MOV指令加一条普通的CALL指令相比，reflectcall()函数的开销太大了，所以Go语言选择为值
  接收者方法生成包装方法。
- 对于代码中的值接收者方法，类型元数据method结构中的ifn和tfn的值是不一样的，指针接收者方法的ifn和tfn是一样的。

从类型元数据来看，T和∗T是不同的两种类型。

- 接收者类型为T的所有方法，属于T的方法集。
- 因为编译器自动包装指针接收者方法的关系，∗T的方法集包含所有方法，也就是所有接收者类型为T的方法加上所有接收者类型为∗T的方法。

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P147_7712.jpg)
![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P148_7722.jpg)
![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P149_7736.jpg)
![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P149_7749.jpg)

- 第1行输出打印出了Integer类型的方法集，String()和Value()这两个方法各自的IFn和TFn都不相等，这是因为IFn指向接收者为指针类型的方法代码，
  而TFn指向接收者为值类型的方法代码。
- 第2行输出打印出了∗Integer类型的方法集，这两个方法各自的IFn和TFn是相等的，都与第1条指令中同名方法的IFn的值相等。
- 第3行输出打印出了Number接口itab中fun数组中的两个方法地址，与第1行输出Integer方法集中对应方法的IFn的值一致。

图5-16 Integer和∗Integer类型的方法集

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P150_7755.jpg)

### 5.2.4 组合式继承

从方法集的角度进行分析

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P150_7765.jpg)

看一下B、C、∗B和∗C会继承哪些方法

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P151_7786.jpg)

表5-2 示例程序中各自定义类型包含的方法的情况

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-T152_16082.jpg)

Go语言不允许为T和∗T定义同名方法，实际上并不是因为不支持函数重载，前面已经看到了A.Value()方法和(∗A).Value()方法是可以区分的。其根本原因就
是编译器要为值接收者方法生成指针接收者包装方法，要保证两者的逻辑一致，所以不允许用户同时实现，用户可能会实现成不同的逻辑。

## 5.3 类型断言

### 5.3.1 E To具体类型

```go
func normal(a interface{}) int {
	return a.(int)
}
```

反编译

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P153_7899.jpg)

等价的伪代码

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P154_7912.jpg)

comma ok风格的断言

```go
func commaOk(a interface{}) (n int, ok bool) {
	n, ok = a.(int)
	return
}
```

反编译

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P154_7931.jpg)

等价的伪代码

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P155_7951.jpg)

从interface{}到具体类型的断言基本上就是一个指针比较操作加上一个具体类型相关的复制操作

图5-17 从interface{}到具体类型的断言

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P155_7955.jpg)

### 5.3.2 E To I

```go
func normal(a interface{}) io.ReadWriter {
	return a.(io.ReadWriter)
}
```

反编译

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P156_7974.jpg)

伪代码

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P156_7985.jpg)

runtime.assertE2I()函数代码

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P157_7999.jpg)

1. 函数先校验了E的具体类型元数据指针不可为空，没有具体类型的元数据是无法进行断言的
2. 然后通过调用getitab()函数来得到对应的itab，data字段直接复制。
3. 注意调用getitab()函数时最后一个参数为false，根据之前的源码分析已知这个参数是canfail。canfail为false时，如果t没有实现inter要求的所有
   方法，getitab()函数就会造成panic。

comma ok风格的断言

```go
func commaOk(a interface{}) (i io.ReadWriter, ok bool) {
	i, ok = a.(io.ReadWriter)
	return
}
```

反编译

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P157_8015.jpg)

伪代码

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P158_8035.jpg)

runtime.assertE2I2()函数代码

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P158_8046.jpg)

E To I形式的类型断言，主要通过runtime中的assertE2I()和assertE2I2()这两个函数实现，底层的主要任务如图5-18所示，都是通过getitab()函数
完成的方法集遍历及itab分配和初始化。因为getitab()函数中用到了全局的itab缓存，所以性能方面应该也是很高效的。

图5-18 从interface{}到非空接口的类型断言

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P159_8052.jpg)

### 5.3.3 I To具体类型

```go
func normal(i io.ReadWriter) *os.File {
	return i.(*os.File)
}
```

反汇编

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P159_8069.jpg)

伪代码

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P160_8089.jpg)

- 其中的go.itab.∗os.File，io.ReadWriter指的就是全局itab缓存中与∗os.File和io.ReadWriter这一对类型对应的itab。这个itab是在编译阶段就
  被编译器生成的，所以代码中可以直接链接到它的地址。
- 这个断言的核心逻辑就是比较iface中tab字段的地址是否与目标itab地址相等。如果不相等就调用panicdottypeI，如果相等就把iface的data字段返回。
- 注意这里因为∗os.File是指针类型，所以不涉及自动拆箱，也就没有与具体类型相关的复制操作，如果具体类型为值类型就不然了。

comma ok风格的断言

```go
func commaOk(i io.ReadWriter) (f *os.File, ok bool) {
	f, ok = i.(*os.File)
	return
}
```

反编译

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P161_8117.jpg)

伪代码

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P162_8130.jpg)

I To具体类型的断言与E To具体类型的断言在实现上极其相似，核心逻辑如图5-19所示，都是一个指针的相等判断。

图5-19 从非空接口到具体类型的类型断言

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P162_8134.jpg)

### 5.3.4 I To I

```go
func normal(rw io.ReadWriter) io.Reader {
	return rw.(io.ReadWriter)
}
```

反编译

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P163_8154.jpg)

伪代码

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P163_8165.jpg)

runtime.assertI2I()函数代码

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P164_8178.jpg)

1. 先校验i.tab不为nil，否则就意味着没有类型元数据，类型断言也就无从谈起
2. 然后检测i.tab.inter是否等于inter，相等就意味着源接口和目标接口类型相同，直接复制就可以了。
3. 最后才调用getitab()函数，根据inter和i.tab._type获取对应的itab。canfail参数为false，所以如果getitab()函数失败就会造成panic。

comma ok风格的断言

```go
func commaOk(rw io.ReadWriter) (r io.Reader, ok bool) {
	r, ok = rw.(io.Reader)
	return
}
```

反编译

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P164_8194.jpg)

伪代码

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P165_8214.jpg)

runtime.assertI2I2()函数代码

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P165_8225.jpg)

- 如果i.tab为nil，则直接返回false。
- 只有在i.tab.inter与inter不相等时才调用getitab()函数，而且canfail为true，如果getitab()函数失败，则不会造成panic，而是返回nil。

I To I的类型断言，实际上是通过runtime.assertI2I()函数和runtime.assertI2I2()函数实现的，底层也都是基于getitab()函数实现的。

图5-20 从非空接口到非空接口的类型断言

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P166_8239.jpg)

## 5.4 反射

### 5.4.1 类型系统

#### 1. 类型信息的萃取

TypeOf()函数所做的事情如图5-21所示，就是找到传入参数的类型元数据，并以reflect.Type形式返回。

图5-21 由一个∗_type和一个∗itab组建一个iface

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P167_8253.jpg)

TypeOf()函数的代码

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P168_8267.jpg)

emptyInterface类型和5.1节介绍过的eface类型在内存布局上等价，emptyInterface类型定义

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P168_8275.jpg)

其中的rtype类型与runtime._type类型在内存布局方面也是等价的，只不过因为无法使用其他包中未导出的类型定义，所以需要在reflect包中重新定义一下。
代码中的eface.typ实际上就是从interface{}变量中提取出的类型元数据地址

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P168_8283.jpg)

- 先判断了一下传入的rtype指针是否为nil，如果不为nil就把它作为Type类型返回，否则返回nil。
- 从这里可以知道∗rtype类型肯定实现了Type接口，之所以要加上这个nil判断，需要考虑到Go的接口类型是个双指针结构，一个指向itab，另一个指向实际
  的数据对象。只有在两个指针都为nil的时候，接口变量才等于nil。

图5-22 萃取前判断非空

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P169_8297.jpg)

通过代码说明接口何时为空

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P168_8291.jpg)

在上述代码中第1个if处判断结果为真，所以会打印出1。第2个if处rw不再为nil，所以不会打印2。

interface{}中的类型元数据地址是从哪里来的呢？

- 当然是在编译阶段由编译器赋值的，实际的地址可能是由链接器填写的，也就是说源头还是要追溯到最初的源码中。

#### 2.类型系统的初始化

## 5.5 本章小结
