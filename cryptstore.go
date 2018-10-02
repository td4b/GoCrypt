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
	fmt.Println("Storing encrypted files on the blockchain...Adding to Swarm.")
	sh := shell.NewShell("localhost:5001")
	e, _ := os.Open(".decrypt")
  defer e.Close()
	scanner := bufio.NewScanner(e)
	for scanner.Scan() {
		ef, _ := ioutil.ReadFile(scanner.Text())
		hash, err := sh.Add(strings.NewReader(string(ef)))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s", err)
			os.Exit(1)
		}
		fmt.Printf("%s -> %s", scanner.Text(), hash)
		fmt.Println("\nPulling copy of file(s) from swarm.")
		sh.Get(hash, "./")
	}
	fmt.Println("Process completed.")
}
