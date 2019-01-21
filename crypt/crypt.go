package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
)

func createHash(key string) string {
	hasher := sha256.New()
	hasher.Write([]byte(key))
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash[0 : len(hash)/2]
}

func sha1Hash(key string) string {
	hasher := sha1.New()
	hasher.Write([]byte(key))
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash
}

func signature(data []byte) string {
	factor := len(data) / 256
	var sig string
	for i := 1; i <= 256; i++ {
		sig += string(data[factor*i : len(data)])
	}
	return sha1Hash(sig)
}

func Encrypt(data []byte, passphrase string) []byte {
	block, err := aes.NewCipher([]byte(createHash(passphrase)))
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

func Decrypt(data []byte, passphrase string) []byte {
	block, err := aes.NewCipher([]byte(createHash(passphrase)))
	if err != nil {
		log.Fatal(err)
	}
	gcm, _ := cipher.NewGCM(block)
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, _ := gcm.Open(nil, nonce, ciphertext, nil)
	return plaintext
}
