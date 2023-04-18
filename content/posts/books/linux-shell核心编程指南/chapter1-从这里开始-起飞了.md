---
title: "第1章 从这里开始 起飞了"
date: 2023-04-18T09:11:27+08:00
draft: true
---

## 1.1 脚本文件的书写格式

多行注释, `<<`后的字符串区分大小写

```shell
#!/usr/bin/env bash
<<comment
something
comment
```

## 1.2 脚本文件的各种执行方式

1. 脚本文件自身没有可执行权限

```bash
$ bash xxx.sh
$ sh xxx.sh
```

2. 脚本文件具有可执行权限

```bash
$ chmod +x xxx.sh
$ xxx.sh
```

3. 开启子进程执行的方式

不管是直接执行脚本，还是使用bash或sh这样的解释器执行脚本，都是会开启子进程的。

4. 不开启子进程的执行方式

```bash
$ source xxx.sh
# 或者使用 . xxx.sh
```

## 1.3 如何在脚本文件中实现数据的输入与输出

1. 使用echo命令创建一个脚本文件菜单

```bash
$ echo [选项] 字符串
```

{{< embedcode shell "chapter1-code/echo-menu-v1.sh" >}}

表1-1 常见转义符号

![](https://res.weread.qq.com/wrepub/epub_27741237_8)

2. 扩展知识，使用printf命令创建一个脚本菜单

```bash
$ printf [格式] 参数
```

表1-2 常用的格式字符串及功能描述

![](https://res.weread.qq.com/wrepub/epub_27741237_10)

```bash
$ printf "%d" 12
# 左对齐
$ printf "%-5d" 12
```

3. 使用read命令读取用户的输入信息

```bash
$ read [选项] [变量名]
```

- 如果未指定变量名，则默认变量名称为REPLY

表1-3 read命令常用的选项

![](https://res.weread.qq.com/wrepub/epub_27741237_12)

```bash
$ read input1 input2                       
abc def
$ echo $input1 $input2                     
abc def
```

## 1.4 输入与输出的重定向
