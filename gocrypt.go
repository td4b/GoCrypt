package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
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

func encryptFile(filename string, passphrase string) []byte {
	data, _ := ioutil.ReadFile(filename)
	f, _ := os.Create(filename)
	defer f.Close()
	f.Write(encrypt(data, passphrase))
	return data
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

	fmt.Print("[Verify SecretKey]: ")
	vrfybytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	vrfysecret := string(vrfybytePassword)

	if vrfysecret != secret {
		fmt.Print("\nError, Key Verification Failed.\n")
		return
	}

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
				rawdata := encryptFile(scanner.Text(), secret)
				fmt.Println("\nEncrypting File: " + scanner.Text() + " Signature: " + signature(rawdata))
				filename := scanner.Text() + ":" + signature(rawdata) + "\n"
				d.Write([]byte(filename))
			}
			d.Close()
			e.Close()
			os.Remove(f.Name())
		case ".decrypt":
			d, _ := os.Open(f.Name())
			e, _ := os.Create(".encrypt")
			scanner := bufio.NewScanner(d)
			for scanner.Scan() {
				filename := strings.Split(scanner.Text(), ":")[0]
				decryptFile(filename, secret)
				name := strings.Split(scanner.Text(), ":")[0]
				signature := strings.Split(scanner.Text(), ":")[1]
				fmt.Println("\nDecrypting File: " + name + " Signature: " + signature)
				e.Write([]byte(filename + "\n"))
			}
			e.Close()
			d.Close()
			os.Remove(f.Name())
		}
	}
}
