package cryptstore

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
	_ "github.com/lib/pq" // Blank import required for pq lib.
)

// Get key for db values.
func Get(key string) bool {
	connStr := "gocryptdb-svc://docker:docker@db/filehashes?sslmode=disable"
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
			return false
		}
	}
	return false
}

func Update(key string, value string) {
	connStr := "gocryptdb-svc://docker:docker@db/filehashes?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	// need to find ID dynamically for now it's static..

	sqlStatement := `
	INSERT INTO filemaps (id, fileid, ipfshash)
	VALUES ($1, $2, $3)
	`
	id := 1
	_, err = db.Exec(sqlStatement, id, key, value)
	if err != nil {
		panic(err)
	}
}

func Store(filehash string, data []byte) {

	fmt.Println("Storing encrypted files on the blockchain...Adding to Swarm.")
	sh := shell.NewShell("ipfs-svc:5001")

	ipfshash, err := sh.Add(strings.NewReader(string(data)))

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}

	fmt.Printf("File:Hash = (%s) ipfs.Hash = %s", filehash, ipfshash)
	Update(filehash, ipfshash)
	fmt.Println("\nProcess completed.")
}
