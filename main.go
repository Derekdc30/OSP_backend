package main

import (
	"backend/db"
	"backend/routes"
	"fmt"
	"log"
	"net/http"
)

// main initializes database and starts HTTP server
func main() {
	db.InitDB()
	// Create router with all routes
	router := routes.NewRouter()
	// Start server
	fmt.Println("Server listening on :8080...")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", router))
}
