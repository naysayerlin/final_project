package db

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"os"
	"strings"
	"sync"

	"modernc.org/sqlite"
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
var registerFuncsOnce sync.Once

func registerFuncs() {
	sqlite.MustRegisterDeterministicScalarFunction(
		"lower_unicode",
		1,
		func(ctx *sqlite.FunctionContext, args []driver.Value) (driver.Value, error) {
			s, ok := args[0].(string)
			if !ok {
				return nil, nil
			}
			return strings.ToLower(s), nil
		},
	)
}

func Init(dbFile string) error {
	registerFuncsOnce.Do(registerFuncs)
	_, err := os.Stat(dbFile)
	if err != nil {
		if os.IsNotExist(err) {
			install = true
		} else {
			return fmt.Errorf("can not check db file: %w", err)
		}
	}
	database, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("not able to open database: %w", err)
	}
	if install {
		if _, err = database.Exec(schema); err != nil {
			database.Close()
			return fmt.Errorf("not able to create schema: %w", err)
		}
	}
	db = database
	return nil
}

func GetDB() *sql.DB {
	return db
}

func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
