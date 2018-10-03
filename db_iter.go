package main

import (
	"fmt"
	"github.com/tidwall/buntdb"
	"log"
)

func main() {
	db, err := buntdb.Open("data.db")
	if err != nil {
		log.Fatal(err)
	}
	err = db.View(func(tx *buntdb.Tx) error {
		tx.Ascend("", func(key, value string) bool {
			fmt.Printf("key: %s, value: %s\n", key, value)
			return true
		})
		return err
	})
	defer db.Close()
}
