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
	_ "github.com/lib/pq"
)

func get(key string) bool {
	connStr := "user=docker password=docker dbname=filehashes sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT id, fileid, ipfshash FROM filemaps")
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
	connStr := "user=docker password=docker dbname=filehashes sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	sqlStatement := `
	INSERT INTO filemaps (id, fileid, ipfshash)
	VALUES ($1, $2, $3)
	RETURNING id`
	err = db.QueryRow(sqlStatement, id, key, value).Scan(&id)
	if err != nil {
		panic(err)
	}
}

func main() {

	fmt.Println("Storing encrypted files on the blockchain...Adding to Swarm.")
	sh := shell.NewShell("datanode:5001")

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

		fmt.Printf("\n(%d) File:Hash = %s \n(%d) ipfs.Hash = %s", count, scanner.Text(), count, hash)
		if get(scanner.Text()) == true {
			continue
		} else {
			update(count, scanner.Text(), hash)
		}
		count++
	}
	fmt.Println("\nProcess completed.")
}
