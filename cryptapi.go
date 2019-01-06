package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type jsondata struct {
	File string
	Hash string
}

func apicall(w http.ResponseWriter, r *http.Request) {
	connStr := "user=docker password=docker dbname=filehashes sslmode=disable"
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func main() {
	http.HandleFunc("/api/", apicall)
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", nil))
}
