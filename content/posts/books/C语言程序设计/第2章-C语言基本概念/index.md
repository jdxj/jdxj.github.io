---
title: "第2章 C语言基本概念"
date: 2023-06-30T16:27:01+08:00
---

## 2.1 编写一个简单的C程序

{{< embedcode c "pun.c" >}}

### 2.1.1 编译和链接

对于C程序来说, 把程序转化为机器可以执行的形式通常包含下列3个步骤

- 预处理。首先程序会被送交给预处理器（preprocessor）。预处理器执行以#开头的命令（通常称为指令）。预处理器有点类似于编辑器，它可以给程序添
  加内容，也可以对程序进行修改。
- 编译。修改后的程序现在可以进入编译器（compiler）了。编译器会把程序翻译成机器指令（即目标代码）。然而，这样的程序还是不可以运行的。
- 链接。在最后一个步骤中，链接器（linker）把由编译器产生的目标代码和所需的其他附加代码整合在一起，这样才最终产生了完全可执行的程序。这些附加
  代码包括程序中用到的库函数（如printf函数）。

预处理器通常会和编译器集成在一起, 在UNIX系统环境下，通常把C编译器命名为cc

```bash
$ cc pun.c
$ gcc -o pun pun.c
```

### 2.1.2 集成开发环境

## 2.2 简单程序的一般形式

```
指令
int main(void)
{
  语句
}
```

### 2.2.1 指令

```c
#include <stdio.h>
```

### 2.2.2 函数

main函数是非常特殊的：在执行程序时系统会自动调用main函数。

- 它会在程序终止时向操作系统返回一个状态码

{{< embedcode c "pun.c" >}}

- main前面的int表明该函数将返回一个整数值
- 圆括号中的void表明main函数没有参数

```c
return 0
```

有两个作用：

- 一是使main函数终止（从而结束程序）
- 二是指出main函数的返回值是0。

如果main函数的末尾没有return语句，程序仍然能终止。但是，许多编译器会产生一条警告信息（因为函数应该返回一个整数却没有这么做）。

### 2.2.3 语句

```c
printf("To C, or not to C： that is the question.\n");
```

- 以`;`结束
- 指令不用`;`

### 2.2.4 显示字符串

```c
printf("Brevity is the soul of wit.\n  --Shakespeare\n");
```

## 2.3 注释

```c
/* This is a comment */
```

建议把*/放在单独一行

```c
/* Name: pun.c
   Purpose: Prints a bad pun.
   Author: K. N. King
*/
```

C99提供了另一种类型的注释, 这种风格的注释会在行末自动终止。

```c
// This is a comment
```

## 2.4 变量和赋值

### 2.4.1 类型

- 进行算术运算时float型变量通常比int型变量慢
- float型变量所存储的数值往往只是实际数值的一个近似值

### 2.4.2 声明

```c
int height;
float profit;

// 声明同类型的变量
int height, length, width, volume;
float profit, loss;
```

C99之前

```c
int main(void)
{
  声明 // 声明在语句之前
  语句
}
```

**在C99中，声明可以不在语句之前。**

### 2.4.3 赋值

```c
height = 8;
length = 12;
width = 10;
```

当我们把一个包含小数点的常量赋值给float型变量时，最好在该常量后面加一个字母f（代表float）：

- 不加f可能会引发编译器的警告。

```c
profit = 2150.48f;
```

混合类型赋值是可以的，但不一定安全

```c
height = 8;
length = 12;
width = 10;
volume = height * length * width;     /* volume is now 960 */
```

### 2.4.4 显示变量的值

```c
printf("Height: %d\n", height);
```

### 2.4.5 初始化

当程序开始执行时，某些变量会被自动设置为零，而大多数变量则不会

```c
int height = 8;

int height = 8, length = 12, width = 10;
```

### 2.4.6 显示表达式的值

```c
printf("%d\n", height * length * width);
```

## 2.5 读入输入

```c
scanf("%d", &i);  /* reads an integer; stores into i */
```

## 2.6 定义常量的名字

宏定义（macro definition）当对程序进行编译时，预处理器会把每一个宏替换为其表示的值。

```c
#define INCHES_PER_POUND 166
```

可以利用宏来定义表达式

```c
#define RECIPROCAL_OF_PI (1.0f / 3.14159f)
```

## 2.7 标识符

在C语言中，标识符可以含有字母、数字和下划线，但是必须以字母或者下划线开头。在C99中，标识符还可以使用某些“通用字符名”

不合法的标识符

```c
10times  get-next-char
```

- C语言是区分大小写的
- C对标识符的最大长度没有限制

### 关键字

![](https://res.weread.qq.com/wrepub/epub_31359737_32)

## 2.8 C程序的书写规范

大多数情况下，程序中记号之间的空格数量没有严格要求。除非两个记号合并后会产生第三个记号，否则在一般情况下记号之间根本不需要留有间隔。

- 事实上，添加足够的空格和空行可以使程序更便于阅读和理解。

非法的分隔

```c
printf("To C, or not to C:
that is the question.\n");    /*** WRONG ***/
```
