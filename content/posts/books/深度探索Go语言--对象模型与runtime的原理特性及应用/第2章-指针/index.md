---
title: "第2章 指针"
date: 2023-02-05T13:55:47+08:00
draft: true
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
