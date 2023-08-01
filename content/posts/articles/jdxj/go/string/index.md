---
title: "Go String"
date: 2023-08-01T10:28:51+08:00
tags:
  - go
---

# String的特点

## string类型的数据是不可变的

一旦声明了一个string类型的标识符，无论是常量还是变量，该标识符所指代的数据在整个程序的生命周期内便无法更改。

{{< embedcode go "immutable1/main.go" >}}

对string进行切片化后，Go编译器会为切片变量重新分配底层存储而不是共用string的底层存储，因此对切片的修改并未对原string的数据产生任何影响。

{{< embedcode go "immutable2/main.go" >}}

对string的底层的数据存储区仅能进行只读操作，一旦试图修改那块区域的数据程序将崩溃.

## 零值可用

```go
var s string
fmt.Println(s) // s = ""
fmt.Println(len(s)) // 0
```

## 获取长度的时间复杂度是O(1)级别

Go string类型数据是不可变的，因此一旦有了初值，那块数据就不会改变，其长度也不会改变。Go将这个长度作为一个字段存储在运行时的string类型的内部表示
结构中. len(s)实际上就是读取存储在运行时中的那个长度值

## 支持通过+/+=操作符进行字符串连接

```go
s := "Rob Pike, "
s = s + "Robert Griesemer, "
s += " Ken Thompson"

fmt.Println(s) // Rob Pike, Robert Griesemer, Ken Thompson
```

## 支持各种比较关系操作符：==、!= 、>=、<=、>和<

```go
// chapter3/sources/string_compare.go

func main() {
    // ==
    s1 := "世界和平"
    s2 := "世界" + "和平"
    fmt.Println(s1 == s2) // true

    // !=
    s1 = "Go"
    s2 = "C"
    fmt.Println(s1 != s2) // true

    // < 和 <=
    s1 = "12345"
    s2 = "23456"
    fmt.Println(s1 < s2)  // true
    fmt.Println(s1 <= s2) // true

    // > 和 >=
    s1 = "12345"
    s2 = "123"
    fmt.Println(s1 > s2)  // true
    fmt.Println(s1 >= s2) // true
}
```

由于Go string是不可变的，因此如果两个字符串的长度不相同，那么无须比较具体字符串数据即可断定两个字符串是不同的。

## 对非ASCII字符提供原生支持

Go语言源文件默认采用的Unicode字符集。

{{< embedcode go "nonascii/main.go" >}}

## 原生支持多行字符串

Go语言直接提供了通过反引号构造“所见即所得”的多行字符串的方法

```go
// chapter3/sources/string_multilines.go

const s = `好雨知时节，当春乃发生。
随风潜入夜，润物细无声。
野径云俱黑，江船火独明。
晓看红湿处，花重锦官城。`

func main() {
    fmt.Println(s)
}
```

# string的内部表示

```go
// $GOROOT/src/runtime/string.go
type stringStruct struct {
    str unsafe.Pointer
    len int
}
```

runtime包中实例化一个字符串对应的函数

```go
// $GOROOT/src/runtime/string.go

func rawstring(size int) (s string, b []byte) {
    p := mallocgc(uintptr(size), nil, false)
    stringStructOf(&s).str = p
    stringStructOf(&s).len = size

    *(*slice)(unsafe.Pointer(&b)) = slice{p, size, size}

    return
}
```

![](https://res.weread.qq.com/wrepub/epub_42557145_32)

rawstring调用后，新申请的内存区域还未被写入数据，b就是供后续运行时层向其中写入数据（"hello"）用的。写完数据后，该slice就可以被回收掉了

# 字符串的高效构造

Go还提供了其他一些构造字符串的方法

- fmt.Sprintf
- strings.Join
- strings.Builder
- bytes.Buffer

结论

- 在能预估出最终字符串长度的情况下，使用预初始化的strings.Builder连接构建字符串效率最高；
- strings.Join连接构建字符串的平均性能最稳定，如果输入的多个字符串是以[]string承载的，那么strings.Join也是不错的选择；
- 使用操作符连接的方式最直观、最自然，在编译器知晓欲连接的字符串个数的情况下，使用此种方式可以得到编译器的优化处理；
- fmt.Sprintf虽然效率不高，但也不是一无是处，如果是由多种不同类型变量来构建特定格式的字符串，那么这种方式还是最适合的。

# 字符串相关的高效转换

> 注意新版本的Go可能提供相关转换

string和[]rune、[]byte可以双向转换

```go
// chapter3/sources/string_slice_to_string.go
func main() {
    rs := []rune{
        0x4E2D,
        0x56FD,
        0x6B22,
        0x8FCE,
        0x60A8,
    }

    s := string(rs)
    fmt.Println(s)

    sl := []byte{
        0xE4, 0xB8, 0xAD,
        0xE5, 0x9B, 0xBD,
        0xE6, 0xAC, 0xA2,
        0xE8, 0xBF, 0x8E,
        0xE6, 0x82, 0xA8,
    }

    s = string(sl)
    fmt.Println(s)
}

/*输出
中国欢迎您
中国欢迎您
*/
```

无论是string转slice还是slice转string，转换都是要付出代价的，这些代价的根源在于string是不可变的，运行时要为转换后的类型分配新内存。

在Go运行时层面，字符串与rune slice、byte slice相互转换对应的函数如下

```
// $GOROOT/src/runtime/string.go
slicebytetostring: []byte -> string
slicerunetostring: []rune -> string
stringtoslicebyte: string -> []byte
stringtoslicerune: string -> []rune
```

slice类型是不可比较的，而string类型是可比较的，因此在日常Go编码中，我们会经常遇到将slice临时转换为string的情况。

在运行时中有一个名为slicebytetostringtmp的函数就是协助实现这一优化的

```go
// $GOROOT/src/runtime/string.go
func slicebytetostringtmp(b []byte) string {
    if raceenabled && len(b) > 0 {
        racereadrangepc(unsafe.Pointer(&b[0]),
            uintptr(len(b)),
            getcallerpc(),
            funcPC(slicebytetostringtmp))
    }
    if msanenabled && len(b) > 0 {
        msanread(unsafe.Pointer(&b[0]), uintptr(len(b)))
    }
    return *(*string)(unsafe.Pointer(&b))
}
```

使用这个函数的前提是：在原slice被修改后，这个string不能再被使用了。因此这样的优化是针对以下几个特定场景的

1. string(b)用在map类型的key中

```go
b := []byte{'k', 'e', 'y'}
m := make(map[string]string)
m[string(b)] = "value"
```

2. string(b)用在字符串连接语句中

```go
b := []byte{'t', 'o', 'n', 'y'}
s := "hello " + string(b) + "!"
```

3. string(b)用在字符串比较中

```go
s := "tom"
b := []byte{'t', 'o', 'n', 'y'}

if s < string(b) {
    ...
}
```

Go编译器对用在for-range循环中的string到[]byte的转换也有优化处理，它不会为[]byte进行额外的内存分配

```go
// chapter3/sources/string_for_range_covert_optimize.go

func convert() {
    s := "中国欢迎您，北京欢迎您"
    sl := []byte(s)
    for _, v := range sl {
        _ = v
    }
}
func convertWithOptimize() {
    s := "中国欢迎您，北京欢迎您"
    for _, v := range []byte(s) {
        _ = v
    }
}

func main() {
    fmt.Println(testing.AllocsPerRun(1, convert))
    fmt.Println(testing.AllocsPerRun(1, convertWithOptimize))
}

/*输出
1
0
*/
```
