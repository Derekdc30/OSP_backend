package routes

import (
	"backend/handlers"
	"net/http"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("POST /api/check-token", handlers.HandleCheckToken)
	mux.HandleFunc("GET /api/surveys/{token}", handlers.HandleGetSurvey)
	mux.HandleFunc("POST /api/responses", handlers.HandleSubmitResponse)

	// Admin routes
	mux.HandleFunc("POST /api/admin/verify", handlers.HandleVerifyAdmin)
	mux.HandleFunc("GET /api/admin/surveys", handlers.HandleGetAllSurveys)
	mux.HandleFunc("DELETE /api/admin/surveys/{token}", handlers.HandleDeleteSurvey)
	mux.HandleFunc("PUT /api/admin/surveys/{token}", handlers.HandleUpdateSurvey)
	mux.HandleFunc("GET /api/admin/responses/{token}", handlers.HandleGetSurveyResponses)
	mux.HandleFunc("POST /api/admin/surveys", handlers.HandleCreateSurvey)

	// Static files
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/", fs)

	return mux
}
