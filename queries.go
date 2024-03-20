package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/tursodatabase/go-libsql"
)

type TursoEmbeddedReplicaConfig struct {
	dbName       string
	primaryUrl   string
	authToken    string
	syncInterval time.Duration
}

type Restaurant struct {
	ID   int
	Name string
}

func (tursoDB *TursoEmbeddedReplicaConfig) createTables() (sql.Result, error) {
	dir, err := os.MkdirTemp(".", "libsql-*")
	if err != nil {
		log.Println("Error creating temporary directory: ", err)
		os.Exit(1)
	}
	defer os.RemoveAll(dir)
	dbPath := filepath.Join(dir, tursoDB.dbName)
	connector, err := libsql.NewEmbeddedReplicaConnector(
		dbPath, tursoDB.primaryUrl, libsql.WithAuthToken(tursoDB.authToken), libsql.WithSyncInterval(tursoDB.syncInterval))
	if err != nil {
		log.Println("Error creating connector: ", err)
		os.Exit(1)
	}
	defer connector.Sync()
	defer connector.Close()
	db := sql.OpenDB(connector)
	defer db.Close()

	result, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS restaurants(
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS menu_items(
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		price DOUBLE NOT NULL,
		restaurant_id INTEGER NOT NULL,
		FOREIGN KEY(restaurant_id) REFERENCES restaurants(id)
		);`)
	return result, err
}

func (tursoDB *TursoEmbeddedReplicaConfig) selectAllRestaurants() ([]Restaurant, error) {
	dir, err := os.MkdirTemp(".", "libsql-*")
	if err != nil {
		log.Println("Error creating temporary directory: ", err)
		os.Exit(1)
	}
	defer os.RemoveAll(dir)
	dbPath := filepath.Join(dir, tursoDB.dbName)
	connector, err := libsql.NewEmbeddedReplicaConnector(
		dbPath, tursoDB.primaryUrl, libsql.WithAuthToken(tursoDB.authToken), libsql.WithSyncInterval(tursoDB.syncInterval))
	if err != nil {
		log.Println("Error creating connector: ", err)
		os.Exit(1)
	}
	db := sql.OpenDB(connector)
	defer connector.Sync()
	defer connector.Close()

	var restaurants []Restaurant
	rows, err := db.Query("select * from restaurants")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to execute query: %v\n", err)
		return restaurants, err
	}
	defer rows.Close()

	for rows.Next() {
		var restaurant Restaurant

		if err := rows.Scan(&restaurant.ID, &restaurant.Name); err != nil {
			log.Println("Error scanning row: ", err)
			return restaurants, err
		}

		restaurants = append(restaurants, restaurant)
		log.Println(restaurant.ID, restaurant.Name)

		if err := rows.Err(); err != nil {
			log.Println("Error during rows iteration: ", err)
		}
	}
	return restaurants, err
}
