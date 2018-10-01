package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
)

func main() {
	fmt.Println("Storing encrypted files on the blockchain...")
	// Where your local node is running on localhost:5001
	sh := shell.NewShell("localhost:5001")
	e, _ := os.Open(".decrypt")
	scanner := bufio.NewScanner(e)
	for scanner.Scan() {
		ef, _ := ioutil.ReadFile(scanner.Text())
		hash, err := sh.Add(strings.NewReader(string(ef)))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s", err)
			os.Exit(1)
		}
		fmt.Printf("File %s Added to Swarm Hash: %s", scanner.Text(), hash)
		fmt.Println("\nPulling copy of file from swarm.")
		sh.Get(hash, "./")
	}

	fmt.Println("\nProcess completed.")
}
