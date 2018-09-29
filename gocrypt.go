package main

import (
        "bufio"
        "crypto/aes"
        "crypto/cipher"
        "crypto/md5"
        "crypto/rand"
        "encoding/hex"
        "fmt"
        "io"
        "io/ioutil"
        "log"
        "os"
)

func createHash(key string) string {
        hasher := md5.New()
        hasher.Write([]byte(key))
        return hex.EncodeToString(hasher.Sum(nil))
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

        reader := bufio.NewReader(os.Stdin)
        fmt.Print("[PrivateKey]: ")
        secret, _ := reader.ReadString('\n')
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
                                fmt.Println("Encrypting File: " + scanner.Text())
                                encryptFile(scanner.Text(), secret)
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
                                fmt.Println("Decrypting File: " + scanner.Text())
                                decryptFile(scanner.Text(), secret)
                                e.Write([]byte(scanner.Text()))
                        }
                        e.Close()
                        d.Close()
                        os.Remove(f.Name())
                }
        }
}
