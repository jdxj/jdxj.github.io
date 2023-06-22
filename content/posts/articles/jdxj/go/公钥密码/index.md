---
title: "公钥密码"
date: 2023-06-22T12:49:30+08:00
tags:
  - go
  - crypto
---

常见的密钥配送方案有事先共享密钥（事先以安全的方式将密钥交给通信方）、密钥分配中心（每个通信方要事先与密钥分配中心共享密钥）、Diffie-Hellman密钥
交换算法、公钥密码等。

图55-5　公钥密钥分发与加解密流程

![](https://res.weread.qq.com/wrepub/epub_42557147_58)

图55-6　RSA密钥对参与加解密的原理图

![](https://res.weread.qq.com/wrepub/epub_42557147_59)

- RSA加解密默认使用PKCS#1 v1.5填充方案，但该方案在面对Chosen Ciphertext Attacks（选择密文攻击）时强度不足（虽然无法破译RSA，但攻击者可能
  获取到密文对应的明文的少量信息）
- RSA-OAEP（Optimal Asymmetric Encryption Padding，最优非对称加密填充）则被认为是一种可信赖、满足强度要求的填充方案。


# 使用RAS加密

{{< embedcode go "public_key.go" >}}

- rsa.EncryptOAEP和rsa.DecryptOAEP的第二个参数都是一个随机数生成器（这里传入rand.Reader），RSA-OAEP会通过随机数使每次生成的密文呈现不同
  的排列方式，因此多次运行上述示例程序所得到的密文结果都是不同的。
- 这两个函数的第一个参数是hash.Hash接口实现的实例，其产生的散列值可作为随机数生成器的种子。这两个函数需要采用同一种hash.Hash接口的实现，Go标准
  库文档推荐使用sha256.New()。

**RSA算法对待处理的数据长度是有要求的**，采用RSA-OAEP填充时，加密函数EncryptOAEP支持的最大明文长度为

- RSA密钥长度（字节数）-单向散列结果长度×2-2。
- 在上面的例子中，加密函数支持的最大明文长度为256-32×2-2=190（字节）。
