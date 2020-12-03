package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
	"main.go/crypt"
	"main.go/cryptstore"
)

type jsondata struct {
	File string
	Hash string
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// This needs to have better URL parsing implemented...
	keys := strings.Split(r.URL.Path, "/")[2]
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	data := crypt.Encrypt(bodyBytes, strings.Split(keys, "=")[1])
	w.Write([]byte("Received File, Encrypting Data.\n"))
	fileid := strings.Split(keys, "=")[0] + ":" + crypt.Signature(bodyBytes)
	if cryptstore.Get(fileid) == true {
		fmt.Println("File has already been uploaded to database..")
	} else {
		fmt.Println("Didnt find file Signature.. Uploading..")
		cryptstore.Store(fileid, data)
	}
}
func apicall(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		switch r.RequestURI {
		case "/api/":
			connStr := "gocryptdb-svc://docker:docker@db/filehashes?sslmode=disable"
			db, err := sql.Open("gocryptdb-svc", connStr)
			if err != nil {
				log.Fatal(err)
			}
			defer db.Close()

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
				data := jsondata{fileid, ipfshash}
				js, err := json.Marshal(data)
				if err != nil {
					log.Fatal(err)
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(js)

			}
		}
	} else {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
}

func main() {
	log.Println("[*] Starting GoCrypt -- Server listening on -- 8000")
	http.HandleFunc("/api/", apicall)
	http.HandleFunc("/upload/", uploadHandler)
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", nil))
}
