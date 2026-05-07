package routes

import (
	"net/http"
	"url_saver/internal/handlers"
)

func RegisterRoutes(mux *http.ServeMux, h *handlers.Handler) {

	mux.HandleFunc("GET /health", h.HealthHandler)
	mux.HandleFunc("POST /save_url", h.SaveHandler)
	mux.HandleFunc("GET /get_all", h.GetAllHandler)
	mux.HandleFunc("PUT /update_link", h.UpdateHandler)
	mux.HandleFunc("DELETE /delete_link", h.DeleteHandler)
}
