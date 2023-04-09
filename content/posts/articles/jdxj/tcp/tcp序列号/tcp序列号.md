---
title: "TCP序列号/确认号"
date: 2023-04-01T18:23:15+08:00
tags:
  - tcp
---

# 序列号 Sequence Number

![](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/7b57f577bc8b460fbf804c728538d230~tplv-k3u1fbpfcp-zoom-in-crop-mark:3024:0:0:0.awebp)

序列号指的是本报文段第一个字节的序列号

- 32位无符号整数

![](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/d658980af614477296e95ae5f3f658f9~tplv-k3u1fbpfcp-zoom-in-crop-mark:3024:0:0:0.awebp)

# 初始序列号

在建立连接之初，通信双方都会各自选择一个序列号，称之为初始序列号。在建立连接时，通信双方通过 SYN 报文交换彼此的 ISN

![](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/0d499ad8e97c4c48942e39a7bd195652~tplv-k3u1fbpfcp-zoom-in-crop-mark:3024:0:0:0.awebp)

## 初始序列号是如何生成的

```c
__u32 secure_tcp_sequence_number(__be32 saddr, __be32 daddr,
				 __be16 sport, __be16 dport)
{
	u32 hash[MD5_DIGEST_WORDS];

	net_secret_init();
	hash[0] = (__force u32)saddr;
	hash[1] = (__force u32)daddr;
	hash[2] = ((__force u16)sport << 16) + (__force u16)dport;
	hash[3] = net_secret[15];
	
	md5_transform(hash, net_secret);

	return seq_scale(hash[0]);
}

static u32 seq_scale(u32 seq)
{
	return seq + (ktime_to_ns(ktime_get_real()) >> 6);
}
```

- 代码中的 net_secret 是一个长度为16的 int 数组，只有在第一次调用 net_secret_init 的时时候会将将这个数组的值初始化为随机值。在系统重启前保
  持不变。
- 可以看到初始序列号的计算函数 secure_tcp_sequence_number() 的逻辑是通过源地址、目标地址、源端口、目标端口和随机因子通过 MD5 进行进行计算。
  如果仅有这几个因子，对于四元组相同的请求，计算出的初始序列号总是相同，这必然有很大的安全风险，所以函数的最后将计算出的序列号通过 seq_scale 函
  数再次计算。
- seq_scale 函数加入了时间因子，对于四元组相同的连接，序列号也不会重复了。

## 序列号回绕了怎么处理

```c
static inline bool before(__u32 seq1, __u32 seq2)
{
    return (__s32)(seq1-seq2) < 0;
}
```

测试

{{< embedcode go "./main.go" >}}

输出

```
case1(未回绕):
seq1: 255, signed: -1
seq2: 1, signed: 1
befort: true

case2(已回绕):
seq1: 255, signed: -1
seq2: 128, signed: -128
befort: false
```

# 确认号 Acknowledgment Number

![](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/09910086f7284eb8892f8531fb8dc24a~tplv-k3u1fbpfcp-zoom-in-crop-mark:3024:0:0:0.awebp)
![](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/b90a1a9ade134affbaa67ea20b9dd3b1~tplv-k3u1fbpfcp-zoom-in-crop-mark:3024:0:0:0.awebp)

- 不是所有的包都需要确认的
- 不是收到了数据包就立马需要确认的，可以延迟一会再确认
- ACK 包本身不需要被确认，否则就会无穷无尽死循环了
- 确认号永远是表示小于此确认号的字节都已经收到
