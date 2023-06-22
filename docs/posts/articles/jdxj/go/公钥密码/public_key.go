package public_key

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func RSAEncrypt() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	publicKey := privateKey.PublicKey

	plaintext := []byte("hello world")
	// 使用公钥加密
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &publicKey, plaintext, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("ciphertext: %s\n", hex.EncodeToString(ciphertext))

	// 使用私钥解密
	plaintext, err = rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, ciphertext, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("plaintext: %s\n", plaintext)
}
