package main

import (
	"bufio"
	"fmt"
  "log"
	"io/ioutil"
	"os"
	"strings"
  "github.com/tidwall/buntdb"
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
    db, err := buntdb.Open("data.db")
    if err != nil {
    	log.Fatal(err)
    }
    err = db.Update(func(tx *buntdb.Tx) error {
	  _, _, err := tx.Set(scanner.Text(), hash, nil)
	  return err
    })
    defer db.Close()
    
	}

	fmt.Println("Process completed.")
}
