package main

import (
  "log"
  "github.com/tidwall/buntdb"
  "fmt"
)

func main() {

    db, err := buntdb.Open("data.db")
    if err != nil {
      log.Fatal(err)
    }
    //err = db.Update(func(tx *buntdb.Tx) error {
    //_, _, err := tx.Set("Hello", "World", nil)
    //return err
    //})

    err = db.View(func(tx *buntdb.Tx) error {
  	val, err := tx.Get("Hello")
  	if err != nil{
  		return err
  	}
  	fmt.Printf("value is %s\n", val)
  	return nil
    })

    defer db.Close()

}
