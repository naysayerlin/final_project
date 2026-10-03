package config

import "os"

type Config struct {
	Port     string
	DBFile   string
	Password string
}

const (
	defaultPort   = "7540"
	defaultDBFile = "scheduler.db"
)

func Load() Config {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}

	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = defaultDBFile
	}

	password := os.Getenv("TODO_PASSWORD")

	return Config{
		Port:     port,
		DBFile:   dbFile,
		Password: password,
	}
}
