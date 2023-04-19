---
title: "第1章 从这里开始, 起飞了"
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

{{< embedcode shell "echo-menu-v1.sh" >}}

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

- 标准输出的文件描述符为1
- 标准错误输出的文件描述符为2
- 标准输入的文件描述符则为0

分别重定向标准输出, 标准错误

```bash
$ ls -l /etc/hosts /nofile > ok.txt 2> error.txt
```

重定向标准输出, 标准错误到同一个文件

```bash
$ ls -l /etc/hosts /nofile &> test.txt
```

将标准错误重定向到标准输出或反过来

```bash
$ ls /nofile 2>&1
```

图1-3 ls命令对比

![](https://res.weread.qq.com/wrepub/epub_27741237_17)

```bash
$ echo "hello" 1>&2
```

图1-4 echo命令对比

![](https://res.weread.qq.com/wrepub/epub_27741237_18)

```bash
$ ls /etc/passwd /nofile >test.txt 2>&1
```

图1-5 标准输出与错误输出

![](https://res.weread.qq.com/wrepub/epub_27741237_19)

输出黑洞`/dev/null`

- 数据一旦导入黑洞将无法找回

用文件重定向输入

```bash
$ mail -s warning root@localhosts < /etc/hosts
```

用`<<`(Here Document)重定向输入

```shell
#!/usr/bin/env bash
#语法格式:
#命令 << 分隔符
#内容
#分隔符
#系统会自动将两个分隔符之间的内容重定向传递给前面的命令，作为命令的输入。
#注意：分隔符是什么都可以，但前后分隔符必须一致。推荐使用EOF(end of file)
mail -s warning root@localhost << EOF
This is content.
This is a test mail for redirect.
EOF
```

同时使用重定向输入, 重定向输出

```shell
#!/usr/bin/env bash
cat > /tmp/test.txt << HERE
该文件为测试文件。
测试完后，记得将该文件删除。
Welcome to Earth.
HERE
```

如果数据和EOF前有Tab, 可以用`<<-`来忽略Tab

```shell
#!/usr/bin/env bash

#不能屏蔽Tab键,缩进将作为内容的一部分被输出
#注意hello和world前面是tab键
cat << EOF
	hello
	world
EOF

#Tab键将被忽略,仅输出数据内容
cat <<- EOF
	hello
	world
EOF
```

## 1.5 各种引号的正确使用姿势

### 单引号与双引号

- `""`
  - 引用一个整体
- `''`
  - 引用一个整体
  - 不解析特殊字符
- `\`
  - 不解析随后的一个特殊字符

### 命令替换

使用``` `` ```

```bash
$ tar -czf  /root/log-`date +%Y%m%d`.tar.gz  /var/log/
```

使用`$()`

```bash
$ echo "当前系统账户登录数量: $(who|wc -l)"
```

## 1.6 千变万化的变量

- 变量名由字母, 数字, `_` 组成
- 不能用数字开头
- 赋值时`=`两边不能有空格.

表1-4 变量名示例

![](https://res.weread.qq.com/wrepub/epub_27741237_24)

- 使用`$var`或`${var}`方式读取变量值

删除变量

```bash
$ test=123
$ unset test
```

表1-5 常见的系统预设变量

![](https://res.weread.qq.com/wrepub/epub_27741237_27)

## 1.7 数据过滤与正则表达式

