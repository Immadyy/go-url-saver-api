package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

type CreateLinkRequest struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

type Link struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Link  string `json:"link"`
}

type Database struct {
	mu    sync.Mutex
	Links []Link
	ID    int64
}

type LinkService struct {
	Store LinkStore
}

type LinkStore interface {
	Save(data Link) (Link, error)
	GetAll() ([]Link, error)
	Update(upId int64, data Link) (Link, error)
	Delete(delId int64) (Link, error)
}

func NewLinkService(l LinkStore) *LinkService {
	return &LinkService{
		Store: l,
	}
}

type Application struct {
	LinkService *LinkService
}

func NewApplication(l *LinkService) *Application {
	return &Application{
		LinkService: l,
	}
}

func (d *Database) Save(data Link) (Link, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.ID++
	data.ID = d.ID
	d.Links = append(d.Links, data)
	return data, nil
}

func (d *Database) GetAll() ([]Link, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	data := make([]Link, len(d.Links))
	copy(data, d.Links)
	return data, nil
}

func (d *Database) Update(updId int64, data Link) (Link, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	for i := range d.Links {
		if d.Links[i].ID == updId {
			d.Links[i].Title = data.Title
			d.Links[i].Link = data.Link
			return d.Links[i], nil
		}
	}
	return Link{}, fmt.Errorf("data not found")
}

func (d *Database) Delete(delId int64) (Link, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for i := range d.Links {
		if d.Links[i].ID == delId {
			data := d.Links[i]
			d.Links = append(d.Links[:i], d.Links[i+1:]...)
			return data, nil
		}
	}
	return Link{}, fmt.Errorf("data not found")
}

func (l *LinkService) ValidateLink(data Link) (Link, error) {
	if data.Link == "" || data.Title == "" {
		return Link{}, fmt.Errorf("Title and link cannot be empty.")
	}

	if !strings.HasPrefix(data.Link, "http://") && !strings.HasPrefix(data.Link, "https://") {
		data.Link = "https://" + data.Link
	}

	if _, err := url.ParseRequestURI(data.Link); err != nil {
		return Link{}, fmt.Errorf("Bad URL format")
	}

	return data, nil
}

func (l *LinkService) CreateLink(data Link) (Link, error) {
	Data, err := l.ValidateLink(data)
	if err != nil {
		return Link{}, err
	}
	return l.Store.Save(Data)
}

func (l *LinkService) GetAllLinks() ([]Link, error) {
	data, err := l.Store.GetAll()
	return data, err
}

func (l *LinkService) UpdateLink(updId int64, data Link) (Link, error) {
	data, err := l.ValidateLink(data)
	if err != nil {
		return Link{}, err
	}

	link, err := l.Store.Update(updId, data)
	return link, err
}

func (l *LinkService) DeleteLink(delId int64) (Link, error) {
	data, err := l.Store.Delete(delId)
	return data, err
}

type APIResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (app *Application) SaveHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateLinkRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	link := Link{
		Title: req.Title,
		Link:  req.Link,
	}

	finalLink, err := app.LinkService.CreateLink(link)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(APIResponse{
		Message: "Link saved successfully.",
		Data:    finalLink,
	})
}

func (app *Application) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	data, err := app.LinkService.GetAllLinks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

type UpdateLinkRequest struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

func (app *Application) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var u UpdateLinkRequest
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

	link := Link{
		Title: u.Title,
		Link:  u.Link,
	}

	data, err := app.LinkService.UpdateLink(number, link)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(APIResponse{
		Message: "Link updated.",
		Data:    data,
	})

}

func (app *Application) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	delId := r.URL.Query().Get("id")
	number, err := strconv.ParseInt(delId, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	data, err := app.LinkService.DeleteLink(number)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(APIResponse{
		Message: "Link deleted",
		Data:    data,
	})

}

func Health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("All is good."))
}

func main() {
	DB := &Database{Links: []Link{}}
	lol := NewLinkService(DB)
	app := NewApplication(lol)

	r := http.NewServeMux()

	r.HandleFunc("GET /health", Health)
	r.HandleFunc("POST /save_url", app.SaveHandler)
	r.HandleFunc("GET /get_all", app.GetAllHandler)
	r.HandleFunc("PUT /update_link", app.UpdateHandler)
	r.HandleFunc("DELETE /delete_link", app.DeleteHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: LoggerMiddleWare(r),
	}

	go func() {
		log.Println("Server starting on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error: %s\n", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	log.Println("Starting gracefull shutdown")

	shutdownctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")

}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggerMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the ResponseWriter to capture the status
		wrapped := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(wrapped, r)

		// Now you have the actual status
		log.Printf("%d %s %s %v", wrapped.status, r.Method, r.URL.Path, time.Since(start))
	})
}
