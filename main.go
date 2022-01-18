package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/gorilla/mux"
)

func initializeRouter() {
	r := mux.NewRouter()

	// get the orders list
	r.HandleFunc("/api/orders", getOrders).Methods("GET")
	r.HandleFunc("/searchOrders/{pn}", searchOrders).Methods("GET")

	log.Fatal(http.ListenAndServe(":1000", r))

}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	initialMigration()
	initializeRouter()
}
