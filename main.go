package main

import (
	"fmt"
)

func main() {
	input := "somelongkeysomelongkeysomelongkey"
	key := "1111111111111111"
	ci := NewCipher(key)
	encryptedText := must(ci.Encrypt([]byte(input)))
	fmt.Printf("input: %s\n encrypted: %s\n", input, encryptedText)
	decryptedText := must(ci.Decrypt([]byte(encryptedText)))
	fmt.Printf("input: %s\n decrypted: %s\n", encryptedText, decryptedText)
}

func must[T any](b T, err error) T {
	if err != nil {
		panic(err)
	}
	return b
}
