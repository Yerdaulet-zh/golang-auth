package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func main() {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)

	// Private Key
	privFile, _ := os.Create("private.pem")
	pem.Encode(privFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

	// Public Key
	pubFile, _ := os.Create("public.pem")
	pubBytes, _ := x509.MarshalPKIXPublicKey(&key.PublicKey)
	pem.Encode(pubFile, &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubBytes,
	})
}
