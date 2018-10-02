package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

func createHash(key string) string {
	hasher := sha256.New()
	hasher.Write([]byte(key))
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash[0 : len(hash)/2]
}

func encrypt(data []byte, passphrase string) []byte {
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

func decrypt(data []byte, passphrase string) []byte {
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

func encryptFile(filename string, passphrase string) {
	data, _ := ioutil.ReadFile(filename)
	f, _ := os.Create(filename)
	defer f.Close()
	f.Write(encrypt(data, passphrase))
}

func decryptFile(filename string, passphrase string) {
	data, _ := ioutil.ReadFile(filename)
	f, _ := os.Create(filename)
	defer f.Close()
	plaintext := decrypt(data, passphrase)
	f.Write(plaintext)
}

func main() {
	fmt.Print("[SecretKey]: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	secret := string(bytePassword)
	fmt.Println()

	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		switch f.Name() {
		case ".encrypt":
			e, _ := os.Open(f.Name())
			d, _ := os.Create(".decrypt")
			scanner := bufio.NewScanner(e)
			for scanner.Scan() {
				encryptFile(scanner.Text(), secret)
				fmt.Println("Encrypting File: " + scanner.Text())
				d.Write([]byte(scanner.Text()))
			}
			d.Close()
			e.Close()
			os.Remove(f.Name())
		case ".decrypt":
			d, _ := os.Open(f.Name())
			e, _ := os.Create(".encrypt")
			scanner := bufio.NewScanner(d)
			for scanner.Scan() {
				decryptFile(scanner.Text(), secret)
				fmt.Println("Decrypting File: " + scanner.Text())
				e.Write([]byte(scanner.Text()))
			}
			e.Close()
			d.Close()
			os.Remove(f.Name())
		}
	}
}
