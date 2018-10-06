package main

import (
	"log"
	"net/http"
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
)

type Data struct {
  File    string
  Hash    string
}

func apicall(w http.ResponseWriter, r *http.Request) {
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
		data := Data{fileid,ipfshash}
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
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
