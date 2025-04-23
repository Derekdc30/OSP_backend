package main

import (
	"backend/routes"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Create router with all routes
	router := routes.NewRouter()

	// Start server
	fmt.Println("Server listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
