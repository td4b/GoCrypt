package createdb

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func create() {
	db, _ := sql.Open("sqlite3", "./data.db")
	sqlStmt := `
        create table filemaps (fileid text, ipfshash text);
	delete from filemaps;
	`
	query, _ := db.Prepare(sqlStmt)
	query.Exec()
}
