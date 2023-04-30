package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

func main() {
	input := "somelongkeysomelongkeysomelongkey"
	key := []byte("1111111111111111")
	encryptedText := encrypt(key, []byte(input))
	fmt.Printf("input: %s\n encrypted: %s\n", input, encryptedText)
	decryptedText := decrypt(key, []byte(encryptedText))
	fmt.Printf("input: %s\n decrypted: %s\n", encryptedText, decryptedText)
}

func encrypt(key, src []byte) string {
	c := must(aes.NewCipher(key))
	gcm := must(cipher.NewGCM(c))
	nonce := make([]byte, gcm.NonceSize())
    must(io.ReadFull(rand.Reader, nonce))
	return string(gcm.Seal(nonce, nonce, src, nil))
}

func decrypt(key, src []byte) string {
    c := must(aes.NewCipher(key))

    gcm := must(cipher.NewGCM(c))
    nonceSize := gcm.NonceSize()
    if len(src) < nonceSize {
        panic(errors.New("ciphertext too short"))
    }

    nonce, ciphertext := src[:nonceSize], src[nonceSize:]
    return string(must(gcm.Open(nil, nonce, ciphertext, nil)))
}

func must[T any](b T, err error) (T) {
	if err != nil {
		panic(err)
	}
	return b
}