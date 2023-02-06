---
title: "第3章 函数"
date: 2023-02-06T21:40:56+08:00
draft: true
---

图3-1 函数调用发生前

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P45_5231.jpg)

1. CALL指令会先把下一条指令的地址(返回地址)压入栈中, IP寄存器存储f1的地址

图3-2 CALL指令执行后

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P46_5242.jpg)

2. 执行f1()
3. f1()最后有条RET指令, 弹出栈顶的返回地址(应该弹到IP中), 跳到返回地址处继续执行

图3-3 RET指令执行后

![](https://res.weread.qq.com/wrepub/CB_3300047233_Figure-P46_5245.jpg)

## 3.1 栈帧
