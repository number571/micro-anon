package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"os"
)

func main() {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile("priv.key", x509.MarshalPKCS1PrivateKey(privKey), 0600); err != nil {
		panic(err)
	}
	if err := os.WriteFile("pub.key", x509.MarshalPKCS1PublicKey(&privKey.PublicKey), 0600); err != nil {
		panic(err)
	}
}
