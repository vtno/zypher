package main

import (
	"fmt"
	"crypto/aes"
	"crypto/cipher"
)

func main() {
	input := "somelongkeysomelongkeysomelongkey"
	key := []byte("1111111111111111")
	if len(input) <= 16 {
		src := []byte(input)
		src = append(src, make([]byte, 16 - len(input))...)
		dst := encrypt(key, src)
		fmt.Printf("input: %s\n", input)
		fmt.Printf("encrypted: %s\n", string(dst))	
		dst = decrypt(key, []byte(dst))
		fmt.Printf("decrypted: %s\n", string(dst))
	} else {
		loopCont := (len(input) / 16) + 1
		leftOver:= len(input) % 16
		var encryptedResult = []byte{}
		for v := 0; v < loopCont; v++ {
			src := []byte(input[(v * 16):])
			fmt.Printf("src loop: %s\n", string(src))
			encryptedResult = append(encryptedResult, []byte(encrypt(key, src))...)
		}
		leftOverSrc := []byte(input[len(input) - leftOver:])
		fmt.Printf("left over src: %s\n", string(leftOverSrc))
		encryptedResult = append(encryptedResult, []byte(encrypt(key, leftOverSrc))...)
		fmt.Printf("input: %s\n", input)
		fmt.Printf("encrypted: %s\n", string(encryptedResult))

		loopCont = len(encryptedResult) / 16
		leftOver = len(encryptedResult) % 16
		var decryptedResult = []byte{}
		for v := 1; v < loopCont; v++ {
			src := encryptedResult[:(16*loopCont)]
			decryptedResult = append(decryptedResult, []byte(decrypt(key, src))...)
		}
		leftOverSrc = make([]byte, leftOver + (16 - leftOver))
		decryptedResult = append(decryptedResult, []byte(decrypt(key, leftOverSrc))...)
		fmt.Printf("decrypted: %s\n", string(decryptedResult))
	}
}

func encrypt(key, src []byte) string {
	dst := make([]byte, len(src))
	block := must(aes.NewCipher(key))
	block.Encrypt(dst, src)
	return string(dst)
}

func decrypt(key, src []byte) string {
	dst := make([]byte, len(src))
	block := must(aes.NewCipher(key))
	block.Decrypt(dst, []byte(src))
	return string(dst)
}

func must(b cipher.Block, err error) (cipher.Block) {
	if err != nil {
		panic(err)
	}
	return b
}