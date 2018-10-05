package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/tidwall/buntdb"
)

func main() {
	// will break these up into functions eventually.
	fmt.Println("Storing encrypted files on the blockchain...Adding to Swarm.")
	sh := shell.NewShell("localhost:5001")

	e, _ := os.Open(".decrypt")
	defer e.Close()

	scanner := bufio.NewScanner(e)
	for scanner.Scan() {
		ef, _ := ioutil.ReadFile(strings.Split(scanner.Text(), ":")[0])
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
		sh.Get(hash, "./")
		err = db.View(func(tx *buntdb.Tx) error {
			_, err := tx.Get(scanner.Text())
			if err != nil {
				return err
			}
			// issue here, find out why element is not added.
			err = db.Update(func(tx *buntdb.Tx) error {
				fmt.Println("test")
				_, _, err := tx.Set(scanner.Text(), hash, nil)
				return err
			})

			return nil
		})

		defer db.Close()

	}

	fmt.Println("\nProcess completed.")
}
