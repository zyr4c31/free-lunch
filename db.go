package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/tursodatabase/go-libsql"
)

func db() (*sql.DB, string, *libsql.Connector) {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}
	dbName := "local.db"
	primaryUrl := os.Getenv("TURSO_URL")
	authToken := os.Getenv("TURSO_TOKEN")

	dir, err := os.MkdirTemp(".", "libsql-*")
	if err != nil {
		log.Println("Error creating temporary directory: ", err)
		os.Exit(1)
	}

	dbPath := filepath.Join(dir, dbName)

	syncInterval := time.Minute
	connector, err := libsql.NewEmbeddedReplicaConnector(
		dbPath, primaryUrl, libsql.WithAuthToken(authToken), libsql.WithSyncInterval(syncInterval))
	if err != nil {
		log.Println("Error creating connector: ", err)
		os.Exit(1)
	}

	db := sql.OpenDB(connector)
	return db, dir, connector
}
