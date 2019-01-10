package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, _ := sql.Open("sqlite3", "./data.db")
	sqlStmt := `
        create table filemaps (id integer not null primary key, fileid text, ipfshash text);
	delete from filemaps;
	`
	query, _ := db.Prepare(sqlStmt)
	query.Exec()
}
