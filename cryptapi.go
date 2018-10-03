package main

import (
	"fmt"
	"github.com/tidwall/buntdb"
	"log"
	"net/http"
)

func apicall(w http.ResponseWriter, r *http.Request) {
	db, err := buntdb.Open("data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.View(func(tx *buntdb.Tx) error {
		tx.Ascend("", func(key, value string) bool {
			data := "key: " + string(key) + " value: " + string(value)
			fmt.Fprintf(w, data)
			return true
		})
		return err
	})
}

func main() {
	http.HandleFunc("/api/", apicall)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
