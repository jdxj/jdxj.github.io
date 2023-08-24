---
title: "nasm指令"
date: 2023-08-23T16:14:51+08:00
draft: true
tags:
  - x86
  - asm
---

# 声明空间

- db(declare byte) 声明字节

```nasm
; 声明5个初始值为0的字节
db 0,0,0,0,0
```

- dw(declare word) 声明字
- dd(declare double) 声明双字
- dq(declare quad word) 声明四字

# div 无符号除法指令

可以做两种类型的除法

- 16位被除数(必须放在ax中), 8位除数(通用寄存器/内存), 执行后商在al中, 余数在ah中

```nasm
div cl
div byte [0x0023]

; 使用标号
dividnd dw 0x3f0
divisor db 0x3f
...
mov ax, [dividnd] ; 相当于 mov ax, [0xf000], 相对于ds
div byte [divisor]; 相当于 div byte [0xf002]
```

- 32位被除数, 16位除数

因为16位的处理器无法直接提供32位的被除数，故要求被除数的高16位在寄存器DX中，低16位在寄存器AX中。
指令执行后，商在寄存器AX中，余数在寄存器DX中

![](https://res.weread.qq.com/wrepub/CB_3300050845_txt008_36.jpg)

```nasm
div cx
div word [0x0230]
```

# idiv 有符号除法

用法和div相同

```nasm
mov ax, 0x0400
mov bl, 0xf0
idiv bl
```

# xor 异或

目的操作数可以是通用寄存器和内存单元，源操作数可以是通用寄存器、内存单元和立即数（不允许两个操作数同时为内存单元）

```nasm
; 清零
xor dx, dx
```

`mov dx，0`的机器码是BA 00 00；而`xor dx，dx`的机器码则是31 D2，不但较短，而且，因为`xor dx，dx`的两个操作数都是通用寄存器，所以执行速
度最快。

![](https://res.weread.qq.com/wrepub/CB_3300050845_txt008_42.jpg)

# add

- 目的操作数可以是8位或者16位的通用寄存器，或者指向8位或者16位实际操作数的内存地址；
- 源操作数可以是相同数据宽度的8位或者16位通用寄存器、指向8位或者16位实际操作数的内存地址，或者立即数，但不允许两个操作数同时为内存单元。
- 相加后，结果保存在目的操作数中。

![](https://res.weread.qq.com/wrepub/CB_3300050845_txt008_44.jpg)

# sub

目的操作数可以是8位或者16位通用寄存器，也可以是8位或者16位的内存单元；源操作数可以是通用寄存器，也可以是内存单元或者立即数（不允许两个操作数
同时为内存单元）。

```nasm
sub ah, al
sub dx, ax
sub [label_a], ch
```

# jmp

绝对地址

```nasm
jmp 0x5000:0xf0c0
```

near关键字

```nasm
infi: jmp near infi
```
- 在编译阶段，编译器是这么做的：用标号（目标位置）处的汇编地址减去当前指令的下一条指令的汇编地址，就得到了jmp near infi指令的实际操作数。
- near用于截断相减结果为16位
- 在指令执行阶段，处理器用指令指针寄存器IP的内容（它已经指向下一条指令）加上该指令的操作数，就得到了要转移的实际偏移地址，同时寄存器CS的内容
  不变。因为改变了指令指针寄存器IP的内容，这直接导致处理器的指令执行流程转向目标位置。

# 条件转移指令

## jns

```nasm
; 如果未设置符号位，则转移到标号“show”所在的位置处执行。
jns show
```

它和jmp指令很相似，它也是相对转移指令，编译后的机器指令操作数也是一个相对偏移量，是用标号处的汇编地址减去当前指令的下一条指令的汇编地址得到的。

相反的指令js

## jz

结果为零（ZF标志为1）则转移；jnz的意思是结果不为零（ZF标志为0）则转移。

## jo

结果溢出（OF标志为1）则转移，jno的意思是结果未溢出（OF标志为0）则转移。

## jc

有进位（CF标志为1）则转移，jnc的意思是没有进位（CF标志为0）则转移。

## jp

如果PF标志为1则转移，jnp的意思是如果PF标志不为1（为0）则转移。

## 与cmp配合的条件转移指令

![](https://res.weread.qq.com/wrepub/CB_3300050845_txt009_52.jpg)
![](https://res.weread.qq.com/wrepub/CB_3300050845_txt009_53.jpg)

## jcxz(jump if CX is zero)

当寄存器CX的内容为零时则转移。

如果寄存器CX的内容为零，则转移到标号show；否则不转移，继续往下执行。

```nasm
jcxz show
```

# cmp

目的操作数可以是8位或者16位通用寄存器，也可以是8位或者16位内存单元；源操作数可以是与目的操作数宽度一致的通用寄存器、内存单元或者立即数，但两
个操作数同时为内存单元的情况除外。

```nasm
cmp al, 0x08
cmp dx, bx
cmp [label_a], cx
```

cmp指令在功能上和sub指令相同，唯一不同之处在于，cmp指令仅仅根据计算的结果设置相应的标志位，而不保留计算结果，因此也就不会改变两个操作数的原
有内容。cmp指令将会影响到CF、OF、SF、ZF、AF和PF标志位。

# times 重复后面的指令

```nasm
times 20 mov ax, bx
```

# 批量数据传送

- movsb
- movsw

movsb和movsw指令执行时，原始数据串的段地址由DS指定，偏移地址由SI指定，简写为DS:SI；要传送到的目的地址由ES:DI指定；传送的字节数(movsb)或
者字数(movsw)由CX指定。

除此之外，还要指定是正向传送还是反向传送，正向传送是指传送操作的方向是从内存区域的低地址端到高地址端；反向传送则正好相反。正向传送时，每传送
一字节(movsb)或者一个字(movsw)，SI和DI加1或者加2；反向传送时，每传送一字节(movsb)或者一个字(movsw)时，SI和DI减去1或者减去2。不管是正向
传送还是反向传送，也不管每次传送的是字节还是字，每传送一次，CX的内容自动减1

# 方向控制

- cld, 传送方向从低到高
- std, 传送方向从高到低

# rep 重复执行指令

cx不为零则重复

```nasm
rep movsw
```

# loop 循环

```nasm
digit: 
         xor dx,dx
         div si
         mov [bx],dl ;保存数位
         inc bx 
         loop digit
```

在执行loop时

- 将cx减1
- 如果cx不为零, 则转移到指定位置, 否则执行后面的指令

# inc 加1

```nasm
inc al
inc byte [bx]; ds
inc word [label_a]; ds
```

# dec 减1

# neg 求相反数

```nasm
neg al
neg dx
neg word [label_a]
```

# 扩展有符号数

- cbw(convert byte to word)
- cwd(convert word to double-world)

cbw操作al, 扩展到ax

- 如果寄存器AL中的内容为01001111，那么执行该指令后，寄存器AX中的内容为0000000001001111；
- 如果寄存器AL中的内容为10001101，执行该指令后，寄存器AX中的内容为1111111110001101。

cwd操作ax, 扩展到dx:ax

- 如果寄存器AX中的内容为0100111101111001，那么执行该指令后，寄存器DX中的内容为0000000000000000，寄存器AX中的内容不变；
- 如果寄存器AX中的内容为1000110110001011，那么执行该指令后，寄存器DX中的内容为1111111111111111，寄存器AX中的内容同样不变。

**事实上，符号位是数的一部分，和其他比特一起共同表示数的大小，同时又用来判断数的正负。**

# 其他

在8086处理器上，如果要用寄存器来提供偏移地址，只能使用寄存器BX、SI、DI、BP，不能使用其他寄存器。

INTEL8086处理器只允许以下几种基址寄存器和变址寄存器的组合

```nasm
[bx + si]
[bx + di]
[bp + si]
[bp + di]
```

