package block_cipher

import (
	"encoding/base64"
	"fmt"
	"testing"
)

func TestAESEncrypt(t *testing.T) {
	plaintext := []byte("hell world")
	key, err := RandBytes(32)
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	ciphertext, err := AESEncrypt(plaintext, key)
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	fmt.Printf("ciphertext: %s\n", base64.StdEncoding.EncodeToString(ciphertext))

	plaintextStr, err := AESDecrypt(ciphertext, key)
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	fmt.Printf("plaintext: %s\n", plaintextStr)
}
