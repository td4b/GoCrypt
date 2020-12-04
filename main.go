package main

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
	"main.go/crypt"
	"main.go/cryptstore"
)

// Data .. json payload
// type Data struct {
// 	Hashes Hashes
// }

// Hashes .. json payload
type Hashes struct {
	File string
	Hash string
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// This needs to have better URL parsing implemented...
	keys := strings.Split(r.URL.Path, "/")[2]
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	data := crypt.Encrypt(bodyBytes, []byte(strings.Split(keys, "=")[1]))
	w.Write([]byte("200 OK\n"))
	fileid := strings.Split(keys, "=")[0] + ":" + hex.EncodeToString(crypt.Signature(bodyBytes))
	// Only upload if data key value (name:SHA256) does not exist.
	if cryptstore.Get(fileid) == true {
		log.Println("[*] " + r.RemoteAddr + "- Got Request but discarding.")
	} else {
		log.Println("[*] " + r.RemoteAddr + "- Got Request, Uploading File.")
		cryptstore.Store(r.RemoteAddr, fileid, data)
	}
}
func apicall(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		switch r.RequestURI {
		case "/api/":
			connStr := "postgres://docker:docker@gocryptdb-svc/filehashes?sslmode=disable"
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

				data := Hashes{fileid, ipfshash}
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
	log.Fatal(http.ListenAndServe("0.0.0.0:8888", nil))
}
