package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
	CREATE TABLE scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT "",
		title VARCHAR(128) NOT NULL DEFAULT "",
		comment TEXT(128) NOT NULL DEFAULT "",
		repeat VARCHAR(128) NOT NULL DEFAULT ""
);
	CREATE INDEX date_index ON scheduler (date);
	`

var install bool
var db *sql.DB

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	if err != nil {
		install = true
	}

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if install {
		_, err = db.Exec(schema)
		if err != nil {
			db.Close()
			return err
		}
	}
	defer db.Close()
	return nil
}
