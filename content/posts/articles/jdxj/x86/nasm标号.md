---
title: "nasm标号"
date: 2023-08-23T15:59:43+08:00
tags:
  - x86
  - asm
---

```nasm
infi: jmp near infi
```

标号可以由字母、数字、`_`、`$`、`#`、`@`、`~`、`.`、`?`组成，但必须以字母、`.` `_`和`?`中的任意一个打头。

# $标记

```nasm
jmp near $; 等同于
infi: jmp near infi
```

# $$标记

代表当前汇编节（段）的起始汇编地址。
