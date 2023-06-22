package block_cipher

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

func RandBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	return b, err
}

func PKCS7Padding(plaintext []byte, blockSize int) []byte {
	paddingSize := blockSize - len(plaintext)%blockSize
	padding := bytes.Repeat([]byte{byte(paddingSize)}, paddingSize)
	return append(plaintext, padding...)
}

func PKCS7UnPadding(padded []byte) []byte {
	size := len(padded)
	paddingSize := int(padded[size-1])
	return padded[:size-paddingSize]
}

func AESEncrypt(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv, err := RandBytes(aes.BlockSize)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCEncrypter(block, iv)

	plaintext = PKCS7Padding(plaintext, aes.BlockSize)
	// iv+ciphertext
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	copy(ciphertext[:aes.BlockSize], iv)
	blockMode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)
	return ciphertext, nil
}

func AESDecrypt(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv := ciphertext[:aes.BlockSize]
	blockMode := cipher.NewCBCDecrypter(block, iv)

	plaintext := make([]byte, len(ciphertext)-aes.BlockSize)
	blockMode.CryptBlocks(plaintext, ciphertext[aes.BlockSize:])
	return PKCS7UnPadding(plaintext), nil
}
