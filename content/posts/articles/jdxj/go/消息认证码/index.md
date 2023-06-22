---
title: "消息认证码"
date: 2023-06-22T13:45:53+08:00
tags:
  - go
  - crypto
---

单向散列函数虽然能辨别出数据是否被篡改，但却无法辨别出数据是不是伪装的。因此，在这样的场合下，我们还需要对消息进行认证（Authentication），即
**校验消息的来源**是不是我们所期望的。而用于解决这一问题的常见密码技术就是消息认证码（Message Authentication Code，MAC）。

消息认证码技术是以通信双方共享密钥为前提的。对于任意长度的消息，我们都可以计算出一个固定长度的消息认证码数据，这个数据被称为MAC值。

我们可以将消息认证码理解成一种与密钥相关联的单向散列函数。消息认证码有多种实现方式，包括使用单向散列函数实现、使用分组密码实现、公钥密码实现等。

使用SHA-256单向散列函数的HMAC示例

```go
// chapter9/sources/go-crypto/hmac_generate.go

func main() {
    // 密钥(key) 32字节
    key := []byte("12345678123456781234567812345678")

    // 要传递的消息
    message := []byte("I love go programming language!!")

    // 创建hmac实例（使用SHA-256单向散列函数）
    mac := hmac.New(sha256.New, key)
    mac.Write(message)

    // 计算mac值
    m := mac.Sum(nil)
    ms := fmt.Sprintf("%x", m) // mac到string
    fmt.Printf("mac值 = %s\n", ms)
}
```

在实际使用中，对数据进行对称加密且携带MAC值的方式被称为“认证加密”（Authenticated Encryption with Associated Data，AEAD）。认证加密同时
满足了机密性（对称加密）、完整性（MAC中的单向散列）以及认证（MAC）的特性，在生产中有着广泛的应用。认证加密主要有以下三种方式。

认证加密主要有以下三种方式

- Encrypt-then-MAC：先用对称密码对明文进行加密，然后计算密文的MAC值。
- Encrypt-and-MAC：将明文用对称密码加密，并计算明文的MAC值。
- MAC-then-Encrypt：先计算明文的MAC值，然后将明文和MAC值一起用对称密码加密。

分组密码中的GCM（Galois Counter Mode）就是一种认证加密模式，它使用CTR（计数器）分组模式和128比特分组长度的AES加密算法进行加密，并使用
Carter-Wegman MAC算法实现MAC值计算。
