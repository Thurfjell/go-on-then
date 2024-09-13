package main

import (
	"context"
	"goonthen/internal/data"
	"goonthen/internal/game"
	"goonthen/internal/webroot"
	"log"
	"mime"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	_ = mime.AddExtensionType(".js", "text/javascript")
}

func main() {
	gameService, err := game.NewService(data.NewInMemoryState())

	if err != nil {
		log.Panic(err.Error())
	}

	rootHandler := webroot.New(
		webroot.WithRouterMeta(gameService.Routes()),
	)

	server := http.Server{
		Addr:        "localhost:1234",
		Handler:     rootHandler,
		IdleTimeout: 30 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("server stopped")
		}
		log.Printf("Server running at %s\n", server.Addr)
	}()

	// Tidy tidy tidy
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("shut down error: %v", err)
	}

}
