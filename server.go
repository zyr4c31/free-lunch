package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/zyr4c31/free-lunch/sqlc"
)

func newServer(db *sql.DB) *http.Server {

	sm := http.NewServeMux()
	fs := http.FileServer(http.Dir("."))
	sm.Handle("GET /", fs)

	sm.HandleFunc("GET /restaurants", func(w http.ResponseWriter, r *http.Request) {
		queries := sqlc.New(db)
		restaurants, err := queries.ListRestaurants(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
		template := Restaurants(restaurants)
		template.Render(r.Context(), w)
	})

	sm.HandleFunc("POST /restaurants", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		form := r.PostForm
		if form.Has("name") == false {
			w.Header().Add("HX-Retarget", "#error")
			template := htmlError("no input")
			template.Render(context.Background(), w)
			return
		}

		name := form.Get("name")

		queries := sqlc.New(db)
		err := queries.CreateRestaurant(context.Background(), name)
		if err != nil {
			log.Fatalln(err)
		}
		w.Header().Add("Content-Type", "text/html")
		w.Header().Add("HX-Refresh", "true")
	})

	sm.HandleFunc("GET /restaurants/{id}", func(w http.ResponseWriter, r *http.Request) {
		pathRestaurantID := r.PathValue("id")
		r.ParseForm()
		form := r.Form
		menuItem := form.Get("menu-item")
		price := form.Get("price")

		restaurantID, _ := strconv.ParseInt(pathRestaurantID, 10, 64)

		queries := sqlc.New(db)
		menuItems, err := queries.ListMenuItemsForRestaurant(r.Context(), restaurantID)
		if err != nil {
			htmlError(err.Error()).Render(r.Context(), w)
			return
		}

		Menu(menuItem, price, menuItems).Render(r.Context(), w)
	})

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalln(err)
	}
	port := os.Getenv("PORT")
	addr := fmt.Sprintf("%v:%v", hostname, port)
	log.Println("hosted on", addr)
	server := http.Server{
		Addr:    addr,
		Handler: logging(sm),
	}

	return &server
}
