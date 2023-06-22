---
title: "Go密码学包"
date: 2023-06-22T00:46:40+08:00
tags:
  - go
  - crypto
---

官方crypto包

- src/crypto
- golang.org/x/crypto

分类

- 分组密码
  - cipher 五种分组模式
    - ECB
    - CBC
    - CFB
    - OFB
    - CTR
  - des 对称密码
    - DES
    - TDEA
  - aes 对称密码
    - AES
- 公钥密码
  - tls
    - TLS 1.2
    - TLS 1.3
  - x509
    - 编码格式的密钥和证书的解析
  - rsa
    - RSA
  - elliptic
    - 标准椭圆曲线算法
  - dsa
    - DSA
  - ecdsa
    - 数字签名算法
  - ed25519
    - 椭圆曲线签名算法Ed25519
- 单向散列函数 (消息摘要)
  - md5
  - sha1
  - sha256
  - sha512
- 消息认证码
  - hmac
  - rand
