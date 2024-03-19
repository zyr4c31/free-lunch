package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

type Restaurant struct {
	ID   int
	Name string
}

func createTables(db *sql.DB) (*sql.Rows, error) {
	var bytes []byte
	file, _ := os.Open("schema.sql")
	file.Read(bytes)
	rows, err := db.Query(
		`CREATE TABLE IF NOT EXISTS restaurants(
		id INTEGER PRIMARY KEY,
		name TEXT
		);`)
	return rows, err
}

func selectAllRestaurants(db *sql.DB) ([]Restaurant, error) {
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
