package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"url_saver/internal/handlers"
	"url_saver/internal/middleware"
	"url_saver/internal/routes"
	"url_saver/internal/service"
	"url_saver/internal/store"
)

func main() {
	//memoryStore := &store.MemoryStore{}
	connStr := "postgres://immady:kabutar@localhost:5432/dbname?sslmode=disable"
	PostgresStore, err := store.NewPostgresStore(connStr)
	if err != nil {
		log.Fatal("Database could'nt start: ", err)
	}
	linkService := service.NewLinkService(PostgresStore)
	handler := handlers.NewHandler(linkService)

	r := http.NewServeMux()
	routes.RegisterRoutes(r, handler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: middleware.LoggerMiddleWare(r),
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
	log.Println("Starting graceful shutdown")

	shutdownctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")

}
