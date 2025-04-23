package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

// Message represents the data we'll send to the frontend
type Message struct {
	Text string
}

func main() {
	// Serve static files from the "static" directory
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Handle the root route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Redirect to our index page
		http.Redirect(w, r, "/static/index.html", http.StatusSeeOther)
	})

	// Handle form submissions
	http.HandleFunc("/greet", greetHandler)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func greetHandler(w http.ResponseWriter, r *http.Request) {
	// Only handle POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Get the name from the form
	name := r.FormValue("name")
	if name == "" {
		name = "Anonymous"
	}

	// Create a response message
	message := Message{
		Text: fmt.Sprintf("Hello, %s! Welcome to our simple website.", name),
	}

	// Parse the template
	tmpl, err := template.ParseFiles("static/response.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Execute the template with our message
	err = tmpl.Execute(w, message)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
