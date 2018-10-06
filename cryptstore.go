package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
	_ "github.com/mattn/go-sqlite3"
)

func get(key string) bool {
	db, _ := sql.Open("sqlite3", "./data.db")
	defer db.Close()
	rows, err := db.Query("select id, fileid, ipfshash from filemaps")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var fileid string
		var ipfshash string
		err = rows.Scan(&id, &fileid, &ipfshash)
		if err != nil {
			log.Fatal(err)
		}
		if key == fileid {
			return true
		}
	}
	return false
}

func update(id int, key string, value string) {
	db, _ := sql.Open("sqlite3", "./data.db")
	defer db.Close()
	query, _ := db.Prepare("insert into filemaps (id, fileid, ipfshash) values (?, ?, ?)")
	defer query.Close()
	query.Exec(id, key, value)

}

func main() {

	fmt.Println("Storing encrypted files on the blockchain...Adding to Swarm.")
	sh := shell.NewShell("localhost:5001")

	e, _ := os.Open(".decrypt")
	defer e.Close()

	scanner := bufio.NewScanner(e)
	count := 0
	for scanner.Scan() {
		ef, _ := ioutil.ReadFile(strings.Split(scanner.Text(), ":")[0])
		hash, err := sh.Add(strings.NewReader(string(ef)))

		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s", err)
			os.Exit(1)
		}

		fmt.Printf("(%d) File:Hash = %s \n(%d) ipfs.Hash = %s", count, scanner.Text(), count, hash)
		if get(scanner.Text()) == true {
			continue
		} else {
			update(count, scanner.Text(), hash)
		}
		count++
	}
	fmt.Println("Process completed.")
}
