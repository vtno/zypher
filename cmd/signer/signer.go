package main

import (
	"crypto"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
	"path"
)

func main() {
	file, err := os.ReadFile("zypher")
	if (err != nil) {
		fmt.Printf("error reading private key: %v", err)
		os.Exit(1)
	}
	block, _ := pem.Decode(file)
	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if (err != nil) {
		fmt.Printf("error parsing private key: %v", err)
		os.Exit(1)
	}
	if len(os.Args) != 2 {
		fmt.Printf("usage: %s <token>", path.Base(os.Args[0]))
		os.Exit(1)
	}
	hashed := sha256.Sum256([]byte(os.Args[1]))
	sig, err := privKey.Sign(rand.Reader, hashed[:], crypto.SHA256)
	if err != nil {
		fmt.Printf("error signing token: %v", err)
		os.Exit(1)
	}
	fmt.Printf("%s", base64.StdEncoding.EncodeToString(sig))
}
