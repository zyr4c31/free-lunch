package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/tursodatabase/go-libsql"
	"github.com/zyr4c31/url-params/sqlc"
)

type Restaurant struct {
	ID   int
	Name string
}

func queryRestaurants(db *sql.DB) []sqlc.Restaurant {
	queries := sqlc.New(db)

	restaurants, err := queries.ListRestaurants(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to execute query: %v\n", err)
		os.Exit(1)
	}

	return restaurants
}

func main() {
	godotenv.Load()
	dbName := "local.db"
	primaryUrl := "libsql://free-lunch-zyr4c31.turso.io"
	authToken := os.Getenv("TURSO_TOKEN")

	dir, err := os.MkdirTemp("", "libsql-*")
	if err != nil {
		fmt.Println("Error creating temporary directory:", err)
		os.Exit(1)
	}
	defer os.RemoveAll(dir)

	dbPath := filepath.Join(dir, dbName)

	connectorStart := time.Now()
	connector, err := libsql.NewEmbeddedReplicaConnector(dbPath, primaryUrl,
		libsql.WithAuthToken(authToken),
	)
	if err != nil {
		fmt.Println("Error creating connector:", err)
		os.Exit(1)
	}
	defer connector.Close()
	connectorDuration := time.Since(connectorStart)
	log.Println(connectorDuration)

	openStart := time.Now()
	db := sql.OpenDB(connector)
	defer db.Close()
	openDuration := time.Since(openStart)
	log.Println(openDuration)

	queryStart := time.Now()
	restaurants := queryRestaurants(db)
	queryDuration := time.Since(queryStart)
	log.Println(queryDuration)
	totalDuration := time.Since(connectorStart)
	log.Println(totalDuration)
	log.Print(restaurants)

	sm := http.NewServeMux()

	server := http.Server{
		Addr:    ":8080",
		Handler: sm,
	}

	log.Fatalln(server.ListenAndServe())
}
