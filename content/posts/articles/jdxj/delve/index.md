---
title: "delve"
date: 2023-07-10T15:02:21+08:00
tags:
  - go
---

# 安装

```bash
$ go install github.com/go-delve/delve/cmd/dlv@latest
```

# 使用

测试代码

{{< embedcode go "cmd/delve-demo1/main.go" >}}
{{< embedcode go "pkg/foo/foo.go" >}}

dlv debug启动调试

在macOS下，执行dlv debug前，我们需要首先通过以下命令赋予dlv使用系统调试API的权限：sudo /usr/sbin/DevToolsSecurity -enable。

设置断点break(b)

```bash
(dlv) b main.go:12
```

查看断点breakpoints(bp)

```bash
(dlv) bp
```

所谓条件断点，指的就是当满足某个条件时，被调试的目标程序才会在该断点处暂停。

```bash
(dlv) b b2 foo.go:6
(dlv) cond b2 sum > 10
```

执行程序continue(c)/下一断点

```bash
(dlv) c
```

查看数据

```bash
(dlv) whatis a
(dlv) p a
(dlv) regs
(dlv) locals
```

常用的查看命令

- print（简写为p）：输出源码中变量的值。
- whatis：输出后面的表达式的类型。
- regs：当前寄存器中的值。
- locals：当前函数栈本地变量列表（包括变量的值）。
- args：当前函数栈参数和返回值列表（包括参数和返回值的值）。
- examinemem（简写为x）：查看某一内存地址上的值。

next(n)断点处的下一行

step(s)单步调试(会进入函数)

输出函数调用栈信息stack(bt)

通过up和down命令，可以在函数调用栈的栈帧间进行跳转

如果要重启调试，无须退出Delve，只需执行restart（简写为r）

delve还支持在调试过程中修改变量的值，并手工调用函数set

```bash
(dlv) set a = 4
(dlv) call foo.Foo(a, b)
```

Delve还可以通过exec子命令直接调试已经构建完的Go二进制程序文件

```bash
$ dlv exec ./main.out
```

# 调试并发程序

通过Delve提供调试命令，我们可以在各个运行的goroutine间切换。

{{< embedcode go "cmd/delve-demo2/main.go" >}}
{{< embedcode go "pkg/bar/bar.go" >}}

启动调试

```bash
$ dlv debug
(dlv) b b1 main.go:19
(dlv) c
(dlv) goroutines # 查看goroutine列表
(dlv) goroutine 1 # 切换goroutine
```

Delve还提供了thread和threads命令，通过这两个命令我们可以查看当前启动的线程列表并在各个线程间切换

# 调试core dump文件

core dump文件是在程序异常终止或崩溃时操作系统对程序当时的内存状态进行记录并保存而生成的一个数据文件，该文件以core命名，也被称为核心转储文件。

- 适用于生产环境中的调试
- Delve目前支持对linux/amd64、linux/arm64架构下产生的core文件的调试，以及Windows/amd64架构下产生的minidump小转储文件的调试。

测试代码

{{< embedcode go "cmd/delve-demo3/main.go" >}}

要想在Linux下让Go程序崩溃时产生core文件，我们需要进行一些设置

```bash
$ulimit -c unlimited # 不限制core文件大小
$go build main.go
$GOTRACEBACK=crash ./main
# 会产生core文件
$ dlv core ./main ./core
(dlv) bt
(dlv) frame 11 # 跳到main.main函数栈帧
```

# 调试运行中的程序

调试器一旦成功挂接到正在运行的进程中，调试器就掌握了进程执行的指挥权，并且正在运行的goroutine都会暂停，等待调试器的进一步指令。

```bash
$ dlv attach <pid> ./main.out
```
