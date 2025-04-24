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
	mux.HandleFunc("GET /api/surveys/{token}", handlers.HandleGetSurvey)
	// Static files
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/", fs)

	return mux
}
