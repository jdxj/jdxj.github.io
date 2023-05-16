---
title: "第3章 根本停不下来的循环和中断控制"
date: 2023-05-14T10:34:52+08:00
draft: true
tags:
  - ""
---

## 3.1 玩转for循环语句

name是可以任意定义的变量名称，word是支持扩展的项目列表，扩展后生成一份完整的项目列表（或值列表）。

```shell
for name [ in [ word ...]]
do
    命令序列
done

for i in 1 2 3 4 5
do
  echo "$i hello world"
done
```

```shell
# 相当于 for name in $@
for name
do
    命令序列
Done
```

Shell支持多种扩展，如变量替换、命令扩展、算术扩展、通配符扩展等。

```bash
# 不能使用变量

$ echo {1..5}                      #生成1～5的数字序列
1 2 3 4 5
$ echo {5..1}                      #生成5～1的数字序列
5 4 3 2 1
$ echo {1..10..2}
1 3 5 7 9

$ echo {a..z}                         #生成字母序列，小写字母表
a b c d e f g h i j k l m n o p q r s t u v w x y z
$ echo {A..Z}
A B C D E F G H I J K L M N O P Q R S T U V W X Y Z
$ echo {x, y{i, j}{1,2,3}, z}         #自动生成组合的字符串序列
x yi1 yi2 yi3 yj1 yj2 yj3 z
```

[seq]({{< ref "" >}})命令生成数字序列，并且可以调用其他变量，但该命令不支持生成字母序列

C风格的for

```shell
for ((expr1 ; expr2 ; expr3))
do
    命令序列
done
```

图3-3 执行流程

![](https://res.weread.qq.com/wrepub/epub_27741237_101)

## 3.2 实战案例：猴子吃香蕉的问题

## 3.3 实战案例：进化版HTTP状态监控脚本

## 3.4 神奇的循环嵌套

在循环嵌套时内层循环和外层循环使用的变量名不能相同。

```shell
#!/bin/bash
#功能描述(Description):显示1和2的所有排列组合.

for i in {1..2}
do
    for j in {1..2}
    do
        echo "${i}${j}"
    done
done
```

## 3.5 非常重要的IFS

