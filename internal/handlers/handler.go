package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"url_saver/internal/models"
	"url_saver/internal/service"
)

type Handler struct {
	LinkService *service.LinkService
}

func NewHandler(l *service.LinkService) *Handler {
	return &Handler{
		LinkService: l,
	}
}

func (app *Handler) SaveHandler(w http.ResponseWriter, r *http.Request) {
	var req models.CreateLinkRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	link := models.Link{
		Title: req.Title,
		Link:  req.Link,
	}

	finalLink, err := app.LinkService.CreateLink(link)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.APIResponse{
		Message: "Link saved successfully.",
		Data:    finalLink,
	})
}

func (app *Handler) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	data, err := app.LinkService.GetAllLinks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (app *Handler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var u models.UpdateLinkRequest
	updateID := r.URL.Query().Get("id")

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	number, err := strconv.ParseInt(updateID, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	link := models.Link{
		Title: u.Title,
		Link:  u.Link,
	}

	data, err := app.LinkService.UpdateLink(number, link)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.APIResponse{
		Message: "Link updated.",
		Data:    data,
	})

}

// func (app *Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
// 	delId := r.URL.Query().Get("id")
// 	number, err := strconv.ParseInt(delId, 10, 64)
// 	if err != nil {
// 		http.Error(w, "invalid id", http.StatusBadRequest)
// 		return
// 	}

// 	data, err := app.LinkService.DeleteLink(number)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(models.APIResponse{
// 		Message: "Link deleted",
// 		Data:    data,
// 	})

// }

func (app *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("All is good."))
}
