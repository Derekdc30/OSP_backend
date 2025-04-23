package routes

import (
	"backend/handlers"
	"net/http"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("POST /api/post", handlers.HandlePost)
	mux.HandleFunc("POST /api/check-token", handlers.HandleCheckToken)

	// Static files
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// // Redirect root
	// mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
	// 	http.Redirect(w, r, "/static/index.html", http.StatusSeeOther)
	// })

	return mux
}
