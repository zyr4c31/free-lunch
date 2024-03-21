package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/tursodatabase/go-libsql"
)

func newDBConnection() *sql.DB {
	dbName := "local.db"
	dir, err := os.MkdirTemp("", "libsql-*")
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		os.Exit(1)
	}
	dbPath := filepath.Join(dir, dbName)

	godotenv.Load()
	primaryUrl := "libsql://free-lunch-zyr4c31.turso.io"
	authToken := os.Getenv("TURSO_TOKEN")

	connectorStart := time.Now()
	connector, err := libsql.NewEmbeddedReplicaConnector(
		dbPath, primaryUrl, libsql.WithAuthToken(authToken), libsql.WithSyncInterval(time.Minute))
	if err != nil {
		fmt.Println("Error creating connector:", err)
		os.Exit(1)
	}
	connectorDuration := time.Since(connectorStart)
	log.Println("connected to turso in", connectorDuration)
	db := sql.OpenDB(connector)
	return db
}
