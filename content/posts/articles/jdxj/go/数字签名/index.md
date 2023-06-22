---
title: "数字签名"
date: 2023-06-22T13:55:03+08:00
tags:
  - go
  - crypto
---

消息认证码虽然解决了消息发送者的身份认证问题，但由于采用消息认证码的通信双方共享密钥，因此对于一条通过了MAC验证的消息，通信双方依旧无法向第三方证
明这条消息就是对方发送的。同时任何一方也都没有办法防止对方否认该条消息是自己发送的。也就是说**单凭消息认证码无法防止否认**（non-repudiation）。

在消息认证码中，生成MAC和验证MAC使用的是同一密钥，这是无法防止否认问题的根源。因此数字签名技术对生成签名的密钥和验证签名的密钥进行了区分，签名密
钥只能由签名一方持有，它的所有通信对端将持有用于验证签名的密钥。

图55-10　签名与验证签名

![](https://res.weread.qq.com/wrepub/epub_42557147_63)

图55-10与公钥密码系统中的“公钥加密，私钥解密”的流程十分相似, 数字签名就是通过将公钥密码反过来用而实现的

图55-11　公钥密码流程与数字签名流程

![](https://res.weread.qq.com/wrepub/epub_42557147_64)

在实际生产应用中，我们通常对消息的摘要进行签名。这是因为公钥密码加密算法本身很慢，如果对消息全文进行加密将非常耗时。

RSA签名默认使用PKCS#1 v1.5方案，但该方案存在潜在伪造签名的可能。为了应对潜在伪造，RSA-PSS算法（Probabilistic Signature Scheme）被设计出
来。RSA-PSS算法通过采用对消息摘要进行签名，并在计算散列值时对消息加盐（salt）的方式来提高安全性（这样对同一条消息进行多次签名，每次得到的签名都
不同）。

```go
// chapter9/sources/go-crypto/cat rsa_pss_sign_and_verify.go

func main() {
    // 生成公钥密码的密钥对
    privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        panic(err)
    }
    publicKey := privateKey.PublicKey

    // 待签名消息
    msg := []byte("I love go programming language!!")

    // 计算摘要
    digest := sha256.Sum256(msg)

    // 用私钥签名
    sign, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, digest[:], nil)
    if err != nil {
        panic(err)
    }
    fmt.Printf("签名：%s\n", fmt.Sprintf("%x", sign))

    // 用公钥验证签名
    err = rsa.VerifyPSS(&publicKey, crypto.SHA256, digest[:], sign, nil)
    if err != nil {
        panic(err)
    }
    fmt.Printf("签名验证成功!\n")
}
```

我们看到数字签名既可以识别篡改和伪装，还可以防止否认，这让计算机网络通信从这一技术中获益匪浅。但数字签名的正确运用有一个大前提，那就是公钥属于真正
的发送者。
