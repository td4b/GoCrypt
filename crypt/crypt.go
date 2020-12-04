package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
	"log"
)

// Signature is..
func Signature(data []byte) []byte {
	h := sha256.New()
	h.Write(data)
	return h.Sum(nil)
}

// Encrypt is..
func Encrypt(data []byte, passphrase []byte) []byte {
	block, err := aes.NewCipher(Signature(passphrase))
	if err != nil {
		log.Fatal(err)
	}
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	if err != nil {
		log.Fatal(err)
	}
	return ciphertext
}

// Decrypt is..
func Decrypt(data []byte, passphrase []byte) []byte {
	block, err := aes.NewCipher(Signature(passphrase))
	if err != nil {
		log.Fatal(err)
	}
	gcm, _ := cipher.NewGCM(block)
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, _ := gcm.Open(nil, nonce, ciphertext, nil)
	return plaintext
}
