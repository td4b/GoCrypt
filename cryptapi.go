package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type jsondata struct {
	File string
	Hash string
}

func getdata() {
	var w http.ResponseWriter
	connStr := "postgres://docker:docker@db/filehashes?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
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

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Create("./result")
	if err != nil {
		panic(err)
	}
	n, err := io.Copy(file, r.Body)
	if err != nil {
		panic(err)
	}

	w.Write([]byte(fmt.Sprintf("%d bytes are recieved.\n", n)))
}

func apicall(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		switch r.RequestURI {
		case "/api/":
			getdata()
		}
	} else {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
}

func main() {
	http.HandleFunc("/api/", apicall)
	http.HandleFunc("/upload/", uploadHandler)
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", nil))
}
