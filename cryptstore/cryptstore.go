package cryptstore

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
	_ "github.com/lib/pq"
)

func Get(key string) bool {
	connStr := "postgres://docker:docker@192.168.0.215/filehashes?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT fileid, ipfshash FROM filemaps")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var fileid string
		var ipfshash string
		err = rows.Scan(&fileid, &ipfshash)
		if err != nil {
			log.Fatal(err)
		}
		if key == fileid {
			return true
		}
	}
	return false
}

func Update(key string, value string) {
	connStr := "postgres://docker:docker@192.168.0.215/filehashes?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	//sqlStatement := `
	//INSERT INTO filemaps (fileid, ipfshash)
	//VALUES ($1, $2)
	//RETURNING id`
	//err = db.QueryRow(sqlStatement, key, value).Scan(&key)
	//if err != nil {
	//	panic(err)
	//}
}

func Store(filehash string, data []byte) {

	fmt.Println("Storing encrypted files on the blockchain...Adding to Swarm.")
	sh := shell.NewShell("datanode:5001")

	ipfshash, err := sh.Add(strings.NewReader(string(data)))

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}

	fmt.Printf("File:Hash = (%s) ipfs.Hash = %s", filehash, ipfshash)
	Update(filehash, ipfshash)
	fmt.Println("\nProcess completed.")
}
